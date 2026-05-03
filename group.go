package thttp

import (
	"fmt"
	"net/http"
)

type Group struct {
	router Router

	parent      *Group
	prefix      string
	middlewares []MiddlewareFunc
}

func (g *Group) Use(middleware ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}

func (g *Group) Handle(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.router.Handle(method, g.formatPattern(pattern), g.getHandler(handler), middleware...)
}

func (g *Group) Get(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodGet, pattern, handler, middleware...)
}

func (g *Group) Post(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodPost, pattern, handler, middleware...)
}

func (g *Group) Put(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodPut, pattern, handler, middleware...)
}

func (g *Group) Patch(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodPatch, pattern, handler, middleware...)
}

func (g *Group) Delete(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodDelete, pattern, handler, middleware...)
}

func (g *Group) Head(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodHead, pattern, handler, middleware...)
}

func (g *Group) Options(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	g.Handle(http.MethodOptions, pattern, handler, middleware...)
}

func (g *Group) Any(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range allowedHttpMethods {
		g.Handle(method, pattern, handler, middleware...)
	}
}

func (g *Group) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	sg := &Group{
		router:      g.router,
		parent:      g,
		prefix:      prefix,
		middlewares: make([]MiddlewareFunc, 0),
	}
	sg.Use(middleware...)
	return sg
}

// formatPattern add prefix
func (g *Group) formatPattern(pattern string) string {
	return fmt.Sprintf("%s%s", g.getPrefix(), pattern)
}

// getPrefix add parent prefix
func (g *Group) getPrefix() string {
	if g.parent != nil {
		return fmt.Sprintf("%s%s", g.parent.getPrefix(), g.prefix)
	}
	return g.prefix
}

func (g *Group) getFullMiddlewares() []MiddlewareFunc {
	if g.parent != nil {
		return append(g.parent.getFullMiddlewares(), g.middlewares...)
	}
	return g.middlewares
}

func (g *Group) getHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		h := applyMiddleware(handler, g.getFullMiddlewares()...)
		return h(ctx)
	}
}
