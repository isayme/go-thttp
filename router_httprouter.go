package thttp

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type httprouterMux struct {
	r           *httprouter.Router
	middlewares []MiddlewareFunc
}

func newHttprouterMux() Router {
	return &httprouterMux{
		r:           httprouter.New(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *httprouterMux) FormatSegment(seg Segment) string {
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

func (router *httprouterMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *httprouterMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Handle(method, pattern, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
	})
}

func (router *httprouterMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
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
	return MustGetHandlerFromCtx(ctx), newHttprouterMuxPathParams, true
}

type httprouterMuxPathParams struct {
	ctx Context
}

func newHttprouterMuxPathParams(ctx Context) PathParams {
	return &httprouterMuxPathParams{
		ctx: ctx,
	}
}

func (pp *httprouterMuxPathParams) Get(name string) string {
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
