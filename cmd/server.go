package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"embed"
	"github.com/getsentry/sentry-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/obada-foundation/registry/api"
	pbacc "github.com/obada-foundation/registry/api/pb/v1/account"
	pbdidoc "github.com/obada-foundation/registry/api/pb/v1/diddoc"
	"github.com/obada-foundation/registry/services/account"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	//	"go.uber.org/zap"
)

//go:embed swagger-ui/*
var swaggerUI embed.FS

// ServerCommand with command line flags and env vars
type ServerCommand struct {
	Port    int    `long:"port" env:"SERVER_PORT" default:"2017" description:"port"`
	Address string `long:"address" env:"SERVER_ADDRESS" default:"" description:"listening address"`

	HTTPServer HTTPServerGroup `group:"http" namespace:"http" env-namespace:"HTTP"`

	GrpcServer GrpcServerGroup `group:"grpc" namespace:"grpc" env-namespace:"GRPC"`

	// Database connection
	DB ImmuDBGroup `group:"db" namespace:"db" env-namespace:"DB"`

	// Timeouts
	ReadTimeout     time.Duration `long:"read-timeout" env:"READ_TIMEOUT" default:"5s" description:"read timeout"`
	WriteTimeout    time.Duration `long:"write-timeout" env:"WRITE_TIMEOUT" default:"10s" description:"write timeout"`
	IdleTimeout     time.Duration `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"120s" description:"idle timeout"`
	ShutdownTimeout time.Duration `long:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" default:"20s" description:"shutdown timeout"`

	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN" default:"" description:"sentry dsn"`

	CommonOpts
}

// HTTPServerGroup contains http server configuration options
type HTTPServerGroup struct {
	Port    int    `long:"port" env:"PORT" default:"2017" description:"port"`
	Address string `long:"address" env:"ADDRESS" default:"" description:"listening address"`
}

// GrpcServerGroup contains grpc server configuration options
type GrpcServerGroup struct {
	Port    int    `long:"port" env:"PORT" default:"2018" description:"port"`
	Address string `long:"address" env:"ADDRESS" default:"" description:"listening address"`
}

// ImmuDBGroup db connection details
type ImmuDBGroup struct {
	Host   string `long:"host" env:"HOST" default:"localhost" description:"immudb host"`
	Port   int    `long:"port" env:"PORT" default:"3322" description:"immudb port"`
	User   string `long:"user" env:"USER" default:"immudb" description:"immudb user"`
	Pass   string `long:"password" env:"PASSWORD" default:"immudb" description:"immudb password"`
	DBName string `long:"dbname" env:"NAME" default:"defaultdb" description:"immudb database name"`
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	//glog.Infof("preflight request for %s", r.URL.Path)
}

// Execute is the entry point for "server" command, called by flag parser
func (s *ServerCommand) Execute(_ []string) error {
	s.Logger.Infow("startup", "status", "server initialization started")

	ctx := context.Background()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	err := sentry.Init(sentry.ClientOptions{
		Dsn: s.SentryDSN,
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}

	dbClient, err := db.NewDBConnection(ctx, db.Connection{
		Host:   s.DB.Host,
		Port:   s.DB.Port,
		User:   s.DB.User,
		Pass:   s.DB.Pass,
		DBName: s.DB.DBName,
	})
	if err != nil {
		return fmt.Errorf("cannot enstalish connection to immudb: %w", err)
	}
	s.Logger.Infow("startup", "status", "immudb connection established")

	didDocSvc := diddoc.NewService(dbClient, s.Logger)

	accountSvc := account.NewService(dbClient, s.Logger)

	grpcServer, srv := api.NewGRPCServer(api.GRPCConfig{
		Log: s.Logger,

		// Services
		DIDDocService:  didDocSvc,
		AccountService: accountSvc,
	})

	grpcMux := runtime.NewServeMux()

	if er := pbacc.RegisterAccountHandlerServer(ctx, grpcMux, srv); er != nil {
		return er
	}

	if er := pbdidoc.RegisterDIDDocHandlerServer(ctx, grpcMux, srv); er != nil {
		return er
	}

	grpcListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.GrpcServer.Address, s.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("cannot select listener address for grpc server: %w", err)
	}

	go func() {
		s.Logger.Infow("startup", "status", "grpc server started", "host", grpcListener.Addr().String())
		serverErrors <- grpcServer.Serve(grpcListener)
	}()

	mux := http.NewServeMux()
	mux.Handle("/", allowCORS(grpcMux))
	mux.Handle("/swagger-ui/", http.FileServer(http.FS(swaggerUI)))

	go func() {
		listenAddr := fmt.Sprintf("%s:%d", s.HTTPServer.Address, s.HTTPServer.Port)
		s.Logger.Infow("startup", "status", "http started", "host", listenAddr)
		// nolint:gosec //check for later, needs refactoring to support timeouts
		serverErrors <- http.ListenAndServe(listenAddr, mux)
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("api error: %w", err)

	case sig := <-shutdown:
		s.Logger.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer s.Logger.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		shutDownCtx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
		defer cancel()

		grpcServer.GracefulStop()

		// if err := apiServer.Shutdown(shutDownCtx); err != nil {
		//	if er := apiServer.Close(); er != nil {
		//		err = fmt.Errorf("%w; %v", err, er)
		//	}

		//	return fmt.Errorf("could not stop server gracefully: %w", err)
		// }

		if err := dbClient.CloseSession(shutDownCtx); err != nil {
			return fmt.Errorf("cannot close db connection: %w", err)
		}
	}

	return nil
}
