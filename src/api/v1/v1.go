package v1

import (
	"net/http"

	"github.com/obada-foundation/registry/api/v1/did"
	"github.com/obada-foundation/registry/system/web"
	"go.uber.org/zap"
)

// Config contains all the mandatory systems required by route handlers.
type Config struct {
	Log *zap.SugaredLogger
}

// Routes binds all the version 1 routes.
func Routes(app *web.App, _ Config) {
	const version = "api/v1"

	didGroup := did.Handlers{}

	app.Handle(http.MethodPost, version, "/did/register", didGroup.Register)
}
