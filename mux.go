package multiplexer

import (
	"net/http"
	"regexp"
)

// reGo122 is a compiled regular expression used to match and extract two groups from a string pattern.
var reGo122 = regexp.MustCompile(`^(\S*)\s+(.*)$`)

// Router represents a multiplexer that routes incoming HTTP requests.
type Router struct {
	mux  *http.ServeMux
	path string
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

// Registers a route with the router, prepending base path for named capture groups.
// pattern (string): URL pattern with optional named capture groups (e.g., `/users/:id`).
// handler (http.HandlerFunc): Function to handle requests matching the pattern.
func (r *Router) register(pattern string, handler http.HandlerFunc) {
	match := reGo122.FindStringSubmatch(pattern)
	if len(match) > 2 {
		pattern = match[1] + " " + r.path + match[2]
	} else {
		pattern = r.path + pattern
	}
	r.mux.HandleFunc(pattern, handler)
}

// HandleFunc adds a new route with the given pattern and handler function.
func (r *Router) HandleFunc(pattern string, handler http.HandlerFunc) {
	r.register(pattern, handler)
}

// Handle adds a new route with the given pattern and handler.
func (r *Router) Handle(pattern string, handler http.Handler) {
	r.register(pattern, handler.ServeHTTP)
}

// Group creates a new sub-router with the given path appended to the base path of the parent router.
func (r *Router) Group(subPath string) *Router {
	return &Router{
		mux:  r.mux,
		path: r.path + subPath,
	}
}
