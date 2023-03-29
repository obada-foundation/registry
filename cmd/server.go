package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/obada-foundation/registry/api"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"go.uber.org/zap"
)

// ServerCommand with command line flags and env vars
type ServerCommand struct {
	Port    int    `long:"port" env:"SERVER_PORT" default:"80" description:"port"`
	Address string `long:"address" env:"SERVER_ADDRESS" default:"" description:"listening address"`

	// Timeouts
	ReadTimeout     time.Duration `long:"read-timeout" env:"READ_TIMEOUT" default:"5s" description:"read timeout"`
	WriteTimeout    time.Duration `long:"write-timeout" env:"WRITE_TIMEOUT" default:"10s" description:"write timeout"`
	IdleTimeout     time.Duration `long:"idle-timeout" env:"IDLE_TIMEOUT" default:"120s" description:"idle timeout"`
	ShutdownTimeout time.Duration `long:"shutdown-timeout" env:"SHUTDOWN_TIMEOUT" default:"20s" description:"shutdown timeout"`

	SentryDSN string `long:"sentry-dsn" env:"SENTRY_DSN" default:"" description:"sentry dsn"`

	CommonOpts
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

	db, err := db.NewDBConnection(ctx, db.Connection{})
	if err != nil {
		return fmt.Errorf("cannot enstalish connection to immudb: %w", err)
	}
	s.Logger.Infow("startup", "status", "immudb connection established")

	didDocSvc := diddoc.NewService(db, s.Logger)

	apiServer := s.makeAPIServer(api.MuxConfig{
		Shutdown: shutdown,
		Log:      s.Logger,

		// Services
		DIDDoc: didDocSvc,
	})

	go func() {
		s.Logger.Infow("startup", "status", "api router started", "host", apiServer.Addr)
		serverErrors <- apiServer.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("api error: %w", err)

	case sig := <-shutdown:
		s.Logger.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer s.Logger.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, s.ShutdownTimeout)
		defer cancel()

		if err := apiServer.Shutdown(ctx); err != nil {
			if er := apiServer.Close(); er != nil {
				err = fmt.Errorf("%w; %v", err, er)
			}

			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		if err := db.CloseSession(ctx); err != nil {
			return fmt.Errorf("cannot close db connection: %w", err)
		}
	}

	return nil
}

func (s *ServerCommand) makeAPIServer(cfg api.MuxConfig) *http.Server {
	apiMux := api.Mux(cfg)

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler:      apiMux,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
		ErrorLog:     zap.NewStdLog(s.Logger.Desugar()),
	}
}
