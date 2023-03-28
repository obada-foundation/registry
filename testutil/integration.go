package testutil

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/obada-foundation/registry/api"
	"github.com/obada-foundation/registry/services/diddoc"
)

// NewIntegrationTest prepares dependencies for the integration test
// nolint
func NewIntegrationTest(t *testing.T, didDoc diddoc.DIDDoc) (*httptest.Server, func()) {
	shutdown := make(chan os.Signal, 1)

	logger, lgDefer := NewTestLoger()

	if didDoc == nil {
		didDoc = diddoc.NewService(logger)
	}

	apiMux := api.Mux(api.MuxConfig{
		Shutdown: shutdown,
		Log:      logger,

		DIDDoc: didDoc,
	})

	srv := httptest.NewServer(apiMux)

	teardown := func() {
		srv.Close()
		lgDefer()
	}

	return srv, teardown
}
