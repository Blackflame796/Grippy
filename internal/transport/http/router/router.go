// router/router.go
package router

import (
	"Grippy/internal/transport/http/middlewares"
	"net/http"
	"strings"
)

type Router struct {
	Prefix      string
	mainRouter  *MainRouter
	middlewares []middlewares.Middleware
}

func New(prefix string, mainRouter *MainRouter) *Router {
	prefixNew := prefix
	if !strings.HasPrefix(prefixNew, "/") {
		prefixNew = "/" + prefixNew
	}
	return &Router{
		Prefix:      prefixNew,
		mainRouter:  mainRouter,
		middlewares: make([]middlewares.Middleware, 0),
	}
}

func (r *Router) Use(mw middlewares.Middleware) *Router {
	r.middlewares = append(r.middlewares, mw)
	return r
}

func (r *Router) wrap(h http.Handler, localMW ...middlewares.Middleware) http.Handler {
	var final http.Handler = h

	allMW := make([]middlewares.Middleware, 0)
	allMW = append(allMW, localMW...)                  // Local middleware
	allMW = append(allMW, r.middlewares...)            // Router middleware
	allMW = append(allMW, r.mainRouter.middlewares...) // Global middleware

	// Apply in reverse order: from last to first
	for i := len(allMW) - 1; i >= 0; i-- {
		final = allMW[i](final)
	}
	return final
}

func (r *Router) Handle(path string, h http.Handler, mw ...middlewares.Middleware) {
	r.mainRouter.mux.Handle(r.Prefix+path, r.wrap(h, mw...))
}

func (r *Router) Get(path string, h http.HandlerFunc, mw ...middlewares.Middleware) {
	r.Handle(path, http.HandlerFunc(h), mw...)
}

func (r *Router) Post(path string, h http.HandlerFunc, mw ...middlewares.Middleware) {
	r.Handle(path, http.HandlerFunc(h), mw...)
}

func (r *Router) Put(path string, h http.HandlerFunc, mw ...middlewares.Middleware) {
	r.Handle(path, http.HandlerFunc(h), mw...)
}

func (r *Router) Delete(path string, h http.HandlerFunc, mw ...middlewares.Middleware) {
	r.Handle(path, http.HandlerFunc(h), mw...)
}
