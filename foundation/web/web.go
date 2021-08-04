// Package web contains a small web framework extension.
package web

import (
	"context"
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
	mw []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		mw:         mw,
	}
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	handler = wrapMiddleware(mw, handler)

	// Add the application's general middleware to the handler chain.
	handler = wrapMiddleware(a.mw, handler)

	// Most Outer Layer
	h := func(w http.ResponseWriter, r *http.Request) {

		// INJECT HERE

		// Most Inner Layer
		if err := handler(r.Context(), w, r); err != nil {
			return
		}

		// INJECT HERE
	}

	a.ContextMux.Handle(method, path, h)
}
