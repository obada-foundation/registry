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
	srv, teardown := testutil.NewIntegrationTest(t)
	defer teardown()

	t.Log("Register not supported DIDs")
	{
		notSupportedDIDs := []string{
			`{}`}

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
}
