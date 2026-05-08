package thttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

func init() {
	RegisterRouter(RouterTypeGorillaMux, newGorillaMux)
}

type gorillaMux struct {
	r           *mux.Router
	middlewares []MiddlewareFunc
}

func newGorillaMux() Router {
	return &gorillaMux{
		r:           mux.NewRouter(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *gorillaMux) FormatSegment(seg Segment) string {
	switch seg.Type {
	case Static:
		return seg.Name
	case Param:
		return "{" + seg.Name + "}"
	case CatchAll:
		return "{" + seg.Name + ":.*}"
	default:
		panic("not supported segment type")
	}
}

func (router *gorillaMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *gorillaMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)

	router.r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := mustGetContextFromRequest(r)
		setHandlerInCtx(ctx, handler)
	}).Methods(method)
}

func (router *gorillaMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamGetter, bool) {
	rm := mux.RouteMatch{}

	found := router.r.Match(r, &rm)
	if !found {
		return nil, nil, false
	}

	rm.Handler.ServeHTTP(w, r)
	ctx := mustGetContextFromRequest(r)
	ctx.Set(pathRawParamsCtxKey, rm.Vars)

	return mustGetHandlerFromCtx(ctx), newGorillaMuxPathParams(ctx), true
}

type gorillaMuxPathParams struct {
	ctx Context
}

func newGorillaMuxPathParams(ctx Context) PathParamGetter {
	return &gorillaMuxPathParams{
		ctx: ctx,
	}
}

func (pp *gorillaMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(pathRawParamsCtxKey)
	if value == nil {
		return ""
	}

	params, ok := value.(map[string]string)
	if ok {
		return params[name]
	}

	return ""
}
