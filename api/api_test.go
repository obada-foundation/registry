package api_test

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/gogo/protobuf/proto"
	"github.com/obada-foundation/registry/api"
	pbdiddoc "github.com/obada-foundation/registry/api/pb/v1/diddoc"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"github.com/obada-foundation/registry/testutil"
	"github.com/obada-foundation/sdkgo/asset"
	"github.com/obada-foundation/sdkgo/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type apiTests struct {
	conn   *grpc.ClientConn
	diddoc pbdiddoc.DIDDocClient
}

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
		grpc.WithInsecure(),
	)

	require.NoError(t, err)

	tests := apiTests{
		conn:   conn,
		diddoc: pbdiddoc.NewDIDDocClient(conn),
	}

	pbdiddoc.NewDIDDocClient(conn)

	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
		conn.Close()
	}()

	t.Run("accountsRpcCalls", tests.accountsRpcCalls)
	t.Run("notSuportedDIDsRegister", tests.notSuportedDIDsRegister)
	t.Run("notRegisteredDIDs", tests.notRegisteredDIDs)
	t.Run("registerDID", tests.registerDID)
	t.Run("saveMetadata", tests.saveMetadata)
}

func (tests apiTests) saveMetadata(t *testing.T) {
	t.Log("Save metadata")

	ctx := context.Background()

	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb82"

	t.Log("\tRegister DID for testing metadata save")
	{
		tests.registerDIDWithNoErrs(t, DID)
	}

	t.Log("\tSaving metadata with an empty signature")
	{
		_, err := tests.diddoc.SaveMetadata(ctx, &pbdiddoc.SaveMetadataRequest{
			Signature: []byte(""),
			Data: &pbdiddoc.SaveMetadataRequest_Data{
				Did: DID,
				Objects: append(make([]*pbdiddoc.Object, 1), &pbdiddoc.Object{
					Url: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
					Metadata: map[string]string{
						"type": string(asset.MainImage),
					},
					HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
				},
				),
			},
		})

		permissionDenied(t, err)
	}

	t.Log("\tSaving metadata without verification method")
	{
		_, err := tests.diddoc.SaveMetadata(ctx, &pbdiddoc.SaveMetadataRequest{
			Signature: []byte("some signature"),
			Data: &pbdiddoc.SaveMetadataRequest_Data{
				Did: DID,
				Objects: append(make([]*pbdiddoc.Object, 1), &pbdiddoc.Object{
					Url: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
					Metadata: map[string]string{
						"type": string(asset.MainImage),
					},
					HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
				},
				),
			},
		})

		permissionDenied(t, err)
	}

	t.Log("\tSaving signed metadata")
	{
		privKey := secp256k1.GenPrivKey()
		pubKey := privKey.PubKey()
		base58PubKey := base58.Encode(pubKey.Bytes())

		newDID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb89"

		regMsg := &pbdiddoc.RegisterRequest{
			Did:                newDID,
			VerificationMethod: make([]*pbdiddoc.VerificationMethod, 0),
			Authentication: []string{
				fmt.Sprintf("%s#keys-1", newDID),
			},
		}

		regMsg.VerificationMethod = append(regMsg.VerificationMethod, &pbdiddoc.VerificationMethod{
			Id:              fmt.Sprintf("%s#keys-1", newDID),
			PublicKeyBase58: base58PubKey,
		})

		_, err := tests.diddoc.Register(context.Background(), regMsg)
		require.NoError(t, err)

		data := &pbdiddoc.SaveMetadataRequest_Data{
			Did:     newDID,
			Objects: make([]*pbdiddoc.Object, 0),
		}

		data.Objects = append(data.Objects, &pbdiddoc.Object{
			Url: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
			Metadata: map[string]string{
				"type":  string(asset.MainImage),
				"color": "red",
			},
			HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
		})

		dataBytes, err := proto.Marshal(data)
		require.NoError(t, err)

		signature, err := privKey.Sign(dataBytes)
		require.NoError(t, err)

		_, er := tests.diddoc.SaveMetadata(ctx, &pbdiddoc.SaveMetadataRequest{
			Signature: signature,
			Data:      data,
		})

		require.NoError(t, er)
	}
}

