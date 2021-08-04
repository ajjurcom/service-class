// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

// APIMux constructs an http.Handler with all application routes defined.
func APIMux() *web.App {
	app := web.NewApp()

	app.Handle(http.MethodGet, "/debug/readiness", readiness)
	app.Handle(http.MethodGet, "/debug/liveness", liveness)

	return app
}
