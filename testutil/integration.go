package testutil

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/obada-foundation/registry/api"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"github.com/stretchr/testify/require"
)

// NewIntegrationTest prepares dependencies for the integration test
// nolint
func NewIntegrationTest(t *testing.T, c *Container, didDoc diddoc.DIDDoc) (*httptest.Server, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	shutdown := make(chan os.Signal, 1)

	logger, lgDefer := NewTestLoger()

	conn := db.Connection{
		Host:   c.Host,
		Port:   c.Port,
		User:   "immudb",
		Pass:   "immudb",
		DBName: "defaultdb",
	}

	dbClient, err := db.NewDBConnection(ctx, conn)
	require.NoErrorf(t, err, "No connection with docker container %+v %+v", c, conn)

	err = dbClient.HealthCheck(ctx)
	require.NoError(t, err, "immudb is not healthy")

	if didDoc == nil {
		didDoc = diddoc.NewService(dbClient, logger)
	}

	apiMux := api.Mux(api.MuxConfig{
		Shutdown: shutdown,
		Log:      logger,

		DIDDoc: didDoc,
	})

	srv := httptest.NewServer(apiMux)

	teardown := func() {
		dbClient.CloseSession(ctx)
		srv.Close()
		lgDefer()
		cancel()
	}

	return srv, teardown
}
