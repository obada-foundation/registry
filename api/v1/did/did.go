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

// Get DID document from the registry
func (h Handlers) Get(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	DID := web.Param(r, "did")

	DIDDoc, err := h.DIDDoc.Get(ctx, DID)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, DIDDoc, http.StatusOK)
}

// Register DID in the registry
func (h Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var registerDID types.RegisterDID

	if err := web.Decode(r, &registerDID); err != nil {
		return fmt.Errorf("unable to decode request data: %w", err)
	}

	if err := h.DIDDoc.Register(ctx, registerDID.DID); err != nil {
		if err != sdkdid.ErrNotSupportedDIDMethod {
			return fmt.Errorf("cannot create DID from string: %w", err)
		}

		return apierrors.NewRequestError(err, http.StatusBadRequest)
	}

	return web.RespondWithNoContent(ctx, w, http.StatusCreated)
}
