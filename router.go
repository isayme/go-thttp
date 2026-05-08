package thttp

import (
	"net/http"
)

type newRouterFunc func() Router
type RouterType string

const (
	RouterTypeStd        RouterType = "net/http"
	RouterTypeHttprouter RouterType = "julienschmidt/httprouter"
	RouterTypeGorillaMux RouterType = "gorilla/mux"
	RouterTypeGin        RouterType = "gin-gonic/gin"
	RouterTypeChi        RouterType = "go-chi/chi"
	RouterTypeEcho       RouterType = "labstack/echo"
)

var allRouterTypes = []RouterType{
	RouterTypeStd,
	RouterTypeHttprouter,
	RouterTypeGorillaMux,
	RouterTypeGin,
	RouterTypeChi,
	RouterTypeEcho,
}

var routerTypeMap = map[RouterType]newRouterFunc{
	RouterTypeStd:        newHttpServeMux,
	RouterTypeHttprouter: newHttprouterMux,
	RouterTypeGorillaMux: newGorillaMux,
	RouterTypeChi:        newChiMux,
	RouterTypeEcho:       newEchoMux,
	RouterTypeGin:        newGinMux,
}

type Router interface {
	Use(middleware ...MiddlewareFunc)

	Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc)

	Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool)

	FormatSegment(seg Segment) string
}

type PathParamsFunc func(ctx Context) PathParams

type PathParams interface {
	Get(name string) string
}
