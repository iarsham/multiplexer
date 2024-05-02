package multiplexer

import "net/http"

// Middlewares represents a list of middleware functions.
type Middlewares []func(http.Handler) http.Handler

// Router represents a multiplexer that routes incoming HTTP requests.
type Router struct {
	mux  *http.ServeMux
	path string
	mws  Middlewares
}

// New creates a new instance of the Router struct.
// It initializes the mux field with a new instance of http.ServeMux.
// Returns a pointer to the newly created Router.
func New(mux *http.ServeMux) *Router {
	return &Router{
		mux: mux,
	}
}

// ServeHTTP implements the http.Handler interface.
// It calls the ServeHTTP method of the underlying http.ServeMux.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
