package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	apierrors "github.com/obada-foundation/registry/api/errors"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/web"
	"github.com/obada-foundation/registry/types"
)

// Authenticate middleware checks signature and verifies that it is valid for the given DID
func Authenticate(svc diddoc.DIDDoc) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if r.Method == http.MethodPost {
				var req map[string]interface{}

				DID := web.Param(r, "did")

				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return err
				}
				_ = r.Body.Close()
				r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

				buf := bytes.NewBuffer(body)
				if er := json.NewDecoder(buf).Decode(&req); er != nil {
					return fmt.Errorf("unable to decode request data: %w", er)
				}

				signature, ok := req["signature"]
				if !ok || signature == "" {
					return apierrors.NewRequestError(types.ErrUnauthorizedNoSignature, http.StatusUnauthorized)
				}

				DIDDoc, err := svc.Get(ctx, DID)
				if err != nil {
					return err
				}

				for _, authKey := range DIDDoc.Authentication {
					for _, method := range DIDDoc.VerificationMethod {
						if method.ID == authKey {
							return handler(ctx, w, r)
						}
					}
				}
			}

			return apierrors.NewRequestError(types.ErrUnauthorized, http.StatusUnauthorized)
		}

		return h
	}

	return m
}
