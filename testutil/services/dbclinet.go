package services

import (
	"context"
	"testing"

	"github.com/codenotary/immudb/pkg/client"
	"github.com/obada-foundation/registry/system/db"
	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/require"
)

// MakeDBClient heler func that creates immudb client
// nolint: gocritic // no needed named vars in return
func MakeDBClient(ctx context.Context, t *testing.T) (client.ImmuClient, func()) {
	c, err := testutil.StartDB()
	require.NoError(t, err)

	dbClient, err := db.NewDBConnection(ctx, db.Connection{
		Host:   c.Host,
		Port:   c.Port,
		User:   "immudb",
		Pass:   "immudb",
		DBName: "defaultdb",
	})
	require.NoError(t, err)

	return dbClient, func() {
		_ = dbClient.CloseSession(ctx)
		testutil.StopDB(c)
	}
}
