package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/require"
)

func Test_Register(t *testing.T) {
	srv, teardown := testutil.NewIntegrationTest(t, nil)
	defer teardown()

	t.Log("Register not supported DIDs")
	{
		notSupportedDIDs := []string{
			`{}`,
			`{"did":"did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De"}`,
		}

		for _, req := range notSupportedDIDs {
			resp, err := testutil.Post(t, srv.URL+"/api/v1/register", req, nil)
			require.NoError(t, err)

			b, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, http.StatusBadRequest, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			c := JSON{}
			err = json.Unmarshal(b, &c)
			require.NoError(t, err)
		}
	}

	t.Log("Register new DID")
	{
		req := `{"did": "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"}`

		resp, err := testutil.Post(t, srv.URL+"/api/v1/register", req, nil)
		require.NoError(t, err)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.NoError(t, resp.Body.Close())

		c := JSON{}
		err = json.Unmarshal(b, &c)
		require.NoError(t, err)
	}
}
