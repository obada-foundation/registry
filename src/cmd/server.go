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

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	err := sentry.Init(sentry.ClientOptions{
		Dsn: s.SentryDSN,
	})
	if err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}

	apiServer := s.makeApiServer(api.APIMuxConfig{
		Shutdown: shutdown,
		Log:      s.Logger,
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

		ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
		defer cancel()

		if err := apiServer.Shutdown(ctx); err != nil {
			apiServer.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

func (s *ServerCommand) makeApiServer(cfg api.APIMuxConfig) *http.Server {
	apiMux := api.APIMux(cfg)

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler:      apiMux,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
		IdleTimeout:  s.IdleTimeout,
		ErrorLog:     zap.NewStdLog(s.Logger.Desugar()),
	}
}
