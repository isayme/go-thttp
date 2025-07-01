package thttp

import (
	"context"
	"net/http"
	"sync"
)

type contextKey int

const ContextKey contextKey = 0
const PathParamsCtxKey contextKey = 0
const RequestIDKey contextKey = 1

type App struct {
	pool sync.Pool

	router      Router
	middlewares []MiddlewareFunc
}

func New() *App {
	app := &App{
		router:      NewMuxRouter(),
		middlewares: make([]MiddlewareFunc, 0),
	}

	app.pool.New = func() any {
		return NewContext(nil, nil)
	}

	// app.router.Match()

	return app
}

func (app *App) Get(pattern string, handler HandlerFunc) {
	app.router.Get(pattern, handler)
}

func (app *App) Post(pattern string, handler HandlerFunc) {
	app.router.Post(pattern, handler)
}

func (app *App) Put(pattern string, handler HandlerFunc) {
	app.router.Put(pattern, handler)
}

func (app *App) Patch(pattern string, handler HandlerFunc) {
	app.router.Patch(pattern, handler)
}

func (app *App) Delete(pattern string, handler HandlerFunc) {
	app.router.Del(pattern, handler)
}

func (app *App) Head(pattern string, handler HandlerFunc) {
	app.router.Head(pattern, handler)
}

func (app *App) Trace(pattern string, handler HandlerFunc) {
	app.router.Trace(pattern, handler)
}

func (app *App) Use(middleware ...MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middleware...)
}

func (app *App) Static(pattern string, root string) {

}

func (app *App) Start(address string) error {
	return http.ListenAndServe(address, app)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	r = r.WithContext(context.WithValue(r.Context(), ContextKey, ctx))
	ctx.SetRequest(r)

	h, ok := app.router.Match(w, r)
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}

	h = applyMiddleware(h, app.middlewares...)

	err := h(ctx)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}

// WrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(c Context) error {
		h(c.Response(), c.Request())
		return nil
	}
}

func WrapMiddleware(m func(http.Handler) http.Handler) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) (err error) {
			m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				c.SetResponse(w)
				err = next(c)
			})).ServeHTTP(c.Response(), c.Request())
			return
		}
	}
}

type HandlerFunc func(ctx Context) error

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Middlewares []func(http.Handler) http.Handler

type Skipper func(ctx Context) bool

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
