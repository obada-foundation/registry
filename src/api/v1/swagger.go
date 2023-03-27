package v1

import (
	"context"
	"net/http"

	"github.com/obada-foundation/registry/system/web"
)

const swaggerVersion = "swagger/v1"

// nolint
func Swagger(app *web.App) {
	app.Handle(http.MethodGet, swaggerVersion, "/", RedirectToIndex)
	app.Handle(http.MethodGet, swaggerVersion, "/index.html", Index)
	app.Handle(http.MethodGet, swaggerVersion, "/doc.json", RedirectToIndex)
}

// nolint
func RedirectToIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, swaggerVersion+"/index.html", http.StatusMovedPermanently)

	return nil
}

// nolint
func Index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	http.Redirect(w, r, swaggerVersion+"/index.html", http.StatusMovedPermanently)

	return nil
}
