package handlers

import (
	middlewares "ToDoApp/internal/core/transport/http/middlewares"
	"fmt"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Init() http.Handler {
	var handler http.Handler = r.mux
	for _, middleware := range r.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) Use(middleware middlewares.Middleware) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) applyMiddlewares(handler http.Handler, mmiddlewares []middlewares.Middleware) http.Handler {
	for _, middleware := range mmiddlewares {
		handler = middleware(handler)
	}
	return handler
}

func (r *Router) Get(path string, handlerFunc http.HandlerFunc, mw ...middlewares.Middleware) {
	handler := r.applyMiddlewares(handlerFunc, mw)
	r.mux.Handle(fmt.Sprintf("GET %s", path), handler)
}

func (r *Router) Post(path string, handlerFunc http.HandlerFunc, mw ...middlewares.Middleware) {
	handler := r.applyMiddlewares(handlerFunc, mw)
	r.mux.Handle(fmt.Sprintf("POST %s", path), handler)
}

func (r *Router) Put(path string, handlerFunc http.HandlerFunc, mw ...middlewares.Middleware) {
	handler := r.applyMiddlewares(handlerFunc, mw)
	r.mux.Handle(fmt.Sprintf("PUT %s", path), handler)
}

func (r *Router) Delete(path string, handlerFunc http.HandlerFunc, mw ...middlewares.Middleware) {
	handler := r.applyMiddlewares(handlerFunc, mw)
	r.mux.Handle(fmt.Sprintf("DELETE %s", path), handler)
}
