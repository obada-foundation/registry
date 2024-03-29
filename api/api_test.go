package api_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/codenotary/immudb/pkg/api/schema"
	"github.com/obada-foundation/registry/api"
	pbacc "github.com/obada-foundation/registry/api/pb/v1/account"
	pbdiddoc "github.com/obada-foundation/registry/api/pb/v1/diddoc"
	"github.com/obada-foundation/registry/services/account"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type apiTests struct {
	conn    *grpc.ClientConn
	diddoc  pbdiddoc.DIDDocClient
	account pbacc.AccountClient
}

func startGRPCServer(t *testing.T) (*grpc.Server, *bufconn.Listener, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	c, err := testutil.StartDB()
	require.NoError(t, err, "cabbot start db")

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

	_, err = dbClient.ServerInfo(ctx, &schema.ServerInfoRequest{})
	require.NoError(t, err, "immudb is not healthy")

	srv, _ := api.NewGRPCServer(api.GRPCConfig{
		Log: logger,

		// Services
		DIDDocService:  diddoc.NewService(dbClient, logger),
		AccountService: account.NewService(dbClient, logger),
	})
	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	return srv, listener, func() {
		logDeferFn()
		_ = dbClient.CloseSession(ctx)
		testutil.StopDB(c)

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

	srv, listener, teardown := startGRPCServer(t)
	t.Cleanup(func() {
		teardown()
		srv.Stop()
		listener.Close()
	})

	conn, err := grpc.DialContext(
		context.Background(),
		"",
		grpc.WithContextDialer(getBufDialer(listener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)

	tests := apiTests{
		conn:    conn,
		diddoc:  pbdiddoc.NewDIDDocClient(conn),
		account: pbacc.NewAccountClient(conn),
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
		conn.Close()
	}()

	// DID Docs test
	t.Run("registerNotSuportedDIDs", tests.registerNotSuportedDIDs)
	t.Run("notRegisteredDIDs", tests.notRegisteredDIDs)
	t.Run("registerDID", tests.registerDID)
	t.Run("saveMetadata", tests.saveMetadata)
	t.Run("SaveVerificationMethods", tests.SaveVerificationMethods)

	// Accounts test
	t.Run("registerAccount", tests.registerAccount)
}

func unauthenticated(t *testing.T, err error) {
	checkErr(t, err, "unauthorized", codes.Unauthenticated)
}

func emptySignature(t *testing.T, err error) {
	checkErr(t, err, "empty signature", codes.InvalidArgument)
}

func emptyAuthentificationID(t *testing.T, err error) {
	checkErr(t, err, "empty autherntication id", codes.InvalidArgument)
}

func unknownVerificationMethod(t *testing.T, err error) {
	t.Log(err)
	checkErr(t, err, "unknown verification method", codes.Unknown)
}

func notFoundDID(t *testing.T, err error) {
	checkErr(t, err, "DID is not registered", codes.NotFound)
}

func checkErr(t *testing.T, err error, errMsg string, code codes.Code) {
	er, ok := status.FromError(err)

	assert.True(t, ok, "error is not a grpc error")
	assert.Contains(t, errMsg, er.Message())
	assert.Equal(t, code, er.Code())
}
