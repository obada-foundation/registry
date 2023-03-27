package tests

import (
	"testing"

	"github.com/obada-foundation/registry/testutil"
	"github.com/stretchr/testify/require"
)

func Test_Register(t *testing.T) {
	srv, teardown := testutil.NewIntegrationTest(t)
	defer teardown()

	t.Log("Register a new DID")
	{
		_, err := testutil.Post(t, srv.URL+"/api/v1/did/register", `{}`, nil)
		require.NoError(t, err)
	}
}
