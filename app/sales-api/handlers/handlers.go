// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"net/http"
	"os"

	"github.com/ardanlabs/service/business/sys/metrics"
	"github.com/ardanlabs/service/business/web/mid"
	"github.com/ardanlabs/service/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Build    string
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	Metrics  *metrics.Metrics
}

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(cfg.Metrics), mid.Panics())

	// Register debug check endpoints.
	cg := checkGroup{
		build: cfg.Build,
		log:   cfg.Log,
	}
	app.Handle(http.MethodGet, "/debug/readiness", cg.readiness)
	app.Handle(http.MethodGet, "/debug/liveness", cg.liveness)
	app.Handle(http.MethodGet, "/testerror", cg.testerror)

	return app
}
