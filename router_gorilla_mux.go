package thttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type GorillaMux struct {
	r           *mux.Router
	middlewares []MiddlewareFunc
}

func newGorillaMux() Router {
	return &GorillaMux{
		r:           mux.NewRouter(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *GorillaMux) FormatSegment(seg Segment) string {
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

func (router *GorillaMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *GorillaMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)

	router.r.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
	}).Methods(method)
}

func (router *GorillaMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	rm := mux.RouteMatch{}

	found := router.r.Match(r, &rm)
	if !found {
		return nil, nil, false
	}

	rm.Handler.ServeHTTP(w, r)
	ctx := MustGetContextFromRequest(r)
	ctx.Set(PathRawParamsCtxKey, rm.Vars)

	return MustGetHandlerFromCtx(ctx), newGorillaMuxPathParams, true
}

type gorillaMuxPathParams struct {
	ctx Context
}

func newGorillaMuxPathParams(ctx Context) PathParams {
	return &gorillaMuxPathParams{
		ctx: ctx,
	}
}

func (pp *gorillaMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(PathRawParamsCtxKey)
	if value == nil {
		return ""
	}

	params, ok := value.(map[string]string)
	if ok {
		return params[name]
	}

	return ""
}
