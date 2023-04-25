package api_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/obada-foundation/registry/api"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func startGRPCServer(t *testing.T) (*grpc.Server, *bufconn.Listener, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	listener := bufconn.Listen(1024 * 1024)

	logger, logDeferFn := testutil.NewTestLoger()

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

	srv, _ := api.NewGRPCServer(api.GRPCConfig{
		Log: logger,

		// Services
		DIDDocService: diddoc.NewService(dbClient, logger),
	})
	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	return srv, listener, func() {
		logDeferFn()
		_ = listener.Close()
		_ = dbClient.CloseSession(ctx)
		srv.Stop()

		cancel()
	}
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func Test_GRPCServer(t *testing.T) {
	t.Parallel()

	_, listener, teardown := startGRPCServer(t)

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
		conn.Close()
		teardown()
	}()
}
