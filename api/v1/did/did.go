package did

import (
	"context"
	"fmt"
	"net/http"

	apierrors "github.com/obada-foundation/registry/api/errors"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/web"
	"github.com/obada-foundation/registry/types"
	sdkdid "github.com/obada-foundation/sdkgo/did"
)

// Handlers contains methods for DID API
type Handlers struct {
	DIDDoc diddoc.DIDDoc
}

// Register DID in the registry
func (h Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var registerDID types.RegisterDID

	if err := web.Decode(r, &registerDID); err != nil {
		return fmt.Errorf("unable to decode request data: %w", err)
	}

	DID, err := sdkdid.FromString(registerDID.DID, nil)
	if err != nil {
		if err != sdkdid.ErrNotSupportedDIDMethod {
			return fmt.Errorf("cannot create DID from string: %w", err)
		}

		return apierrors.NewRequestError(err, http.StatusBadRequest)
	}

	if err := h.DIDDoc.Register(*DID); err != nil {
		return err
	}

	return web.Respond(ctx, w, registerDID, http.StatusOK)
}
