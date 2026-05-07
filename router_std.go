package thttp

import (
	"fmt"
	"net/http"
	"strings"
)

type HttpServeMux struct {
	r           *http.ServeMux
	middlewares []MiddlewareFunc
}

func NewHttpServeMux() Router {
	return &HttpServeMux{
		r:           http.NewServeMux(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *HttpServeMux) FormatSegment(seg Segment) string {
	switch seg.Type {
	case Static:
		return seg.Name
	case Param:
		return "{" + seg.Name + "}"
	case CatchAll:
		return "{" + seg.Name + "...}"
	default:
		panic("not supported segment type")
	}
}

func (router *HttpServeMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *HttpServeMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.Handle(fmt.Sprintf("%s %s", method, pattern), newWrapHandler(handler))
}

func (router *HttpServeMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	handler, pattern := router.r.Handler(r)
	if pattern == "" {
		return nil, nil, false
	}

	wh, ok := handler.(*wrapHandler)
	if !ok {
		panic("handler is not wrapHandler:" + pattern)
	}

	populatePathValues(r, pattern)

	return wh.h, NewHttpServeMuxPathParams, true
}

type HttpServeMuxPathParams struct {
	ctx Context
}

func NewHttpServeMuxPathParams(ctx Context) PathParams {
	return &HttpServeMuxPathParams{ctx: ctx}
}

func (pp *HttpServeMuxPathParams) Get(name string) string {
	return pp.ctx.Request().PathValue(name)
}

func populatePathValues(r *http.Request, pattern string) {
	if pattern == "" {
		return
	}

	if i := strings.Index(pattern, " "); i >= 0 {
		pattern = pattern[i+1:]
	}

	pSeg := split(pattern)
	sSeg := split(r.URL.Path)

	for i := 0; i < len(pSeg) && i < len(sSeg); i++ {
		ps := pSeg[i]

		if strings.HasPrefix(ps, "{") && strings.HasSuffix(ps, "}") {
			name := ps[1 : len(ps)-1]

			if strings.HasSuffix(name, "...") {
				name = strings.TrimSuffix(name, "...")
				r.SetPathValue(name, strings.Join(sSeg[i:], "/"))
				return
			}

			r.SetPathValue(name, sSeg[i])
		}
	}
}

func split(p string) []string {
	if len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}
