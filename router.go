package thttp

import (
	"net/http"
)

var _ Router = &HttpServeMux{}

type Router interface {
	PatternType() PatternType

	Use(middleware ...MiddlewareFunc)

	Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc)

	Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool)
}

type wrapHandler struct {
	h HandlerFunc
}

func newWrapHandler(h HandlerFunc) *wrapHandler {
	return &wrapHandler{h: h}
}

func (wh *wrapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// do nothing
}

type PathParamsFunc func(ctx Context) PathParams

type PathParams interface {
	Get(name string) string
}

type HttpServeMuxPathParams struct {
	ctx Context
}

func NewHttpServeMuxPathParams(ctx Context) PathParams {
	return &HttpServeMuxPathParams{ctx: ctx}
}

func (pp *HttpServeMuxPathParams) Get(name string) string {
	return pp.ctx.Request().PathValue(name)
}
