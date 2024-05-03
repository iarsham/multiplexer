package multiplexer

import "net/http"

// Middleware is a function that takes a http.Handler and returns a http.Handler
type Middleware func(http.Handler) http.Handler

// Chain represents a chain of middlewares.
type Chain struct {
	mws []Middleware
}

// NewChain creates a new Chain with the provided middlewares.
func NewChain(mws ...Middleware) Chain {
	return Chain{append(([]Middleware)(nil), mws...)}
}

// Wrap wraps the provided http.Handler with the chain of middlewares.
func (c Chain) Wrap(h http.Handler) http.Handler {
	for i := range c.mws {
		h = c.mws[len(c.mws)-1-i](h)
	}
	return h
}

// WrapFunc wraps the provided http.HandlerFunc with the chain of middlewares.
func (c Chain) WrapFunc(h http.HandlerFunc) http.Handler {
	return c.Wrap(h)
}

// Append appends the provided middlewares to the chain.
func (c Chain) Append(mws ...Middleware) Chain {
	newMdw := make([]Middleware, 0, len(c.mws)+len(mws))
	newMdw = append(newMdw, c.mws...)
	newMdw = append(newMdw, mws...)
	return Chain{newMdw}
}
