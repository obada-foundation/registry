package testutil

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/obada-foundation/registry/api"
)

// NewIntegrationTest prepares dependencies for the integration test
// nolint
func NewIntegrationTest(t *testing.T) (*httptest.Server, func()) {
	shutdown := make(chan os.Signal, 1)

	logger, lgDefer := NewTestLoger()

	apiMux := api.Mux(api.MuxConfig{
		Shutdown: shutdown,
		Log:      logger,
	})

	srv := httptest.NewServer(apiMux)

	teardown := func() {
		srv.Close()
		lgDefer()
	}

	return srv, teardown
}
