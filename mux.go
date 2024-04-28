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
func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}
