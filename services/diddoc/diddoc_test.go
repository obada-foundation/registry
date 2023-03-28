package diddoc_test

import (
	"testing"

	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/testutil"
)

func Test_Service(t *testing.T) {
	logger, _ := testutil.NewTestLoger()

	service := diddoc.NewService(logger)

	defer func() {
		logger.Sync()
	}()

	t.Logf("Test \"Register\" function")
	{
		t.Logf("\tTest not supported DID methods")
		{

		}

		t.Logf("\tTest DID registration")
		{
			service.Register("did:obada:64925be84b586363670c1f7e5ada86a37904e590d1f6570d834436331dd3eb88)")
		}
	}
}
