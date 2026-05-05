package thttp

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HttprouterMux struct {
	r           *httprouter.Router
	middlewares []MiddlewareFunc
}

func NewHttprouterMux() *HttprouterMux {
	return &HttprouterMux{
		r:           httprouter.New(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *HttprouterMux) PatternType() PatternType {
	return ColonPattern
}

func (router *HttprouterMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *HttprouterMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Handle(method, pattern, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
	})
}

func (router *HttprouterMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	handler, params, redirect := router.r.Lookup(r.Method, r.URL.Path)
	if redirect {
		return nil, nil, false
	}
	if handler == nil {
		return nil, nil, false
	}

	handler(w, r, params)
	ctx := MustGetContextFromRequest(r)
	if len(params) > 0 {
		ctx.Set(httprouter.ParamsKey, params)
	}
	return MustGetHandlerFromCtx(ctx), NewHttprouterMuxPathParams, true
}

type HttprouterMuxPathParams struct {
	ctx Context
}

func NewHttprouterMuxPathParams(ctx Context) PathParams {
	return &HttprouterMuxPathParams{
		ctx: ctx,
	}
}

func (pp *HttprouterMuxPathParams) Get(name string) string {
	value := pp.ctx.Get(httprouter.ParamsKey)
	if value == nil {
		return ""
	}

	params, ok := value.(httprouter.Params)
	if ok {
		return params.ByName(name)
	}

	return ""
}
