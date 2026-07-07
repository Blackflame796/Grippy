// router/main_router.go
package router

import (
	"Grippy/internal/transport/http/middlewares"
	"net/http"
)

type MainRouter struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
}

func NewMainRouter() *MainRouter {
	return &MainRouter{
		mux:         http.NewServeMux(),
		middlewares: make([]middlewares.Middleware, 0),
	}
}

func (r *MainRouter) Use(middleware middlewares.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *MainRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Apply global middleware to the root routes
	var handler http.Handler = r.mux

	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}

	handler.ServeHTTP(w, req)
}

func (r *MainRouter) ServeMux() *http.ServeMux {
	return r.mux
}
