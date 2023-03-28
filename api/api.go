package api

import (
	"net/http"
	"os"

	middleware "github.com/obada-foundation/registry/api/middleware/v1"
	"github.com/obada-foundation/registry/api/v1"
	"github.com/obada-foundation/registry/services/diddoc"
	"github.com/obada-foundation/registry/system/web"
	"go.uber.org/zap"
)

// MuxConfig defines the dependencies for the APIMux
type MuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger

	// Services
	DIDDoc diddoc.DIDDoc
}

// Mux constructs a http.Handler with all application routes defined.
func Mux(cfg MuxConfig) http.Handler {
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(),
	)

	v1.Routes(app, v1.Config{
		Log: cfg.Log,

		// Services
		DIDDoc: cfg.DIDDoc,
	})

	return app
}
