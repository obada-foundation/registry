package did

import (
	"context"
	"fmt"
	"net/http"

	"github.com/obada-foundation/registry/system/web"
	"github.com/obada-foundation/registry/types"
	"github.com/obada-foundation/sdkgo"
)

// Handlers contains methods for DID API
type Handlers struct {
}

// Register DID in the registry
func (h Handlers) Register(_ context.Context, _ http.ResponseWriter, r *http.Request) error {
	var registerDID types.RegisterDID

	if err := web.Decode(r, &registerDID); err != nil {
		return fmt.Errorf("unable to decode request data: %w", err)
	}

	return nil
}
