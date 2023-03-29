package db

import (
	"context"

	immudb "github.com/codenotary/immudb/pkg/client"
)

// Connection represents immudb connection options
type Connection struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DBName string
}

// NewDBConnection creates a client instance for immudb and opens authenticated session
func NewDBConnection(ctx context.Context, conn Connection) (immudb.ImmuClient, error) {
	opts := immudb.DefaultOptions().
		WithAddress(conn.Host).
		WithPort(conn.Port)

	client := immudb.NewClient().WithOptions(opts)
	err := client.OpenSession(
		ctx,
		[]byte(conn.User),
		[]byte(conn.Pass),
		conn.DBName,
	)

	return client, err
}
