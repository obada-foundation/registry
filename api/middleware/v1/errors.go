package v1

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/obada-foundation/registry/api/errors"
	"github.com/obada-foundation/registry/system/validate"
	"github.com/obada-foundation/registry/system/web"
	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				log.Errorw("ERROR", "trace_id", web.GetTraceID(ctx), "message", err)

				var er errors.ErrorResponse
				var status int
				switch {

				case validate.IsFieldErrors(err):
					fieldErrors := validate.GetFieldErrors(err)
					er = errors.ErrorResponse{
						Error:  "data validation error",
						Fields: fieldErrors.Fields(),
					}
					status = http.StatusBadRequest

				case errors.IsRequestError(err):
					reqErr := errors.GetRequestError(err)
					er = errors.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					sentry.CaptureException(err)

					er = errors.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				if errResp := web.Respond(ctx, w, er, status); errResp != nil {
					return errResp
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

			return nil
		}

		return h
	}

	return m
}
