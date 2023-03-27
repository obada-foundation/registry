package did

import (
	"context"
	"net/http"
)

type Handlers struct {
}

// Register DID in the registry
func (h Handlers) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}
