package thttp

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func init() {
	RegisterRouter(RouterTypeGin, newGinMux)
}

type ginMux struct {
	r           *gin.Engine
	middlewares []MiddlewareFunc
}

func newGinMux() Router {
	ginEngine := gin.New()
	gin.SetMode(gin.ReleaseMode)
	// ginEngine.HandleMethodNotAllowed = false
	// ginEngine.NoRoute(func(ctx *gin.Context) {
	// 	// ctx.Writer = newFakeResponseWriter()
	// })
	return &ginMux{
		r:           ginEngine,
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *ginMux) FormatSegment(seg Segment) string {
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

func (router *ginMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *ginMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Handle(method, pattern, func(ginCtx *gin.Context) {
		r := ginCtx.Request
		ctx := mustGetContextFromRequest(r)
		setHandlerInCtx(ctx, handler)

		ctx.Set(pathRawParamsCtxKey, ginCtx.Params)
		ctx.Set(handlerFoundKey, true)
	})
}

func (router *ginMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamGetter, bool) {
	router.r.ServeHTTP(newFakeResponseWriter(), r)

	ctx := mustGetContextFromRequest(r)
	if _, ok := ctx.Get(handlerFoundKey).(bool); !ok {
		return nil, nil, false
	}
	return mustGetHandlerFromCtx(ctx), newGinMuxPathParams(ctx), true
}

type ginMuxPathParams struct {
	ctx Context
}

func newGinMuxPathParams(ctx Context) PathParamGetter {
	return &ginMuxPathParams{
		ctx: ctx,
	}
}

func (pp *ginMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(pathRawParamsCtxKey)
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

// use fake response to avoid gin framework auto write response
type fakeResponseWriter struct {
	header http.Header
}

var _ http.ResponseWriter = (*fakeResponseWriter)(nil)

func (w *fakeResponseWriter) WriteHeader(code int) {
}

func (w *fakeResponseWriter) Header() http.Header {
	return w.header
}

func (w *fakeResponseWriter) Write(data []byte) (n int, err error) {
	return len(data), nil
}

func newFakeResponseWriter() http.ResponseWriter {
	return &fakeResponseWriter{
		header: make(http.Header),
	}
}
