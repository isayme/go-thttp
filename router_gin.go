package thttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GinMux struct {
	r           *gin.Engine
	middlewares []MiddlewareFunc
}

func NewGinMux() Router {
	return &GinMux{
		r:           gin.Default(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *GinMux) FormatSegment(seg Segment) string {
	switch seg.Type {
	case Static:
		return seg.Name
	case Param:
		return ":" + seg.Name
	case CatchAll:
		if seg.Name == "*" {
			return "*"
		} else {
			return "*" + seg.Name
		}
	default:
		panic("not supported segment type")
	}
}

func (router *GinMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *GinMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Handle(method, pattern, func(ginCtx *gin.Context) {
		r := ginCtx.Request
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)

		ctx.Set(PathRawParamsCtxKey, ginCtx.Params)
		ctx.Set(HandlerFoundKey, true)
	})
}

func (router *GinMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	router.r.ServeHTTP(w, r)

	ctx := MustGetContextFromRequest(r)
	if _, ok := ctx.Get(HandlerFoundKey).(bool); !ok {
		return nil, nil, false
	}
	return MustGetHandlerFromCtx(ctx), NewGinMuxPathParams, true
}

type GinMuxPathParams struct {
	ctx Context
}

func NewGinMuxPathParams(ctx Context) PathParams {
	return &GinMuxPathParams{
		ctx: ctx,
	}
}

func (pp *GinMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(PathRawParamsCtxKey)
	if value == nil {
		return ""
	}

	params, ok := value.(gin.Params)
	if ok {
		value := params.ByName(name)
		// catch-all param has a leading slash, remove it.
		value = strings.TrimPrefix(value, "/")
		return value
	}

	return ""
}
