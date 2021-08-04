// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"net/http"

	"github.com/ardanlabs/service/business/web/mid"
	"github.com/ardanlabs/service/foundation/web"
	"go.uber.org/zap"
)

// APIMux constructs an http.Handler with all application routes defined.
func APIMux(log *zap.SugaredLogger) *web.App {
	app := web.NewApp(mid.Logger(log))

	app.Handle(http.MethodGet, "/debug/readiness", readiness)
	app.Handle(http.MethodGet, "/debug/liveness", liveness)

	return app
}
