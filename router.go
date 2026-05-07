package thttp

import (
	"net/http"
)

type newRouterFunc func() Router
type RouterType string

const (
	RouterTypeStd        RouterType = "net/http"
	RouterTypeHttprouter RouterType = "httprouter"
	RouterTypeGorillaMux RouterType = "gorilla/mux"
)

var allRouterTypes = []RouterType{
	RouterTypeStd,
	RouterTypeHttprouter,
	RouterTypeGorillaMux,
}

var routerTypeMap = map[RouterType]newRouterFunc{
	RouterTypeStd:        NewHttpServeMux,
	RouterTypeHttprouter: NewHttprouterMux,
	RouterTypeGorillaMux: NewGorillaMux,
}

type Router interface {
	Use(middleware ...MiddlewareFunc)

	Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc)

	Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool)

	FormatSegment(seg Segment) string
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
