// nolint
package diddoc_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/testutil"
	"github.com/obada-foundation/registry/testutil/services"
	"github.com/obada-foundation/registry/types"
	"github.com/obada-foundation/sdkgo/asset"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	did string
	vm  []types.VerificationMethod
	a   []string
	err error
}

func Test_Service(t *testing.T) {
	logger, deferFn := testutil.NewTestLoger()
	defer deferFn()

	ctx := context.Background()

	dbClient, deferFn := services.MakeDBClient(ctx, t)
	defer deferFn()

	validDID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"

	service := diddoc.NewService(dbClient, logger)
	t.Logf("Test \"Register\" function")
	{

		tcs := []testCase{
			{
				did: "",
				err: sdkdid.ErrNotSupportedDIDMethod,
			},
			{
				did: "did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De",
				err: sdkdid.ErrNotSupportedDIDMethod,
			},
			{
				did: validDID,
				vm: []types.VerificationMethod{
					{
						ID:              fmt.Sprintf("%s#keys-1", validDID),
						Type:            types.Ed25519VerificationKey2018JSONLD,
						Controller:      validDID,
						PublicKeyBase58: "",
					},
				},
				a:   []string{fmt.Sprintf("%s#keys-1", validDID)},
				err: nil,
			},
			{
				did: validDID,
				err: diddoc.ErrDIDAlreadyRegistered,
			},
		}

		for _, tc := range tcs {
			err := service.Register(ctx, tc.did, tc.vm, tc.a)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
				continue
			}

			require.NoError(t, err)

			DIDDoc, err := service.Get(ctx, tc.did)
			require.NoError(t, err)

			require.NoError(t, err)
			assert.Equal(t, tc.did, DIDDoc.ID)
			assert.Equal(t, types.DIDSchemaJSONLD, DIDDoc.Context[0])
			assert.Equal(t, types.Ed25519VerificationKey2018JSONLD, DIDDoc.Context[1])
			assert.Equal(t, tc.vm, DIDDoc.VerificationMethod)
			assert.Equal(t, tc.a, DIDDoc.Authentication)
			assert.Equal(t, 0, len(DIDDoc.Metadata.Objects))

			j, err := json.Marshal(DIDDoc)
			require.NoError(t, err)

			t.Logf("\t DIDDoc: \n%s", string(j))

		}
		t.Log("Test to query not registered DID")
		_, err := service.Get(ctx, "64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb84")
		require.ErrorIs(t, err, diddoc.ErrDIDNotRegistered)
	}

	t.Logf("Test \"SaveMetadata\" function")
	{
		t.Logf("\t Add first object to the metadata. Version: 1")
		{
			err := service.SaveMetadata(
				ctx,
				validDID,
				[]asset.Object{
					{
						URL: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
						Metadata: map[string]string{
							"type": string(asset.MainImage),
						},
						HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
					},
				},
			)
			require.NoError(t, err)

			DIDDoc, err := service.Get(ctx, validDID)
			require.NoError(t, err)
			assert.Equal(t, 1, len(DIDDoc.MetadataHistory))
			assert.Equal(t, 1, len(DIDDoc.MetadataHistory[1].Objects))
			assert.Equal(t, "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE", DIDDoc.MetadataHistory[1].Objects[0].HashUnencryptedObject)

			assert.Equal(t, 1, DIDDoc.Metadata.VersionID)
			assert.Equal(t, DIDDoc.MetadataHistory[1].RootHash, DIDDoc.Metadata.RootHash)
			assert.Equal(t, DIDDoc.MetadataHistory[1].VersionHash, DIDDoc.Metadata.VersionHash)
			assert.Equal(t, DIDDoc.MetadataHistory[1].Objects, DIDDoc.Metadata.Objects)

			j, err := json.Marshal(DIDDoc.MetadataHistory)
			require.NoError(t, err)

			t.Logf("\t\t DIDDoc: \n%s", string(j))
		}

		t.Logf("\t Add second object to the metadata. Version: 2")
		{
			err := service.SaveMetadata(
				ctx,
				validDID,
				[]asset.Object{
					{
						URL: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
						Metadata: map[string]string{
							"type": string(asset.MainImage),
						},
						HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
					},
					{
						URL: "https://ipfs.io/ipfs/QmUyARmq5RUJk5zt7KUeaMLYB8SQbKHp3Gdqy5WSxRtPNa",
						Metadata: map[string]string{
							"type": string(asset.MainImage),
						},
						HashUnencryptedObject: "QmUyARmq5RUJk5zt7KUeaMLYB8SQbKHp3Gdqy5WSxRtPNa",
					},
				},
			)
			require.NoError(t, err)

			DIDDoc, err := service.Get(ctx, validDID)
			require.NoError(t, err)
			assert.Equal(t, 2, len(DIDDoc.MetadataHistory))
			assert.Equal(t, 1, len(DIDDoc.MetadataHistory[1].Objects))
			assert.Equal(t, 2, len(DIDDoc.MetadataHistory[2].Objects))
			assert.Equal(t, "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE", DIDDoc.MetadataHistory[1].Objects[0].HashUnencryptedObject)
			assert.Equal(t, "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE", DIDDoc.MetadataHistory[2].Objects[0].HashUnencryptedObject)
			assert.Equal(t, "QmUyARmq5RUJk5zt7KUeaMLYB8SQbKHp3Gdqy5WSxRtPNa", DIDDoc.MetadataHistory[2].Objects[1].HashUnencryptedObject)

			assert.Equal(t, 2, DIDDoc.Metadata.VersionID)
			assert.Equal(t, DIDDoc.MetadataHistory[2].RootHash, DIDDoc.Metadata.RootHash)
			assert.Equal(t, DIDDoc.MetadataHistory[2].VersionHash, DIDDoc.Metadata.VersionHash)
			assert.Equal(t, DIDDoc.MetadataHistory[2].Objects, DIDDoc.Metadata.Objects)

			j, err := json.Marshal(DIDDoc.MetadataHistory)
			require.NoError(t, err)
			t.Logf("\t\t DIDDoc: \n%s", string(j))
		}

		t.Logf("\t Remove a second object from the metadata. Version: 3")
		{
			err := service.SaveMetadata(
				ctx,
				validDID,
				[]asset.Object{
					{
						URL: "https://ipfs.io/ipfs/QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
						Metadata: map[string]string{
							"type": string(asset.MainImage),
						},
						HashUnencryptedObject: "QmQqzMTavQgT4f4T5v6PWBp7XNKtoPmC9jvn12WPT3gkSE",
					},
				},
			)
			require.NoError(t, err)

			DIDDoc, err := service.Get(ctx, validDID)
			require.NoError(t, err)
			assert.Equal(t, 3, DIDDoc.Metadata.VersionID)
			assert.Equal(t, DIDDoc.MetadataHistory[3].RootHash, DIDDoc.Metadata.RootHash)
			assert.Equal(t, DIDDoc.MetadataHistory[1].VersionHash, DIDDoc.Metadata.VersionHash)
			assert.Equal(t, DIDDoc.MetadataHistory[1].Objects, DIDDoc.Metadata.Objects)

			j, err := json.Marshal(DIDDoc.MetadataHistory)
			require.NoError(t, err)
			t.Logf("\t DIDDoc: \n%s", string(j))
		}

		t.Logf("\t Get metadata history")
		{
			_, err := service.GetMetadataHistory(ctx, validDID)
			require.NoError(t, err)
		}

		t.Logf("Test \"GetVerificationKeyByAuthID\"")
		{
			_, err := service.GetVerificationKeyByAuthID(ctx, validDID, fmt.Sprintf("%s#keys-1", validDID))
			require.NoError(t, err)

			invalidDID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb81"
			_, err = service.GetVerificationKeyByAuthID(ctx, invalidDID, fmt.Sprintf("%s#keys-1", invalidDID))
			require.ErrorIs(t, err, diddoc.ErrDIDNotRegistered)

			_, err = service.GetVerificationKeyByAuthID(ctx, validDID, fmt.Sprintf("%s#keys-2", validDID))
			require.ErrorIs(t, err, diddoc.ErrVerificationKeyNotFound)
		}

		t.Logf("Test \"SaveVerificationMethods\"")
		{
			vms := append(make([]types.VerificationMethod, 0, 1), types.VerificationMethod{
				ID:              fmt.Sprintf("%s#keys-1", validDID),
				Type:            types.Ed25519VerificationKey2018JSONLD,
				Controller:      validDID,
				PublicKeyBase58: "",
			})

			a := []string{fmt.Sprintf("%s#keys-1", validDID)}

			err := service.SaveVerificationMethods(ctx, validDID, vms, a)
			require.NoError(t, err)
		}
	}
}
