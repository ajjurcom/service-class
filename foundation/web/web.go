// Package web contains a small web framework extension.
package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

// A Handler is a type that handles an http request within our own little mini
// framework.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp() *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, path string, handler Handler) {

	// Most Outer Layer
	h := func(w http.ResponseWriter, r *http.Request) {

		// CODE HERE
		fmt.Printf("{\"status\": %q, \"url\": %q}\n", "started", r.URL.String())

		// Most Inner Layer
		if err := handler(r.Context(), w, r); err != nil {

			// NOT SURE WHAT TO DO - ERROR HANDLING
			return
		}

		fmt.Printf("{\"status\": %q, \"url\": %q}\n", "end", r.URL.String())

		// CODE HERE
		// LOGGING HERE
	}

	a.ContextMux.Handle(method, path, h)
}
