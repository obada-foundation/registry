package v1

import (
	"net/http"

	"github.com/obada-foundation/registry/api/v1/did"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/web"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by route handlers.
type Config struct {
	Log *zap.SugaredLogger

	// Services
	DIDDoc diddoc.DIDDoc
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, cfg Config) {
	const version = "api/v1.0"

	didGroup := did.Handlers{
		DIDDoc: cfg.DIDDoc,
	}

	app.Handle(http.MethodPost, version, "/register", didGroup.Register)
	app.Handle(http.MethodGet, version, "/:did", didGroup.Get)
}