func (tests apiTests) registerDID(t *testing.T) {
	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"
	t.Log("Register new DID")
	{
		tests.registerDIDWithNoErrs(t, DID)
	}
	t.Log("Get newly registered DID")
	{
		DIDDoc, err := tests.diddoc.Get(context.Background(), &pbdiddoc.GetRequest{
			Did: DID,
		})
		require.NoError(t, err)
		assert.Equal(t, DID, DIDDoc.GetDocument().GetId())
	}
	t.Log("Cannot register DID that already been registered")
	{
		_, err := tests.diddoc.Register(context.Background(), &pbdiddoc.RegisterRequest{
			Did: DID,
		})

		er, ok := status.FromError(err)

		assert.True(t, ok, "error is not a grpc error")
		assert.Equal(t, diddoc.ErrDIDAlreadyRegistered.Error(), er.Message())
		assert.Equal(t, codes.AlreadyExists, er.Code())
	}
}

func (tests apiTests) notRegisteredDIDs(t *testing.T) {
	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb81"

	ctx := context.Background()

	t.Log("\tQuery not registered DID")
	{
		_, err := tests.diddoc.Get(ctx, &pbdiddoc.GetRequest{Did: DID})
		notFoundDID(t, err)
	}

	t.Log("\tQuery not registered DID history")
	{
		_, err := tests.diddoc.GetMetadataHistory(ctx, &pbdiddoc.GetMetadataHistoryRequest{Did: DID})
		notFoundDID(t, err)
	}

	t.Log("\tSave metadata for not registered DID")
	{
		_, err := tests.diddoc.SaveMetadata(ctx, &pbdiddoc.SaveMetadataRequest{
			Signature: []byte("signature"),
			Data: &pbdiddoc.SaveMetadataRequest_Data{
				Did: DID,
			},
		})
		notFoundDID(t, err)
	}
}

func (tests apiTests) notSuportedDIDsRegister(t *testing.T) {
	t.Log("Register not supported DIDs")

	msgs := []*pbdiddoc.RegisterRequest{
		{},
		{
			Did: "did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De",
		},
	}

	for _, msg := range msgs {
		_, err := tests.diddoc.Register(context.Background(), msg)
		er, ok := status.FromError(err)

		assert.True(t, ok, "error is not a grpc error")
		assert.Equal(t, "given DID method is not supported", er.Message(), msg.GetDid())
	}
}

func (tests apiTests) accountsRpcCalls(t *testing.T) {
	t.Log("test5")

}

func (tests apiTests) registerDIDWithNoErrs(t *testing.T, DID string) {
	_, err := tests.diddoc.Register(context.Background(), &pbdiddoc.RegisterRequest{
		Did: DID,
		VerificationMethod: append(make([]*pbdiddoc.VerificationMethod, 1), &pbdiddoc.VerificationMethod{
			Id: fmt.Sprintf("%s#keys-1", DID),
		}),
		Authentication: []string{
			fmt.Sprintf("%s#keys-1", DID),
		},
	})
	require.NoError(t, err)
}

func notFoundDID(t *testing.T, err error) {
	er, ok := status.FromError(err)

	assert.True(t, ok, "error is not a grpc error")
	assert.Equal(t, "DID is not registered", er.Message())
	assert.Equal(t, codes.NotFound, er.Code())
}

func permissionDenied(t *testing.T, err error) {
	er, ok := status.FromError(err)

	assert.True(t, ok, "error is not a grpc error")
	assert.Equal(t, "unauthorized", er.Message())
	assert.Equal(t, codes.PermissionDenied, er.Code())
}
