package thttp

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type HttprouterMux struct {
	r           *httprouter.Router
	middlewares []MiddlewareFunc
}

func newHttprouterMux() Router {
	return &HttprouterMux{
		r:           httprouter.New(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *HttprouterMux) FormatSegment(seg Segment) string {
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
		ctx.Set(PathRawParamsCtxKey, params)
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
	value := pp.ctx.Get(PathRawParamsCtxKey)
	if value == nil {
		return ""
	}

	params, ok := value.(httprouter.Params)
	if ok {
		value := params.ByName(name)
		// catch-all param has a leading slash, remove it.
		value = strings.TrimPrefix(value, "/")
		return value
	}

	return ""
}
