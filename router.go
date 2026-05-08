package thttp

import (
	"net/http"
)

type NewRouterFunc func() Router

// RouterType specifies the underlying router implementation.
// Supported types: net/http, julienschmidt/httprouter, gorilla/mux, gin-gonic/gin, go-chi/chi, labstack/echo.
type RouterType string

const (
	RouterTypeStd        RouterType = "net/http"
	RouterTypeHttprouter RouterType = "julienschmidt/httprouter"
	RouterTypeGorillaMux RouterType = "gorilla/mux"
	RouterTypeGin        RouterType = "gin-gonic/gin"
	RouterTypeChi        RouterType = "go-chi/chi"
	RouterTypeEcho       RouterType = "labstack/echo"
)

var allRouterTypes = []RouterType{}
var routerTypeMap = map[RouterType]NewRouterFunc{}

// Router is the interface that wraps the routing capabilities.
// Each router implementation (gin, echo, chi, etc.) must implement this interface.
type Router interface {
	// Use registers middlewares for the router.
	Use(middleware ...MiddlewareFunc)

	// Handle registers a handler for the given method and pattern.
	Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc)

	// Match finds a handler for the given request.
	// Returns the handler, path parameter function, and whether a match was found.
	Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamGetter, bool)

	// FormatSegment converts a thttp Segment to the router's native pattern syntax.
	FormatSegment(seg Segment) string
}

// PathParamGetter is an interface for accessing path parameter values.
type PathParamGetter interface {
	// Get returns the path parameter value by name.
	Get(name string) string
}

// RegisterRouter registers a new router implementation.
func RegisterRouter(routerType RouterType, newRouter NewRouterFunc) {
	allRouterTypes = append(allRouterTypes, routerType)
	routerTypeMap[routerType] = newRouter
}
