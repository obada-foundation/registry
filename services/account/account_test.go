package account_test

import (
	"context"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/obada-foundation/common/testutil"
	"github.com/obada-foundation/registry/services/account"
	"github.com/obada-foundation/registry/testutil/services"
	"github.com/obada-foundation/sdkgo/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Service(t *testing.T) {
	ctx := context.Background()

	dbClient, deferFn := services.MakeDBClient(t, ctx)
	defer deferFn()

	logger, deferFn := testutil.NewTestLoger()
	defer deferFn()

	service := account.NewService(dbClient, logger)

	type regTestCase struct {
		pubKey []byte
		error  error
	}

	type getTestCase struct {
		addr string
		err  error
	}

	regTestCases := []regTestCase{
		{
			pubKey: []byte(""),
			error:  account.ErrPublicKeyIsEmpty,
		},
		{
			pubKey: []byte("safafdfsdf"),
			error:  account.ErrInvalidPublicKey,
		},
		{
			pubKey: []byte("mwZ6J1yji1Q1mJEmaCiqunkdtJdGYWNTqCWqFwAR37Bm5231jYPAYzm1vmjJGPC5x1tbdDbBx2wDRe6bzGaX8rdaCu"),
			error:  account.ErrInvalidPublicKey,
		},
		{
			pubKey: []byte("safafdfsddfdsgsdgdsgsdgdsgdasfsdsgfhdfhgkfkhgfjgjgsjfgjfdjhhfjhfhf"),
			error:  account.ErrInvalidPublicKey,
		},
		{
			pubKey: secp256k1.GenPrivKey().PubKey().Bytes(),
			error:  nil,
		},
	}

	for _, tc := range regTestCases {
		b58PubKey := base58.Encode(tc.pubKey)
		pk := &secp256k1.PubKey{
			Key: tc.pubKey,
		}

		err := service.Register(ctx, b58PubKey)

		if tc.error != nil {
			require.ErrorIs(t, err, tc.error)
			continue
		}

		require.NoError(t, err)

		addr := types.AccAddress(pk.Address()).String()

		pubKey, err := service.GetPublicKey(ctx, addr)
		require.NoError(t, err)

		assert.Equal(t, b58PubKey, pubKey)
	}

	getTestCases := []getTestCase{
		{
			addr: "",
			err:  account.ErrInvalidAddress,
		},
		{
			addr: "test",
			err:  account.ErrInvalidAddress,
		},
		{
			addr: "cosmos1hjnfd4v8lenkrfhwapjvu20dqknh0dngceh5um",
			err:  account.ErrPubKeyNotRegistered,
		},
	}

	for _, tc := range getTestCases {
		_, err := service.GetPublicKey(ctx, tc.addr)
		require.ErrorIs(t, err, tc.err)
	}
}
