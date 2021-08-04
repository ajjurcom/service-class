// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// APIMux constructs an http.Handler with all application routes defined.
func APIMux() *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	mux.Handle(http.MethodGet, "/debug/readiness", readiness)
	mux.Handle(http.MethodGet, "/debug/liveness", liveness)

	return mux
}
