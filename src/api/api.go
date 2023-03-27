package api

import (
	"net/http"
	"os"

	middleware "github.com/obada-foundation/registry/api/middleware/v1"
	"github.com/obada-foundation/registry/api/v1"
	"github.com/obada-foundation/registry/system/web"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) http.Handler {
	app := web.NewApp(
		cfg.Shutdown,
		middleware.Logger(cfg.Log),
		middleware.Errors(cfg.Log),
		middleware.Metrics(),
		middleware.Panics(),
	)

	v1.Routes(app, v1.Config{
		Log: cfg.Log,
	})

	return app
}
