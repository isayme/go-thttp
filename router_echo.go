package thttp

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type EchoMux struct {
	r           *echo.Echo
	middlewares []MiddlewareFunc
}

func echoNotFoundHandler(c *echo.Context) error {
	return nil
}

func echoMethodNotImplementedHandler(c *echo.Context) error {
	return nil
}

func NewEchoMux() Router {
	routerConfig := echo.RouterConfig{}
	// routerConfig.NotFoundHandler = echoNotFoundHandler
	// routerConfig.MethodNotAllowedHandler = echoMethodNotImplementedHandler

	return &EchoMux{
		r:           echo.NewWithConfig(echo.Config{Router: echo.NewRouter(routerConfig)}),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *EchoMux) FormatSegment(seg Segment) string {
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

func (router *EchoMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *EchoMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Add(method, pattern, func(echoCtx *echo.Context) error {
		r := echoCtx.Request()
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
		return nil
	})
}

func (router *EchoMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	echoCtx := echo.NewContext(r, w, router.r)
	handler := router.r.Router().Route(echoCtx)

	// notFound, notImplemented will cause error
	if err := handler(echoCtx); err != nil {
		return nil, nil, false
	}

	ctx := MustGetContextFromRequest(r)
	ctx.Set(PathRawParamsCtxKey, echoCtx)

	return MustGetHandlerFromCtx(ctx), NewEchoMuxPathParams, true
}

type EchoMuxPathParams struct {
	ctx Context
}

func NewEchoMuxPathParams(ctx Context) PathParams {
	return &EchoMuxPathParams{
		ctx: ctx,
	}
}

func (pp *EchoMuxPathParams) Get(name string) string {
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
