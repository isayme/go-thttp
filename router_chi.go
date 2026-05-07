package thttp

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ChiMux struct {
	r           *chi.Mux
	middlewares []MiddlewareFunc
}

func NewChiMux() Router {
	return &ChiMux{
		r:           chi.NewMux(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *ChiMux) FormatSegment(seg Segment) string {
	switch seg.Type {
	case Static:
		return seg.Name
	case Param:
		return "{" + seg.Name + "}"
	case CatchAll:
		return "*"
	default:
		panic("not supported segment type")
	}
}

func (router *ChiMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *ChiMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	slog.Info("chi", "method", method, "pattern", pattern)
	router.r.MethodFunc(method, pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
	})
}

func (router *ChiMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	chiCtx := chi.NewRouteContext()
	found := router.r.Match(chiCtx, r.Method, r.URL.Path)
	if !found {
		return nil, nil, false
	}

	router.r.ServeHTTP(w, r)

	ctx := MustGetContextFromRequest(r)
	ctx.Set(PathRawParamsCtxKey, chiCtx)

	return MustGetHandlerFromCtx(ctx), NewChiMuxPathParams, true
}

type ChiMuxPathParams struct {
	ctx Context
}

func NewChiMuxPathParams(ctx Context) PathParams {
	return &ChiMuxPathParams{
		ctx: ctx,
	}
}

func (pp *ChiMuxPathParams) Get(name string) string {
	chiCtx, ok := pp.ctx.Get(PathRawParamsCtxKey).(*chi.Context)
	if ok {
		if v, ok := pp.ctx.Get(CatchAllPathParamCtxKey).(string); ok && v == name {
			name = "*"
		}
		return chiCtx.URLParam(name)
	}

	return ""
}
