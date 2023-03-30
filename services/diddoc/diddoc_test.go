// nolint
package diddoc_test

import (
	"context"
	"testing"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/db"
	"github.com/obada-foundation/registry/testutil"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Service(t *testing.T) {
	c, err := testutil.StartDB()
	logger, deferFn := testutil.NewTestLoger()

	ctx := context.Background()

	dbClient, err := db.NewDBConnection(ctx, db.Connection{
		Host:   c.Host,
		Port:   c.Port,
		User:   "immudb",
		Pass:   "immudb",
		DBName: "defaultdb",
	})
	require.NoError(t, err)

	service := diddoc.NewService(dbClient, logger)

	t.Logf("Test \"Register\" function")
	{
		t.Logf("\tTest not supported DID methods")
		{
			notSupportedDIDs := []string{
				`{}`,
				`{"did":"did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De"}`,
			}

			for _, DID := range notSupportedDIDs {
				err := service.Register(ctx, DID)
				require.ErrorIs(t, err, sdkdid.ErrNotSupportedDIDMethod)
			}
		}

		DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"

		t.Logf("\tTest DID registration")
		{
			err := service.Register(ctx, DID)
			require.NoError(t, err)
		}

		t.Logf("\tTest DIDDoc fetching")
		{
			DIDDoc, err := service.Get(ctx, DID)
			require.NoError(t, err)

			assert.Equal(t, DID, DIDDoc.ID)
		}
	}

	defer func() {
		deferFn()
	}()
}
