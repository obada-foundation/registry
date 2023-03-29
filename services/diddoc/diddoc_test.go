// nolint
package diddoc_test

import (
	"testing"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/testutil"
	sdkdid "github.com/obada-foundation/sdkgo/did"
	"github.com/stretchr/testify/require"
)

func Test_Service(t *testing.T) {
	logger, deferFn := testutil.NewTestLoger()

	service := diddoc.NewService(logger)

	t.Logf("Test \"Register\" function")
	{
		t.Logf("\tTest not supported DID methods")
		{
			notSupportedDIDs := []string{
				`{}`,
				`{"did":"did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De"}`,
			}

			for _, DID := range notSupportedDIDs {
				err := service.Register(DID)
				require.ErrorIs(t, err, sdkdid.ErrNotSupportedDIDMethod)
			}
		}

		t.Logf("\tTest DID registration")
		{
			err := service.Register("did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88")
			require.NoError(t, err)
		}
	}

	defer func() {
		deferFn()
	}()
}
