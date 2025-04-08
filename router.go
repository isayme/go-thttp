package thttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

var _ Router = &MuxRouter{}

var allowedMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type Router interface {
	Use(middlewares ...MiddlewareFunc)

	Get(pattern string, h HandlerFunc)
	Head(pattern string, h HandlerFunc)
	Post(pattern string, h HandlerFunc)
	Put(pattern string, h HandlerFunc)
	Patch(pattern string, h HandlerFunc)
	Del(pattern string, h HandlerFunc)
	Options(pattern string, h HandlerFunc)
	Trace(pattern string, h HandlerFunc)

	Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, bool)
}

type MuxRouter struct {
	r           *mux.Router
	middlewares []MiddlewareFunc
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{
		r:           mux.NewRouter(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *MuxRouter) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

type noopHandler struct {
	h HandlerFunc
}

func newNoopHandler(h HandlerFunc) *noopHandler {
	return &noopHandler{h: h}
}

func (nh *noopHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func (nh *noopHandler) Handler() HandlerFunc {
	return nh.h
}

func (router *MuxRouter) addRoute(method, pattern string, h HandlerFunc) {
	// router.r.Methods(method).Path(pattern).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	value := r.Context().Value(ContextKey)
	// 	ctx, _ := value.(Context)
	// 	h(ctx)
	// })
	router.r.Methods(method).Path(pattern).Handler(newNoopHandler(h))
}

func (router *MuxRouter) Get(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodGet, pattern, h)
}

func (router *MuxRouter) Head(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodHead, pattern, h)
}

func (router *MuxRouter) Post(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodPost, pattern, h)
}

func (router *MuxRouter) Put(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodPut, pattern, h)
}

func (router *MuxRouter) Patch(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodPatch, pattern, h)
}
func (router *MuxRouter) Del(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodDelete, pattern, h)
}
func (router *MuxRouter) Connect(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodConnect, pattern, h)
}
func (router *MuxRouter) Options(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodOptions, pattern, h)
}
func (router *MuxRouter) Trace(pattern string, h HandlerFunc) {
	router.addRoute(http.MethodTrace, pattern, h)
}

func (router *MuxRouter) Any(pattern string, h HandlerFunc) {
	for _, method := range allowedMethods {
		router.addRoute(method, pattern, h)
	}
}

func (router *MuxRouter) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, bool) {
	var match mux.RouteMatch
	var handler http.Handler

	if router.r.Match(r, &match) {
		ctx, _ := r.Context().Value(ContextKey).(Context)
		ctx.Set(PathParamsCtxKey, match.Vars)
		handler = match.Handler
		nh, _ := handler.(*noopHandler)
		return nh.Handler(), true
	}

	return nil, false
}
