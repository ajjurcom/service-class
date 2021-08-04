// Package handlers contains the full set of handler functions and routes
// supported by the web api.
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// APIMux constructs an http.Handler with all application routes defined.
func APIMux() *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status string
		}{
			Status: "OK",
		}
		json.NewEncoder(w).Encode(status)
	}

	mux.Handle(http.MethodGet, "/test", h)

	return mux
}
