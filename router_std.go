package thttp

import (
	"fmt"
	"net/http"
	"strings"
)

type httpServeMux struct {
	r           *http.ServeMux
	middlewares []MiddlewareFunc
}

func newHttpServeMux() Router {
	return &httpServeMux{
		r:           http.NewServeMux(),
		middlewares: make([]MiddlewareFunc, 0),
	}
}

func (router *httpServeMux) FormatSegment(seg Segment) string {
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

func (router *httpServeMux) Use(middlewares ...MiddlewareFunc) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *httpServeMux) Handle(method, pattern string, h HandlerFunc, middleware ...MiddlewareFunc) {
	handler := applyMiddleware(h, middleware...)
	router.r.HandleFunc(fmt.Sprintf("%s %s", method, pattern), func(w http.ResponseWriter, r *http.Request) {
		ctx := MustGetContextFromRequest(r)
		SetHandlerInCtx(ctx, handler)
	})
}

func (router *httpServeMux) Match(w http.ResponseWriter, r *http.Request) (HandlerFunc, PathParamsFunc, bool) {
	handler, pattern := router.r.Handler(r)
	if pattern == "" {
		return nil, nil, false
	}

	handler.ServeHTTP(w, r)
	populatePathValues(r, pattern)

	ctx := MustGetContextFromRequest(r)

	return MustGetHandlerFromCtx(ctx), newHttpServeMuxPathParams, true
}

type httpServeMuxPathParams struct {
	ctx Context
}

func newHttpServeMuxPathParams(ctx Context) PathParams {
	return &httpServeMuxPathParams{ctx: ctx}
}

func (pp *httpServeMuxPathParams) Get(name string) string {
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
