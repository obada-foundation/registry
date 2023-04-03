package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DIDTests struct {
	srv *httptest.Server
}

func Test_Register(t *testing.T) {
	t.Parallel()

	srv, teardown := testutil.NewIntegrationTest(t, c, nil)

	tests := DIDTests{
		srv: srv,
	}

	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
		teardown()
	}()

	t.Run("registerDID", tests.registerDID)
	t.Run("notSuportedDIDsRegister", tests.notSuportedDIDsRegister)
	t.Run("registeredMiddleware", tests.registeredMiddleware)
	t.Run("saveMetadata", tests.saveMetadata)
}

func (dt *DIDTests) saveMetadata(t *testing.T) {
	t.Log("Save metadata")
	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"
	_, err := testutil.Post(t, dt.srv.URL+"/api/v1.0/"+DID+"/metadata", `{}`, nil)
	require.NoError(t, err)
}

func (dt *DIDTests) registerDID(t *testing.T) {
	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88"
	t.Log("Register new DID")
	{
		req := `{"did": "` + DID + `"}`

		resp, err := testutil.Post(t, dt.srv.URL+"/api/v1.0/register", req, nil)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		require.Equal(t, http.StatusCreated, resp.StatusCode)
	}
	t.Log("Get newly registered DID")
	{
		resp, err := testutil.Get(t, dt.srv.URL+"/api/v1.0/"+DID, nil)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		b, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.NoError(t, resp.Body.Close())

		c := JSON{}
		err = json.Unmarshal(b, &c)
		require.NoError(t, err)

		assert.Equal(t, DID, c["id"])
	}
}

func (dt *DIDTests) notSuportedDIDsRegister(t *testing.T) {
	t.Log("Register not supported DIDs")

	notSupportedDIDs := []string{
		`{}`,
		`{"did":"did:bnb:1f4B9d871fed2dEcb2670A80237F7253DB5766De"}`,
	}

	for _, req := range notSupportedDIDs {
		resp, err := testutil.Post(t, dt.srv.URL+"/api/v1.0/register", req, nil)
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

func (dt *DIDTests) registeredMiddleware(t *testing.T) {
	t.Log("Test routes that require registered DID")

	DID := "did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb81"

	t.Log("\tQuery not registered DID")
	{
		resp, err := testutil.Get(t, dt.srv.URL+"/api/v1.0/"+DID, nil)
		require.NoError(t, err)
		DID404(t, resp)
	}

	t.Log("\tQuery not registered DID history")
	{
		resp, err := testutil.Get(t, dt.srv.URL+"/api/v1.0/"+DID+"/metadata-history", nil)
		require.NoError(t, err)
		DID404(t, resp)
	}

	t.Log("\tSave metadata for not registered DID")
	{
		resp, err := testutil.Post(t, dt.srv.URL+"/api/v1.0/"+DID+"/metadata", `{}`, nil)
		require.NoError(t, err)
		DID404(t, resp)
	}

}

func DID404(t *testing.T, resp *http.Response) {
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	b, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())

	c := JSON{}
	err = json.Unmarshal(b, &c)
	require.NoError(t, err)

	assert.Equal(t, diddoc.ErrDIDNotRegitered.Error(), c["error"])
}
