package v1

import (
	"context"
	"errors"
	"net/http"

	apierrors "github.com/obada-foundation/registry/api/errors"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/web"
)

func Registered(svc diddoc.DIDDoc) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			DID := web.Param(r, "did")

			if _, err := svc.Get(ctx, DID); err != nil {
				if errors.Is(err, diddoc.ErrDIDNotRegitered) {
					return apierrors.NewRequestError(err, http.StatusNotFound)
				}

				return err
			}

			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
