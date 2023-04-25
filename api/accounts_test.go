package api_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types"
	pbacc "github.com/obada-foundation/registry/api/pb/v1/account"
	"github.com/obada-foundation/registry/services/account"
	"github.com/obada-foundation/sdkgo/base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (tests apiTests) registerAccount(t *testing.T) {
	ctx := context.Background()

	type tc struct {
		msg *pbacc.RegisterAccountRequest
		err error
	}

	tcs := []tc{
		{
			msg: &pbacc.RegisterAccountRequest{},
			err: account.ErrPublicKeyIsEmpty,
		},
		{
			msg: &pbacc.RegisterAccountRequest{
				Pubkey: "test",
			},
			err: account.ErrInvalidPublicKey,
		},
		{
			msg: &pbacc.RegisterAccountRequest{
				Pubkey: base58.Encode(secp256k1.GenPrivKey().PubKey().Bytes()),
			},
			err: nil,
		},
	}

	for _, tc := range tcs {
		_, err := tests.account.RegisterAccount(ctx, tc.msg)

		if tc.err != nil {
			er, ok := status.FromError(err)
			assert.True(t, ok, "error is not a grpc error")
			assert.Equal(t, tc.err.Error(), er.Message())

			continue
		}

		require.NoError(t, err)

		pk := &secp256k1.PubKey{
			Key: base58.Decode(tc.msg.GetPubkey()),
		}

		resp, err := tests.account.GetPublicKey(ctx, &pbacc.GetPublicKeyRequest{
			Address: types.AccAddress(pk.Address()).String(),
		})
		require.NoError(t, err)

		assert.Equal(t, tc.msg.GetPubkey(), resp.GetPubkey())
	}

	_, err := tests.account.GetPublicKey(ctx, &pbacc.GetPublicKeyRequest{})
	assert.True(t, strings.Contains(err.Error(), account.ErrInvalidAddress.Error()))

	_, err = tests.account.GetPublicKey(ctx, &pbacc.GetPublicKeyRequest{Address: "cosmos13j5sl00cg820r3twh8g4xm4znj0ery45674p3g"})
	er, ok := status.FromError(err)
	assert.True(t, ok, "error is not a grpc error")
	assert.Equal(t, account.ErrPubKeyNotRegistered.Error(), er.Message())
	assert.Equal(t, codes.NotFound, er.Code())
}
