package thttp

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func init() {
	RegisterRouter(RouterTypeEcho, newEchoMux)
}

type echoMux struct {
	r           *echo.Echo
	middlewares []MiddlewareFunc
}

func newEchoMux() Router {
	routerConfig := echo.RouterConfig{}

	return &echoMux{
		r:           echo.NewWithConfig(echo.Config{Router: echo.NewRouter(routerConfig)}),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *echoMux) FormatSegment(seg Segment) string {
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

func (router *echoMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *echoMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Add(method, pattern, func(echoCtx *echo.Context) error {
		r := echoCtx.Request()
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
		return nil
	})
}

func (router *echoMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	echoCtx := echo.NewContext(r, w, router.r)
	handler := router.r.Router().Route(echoCtx)

	// notFound, notImplemented will cause error
	if err := handler(echoCtx); err != nil {
		return nil, nil, false
	}

	ctx := MustGetContextFromRequest(r)
	ctx.Set(PathRawParamsCtxKey, echoCtx)

	return MustGetHandlerFromCtx(ctx), newEchoMuxPathParams, true
}

type echoMuxPathParams struct {
	ctx Context
}

func newEchoMuxPathParams(ctx Context) PathParams {
	return &echoMuxPathParams{
		ctx: ctx,
	}
}

func (pp *echoMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(PathRawParamsCtxKey)
	if value == nil {
		return ""
	}

	echoCtx, ok := value.(*echo.Context)
	if ok {
		if v, ok := pp.ctx.Get(CatchAllPathParamCtxKey).(string); ok && v == name {
			name = "*"
		}

		value := echoCtx.Param(name)
		return value
	}

	return ""
}
