package thttp

import (
	"context"
	"net/http"
	"sync"
)

type contextKey int

const (
	ContextKey contextKey = iota
	PathParamsCtxKey
	PathRawParamsCtxKey
	RequestIDKey
	HandlerKey
)

type App struct {
	pool sync.Pool

	router      Router
	middlewares []MiddlewareFunc

	notFoundHandler HandlerFunc
	errorHandler    func(Context, error) error
}

func New() *App {
	app := &App{
		router: NewHttpServeMux(),
		// router: NewGorillaMux(),
		// router:      NewHttprouterMux(),
		middlewares: make([]MiddlewareFunc, 0),
	}

	app.NotFound(app.defaultNotFoundHandler)
	app.ErrorHandler(app.defaultErrorHandler)

	app.pool.New = func() any {
		return NewContext(nil, nil)
	}

	return app
}

func (app *App) getHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		h := applyMiddleware(handler, app.middlewares...)
		return h(ctx)
	}
}

func (app *App) Handle(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	pattern = convertPattern(pattern, app.router.PatternType())
	app.router.Handle(method, pattern, applyMiddleware(app.getHandler(handler), middleware...))
}

func (app *App) Get(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodGet, pattern, handler, middleware...)
}

func (app *App) Post(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPost, pattern, handler, middleware...)
}

func (app *App) Put(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPut, pattern, handler, middleware...)
}

func (app *App) Patch(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPatch, pattern, handler, middleware...)
}

func (app *App) Delete(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodDelete, pattern, handler, middleware...)
}

func (app *App) Head(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodHead, pattern, handler, middleware...)
}

func (app *App) Options(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodOptions, pattern, handler, middleware...)
}

func (app *App) Any(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range allowedHttpMethods {
		app.Handle(method, pattern, handler, middleware...)
	}
}

func (app *App) Use(middleware ...MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middleware...)
}

func (app *App) Static(pattern string, root string) {

}

func (app *App) Group(prefix string, middleware ...MiddlewareFunc) *Group {
	g := &Group{
		app:         app,
		router:      app.router,
		parent:      nil,
		prefix:      prefix,
		middlewares: make([]MiddlewareFunc, 0),
	}
	g.Use(middleware...)
	return g
}

func (app *App) defaultNotFoundHandler(ctx Context) error {
	return ctx.String(http.StatusNotFound, "404 page not found")
}

func (app *App) NotFound(handler HandlerFunc) {
	app.notFoundHandler = handler
}

func (app *App) defaultErrorHandler(ctx Context, err error) error {
	return ctx.String(http.StatusInternalServerError, err.Error())
}

func (app *App) ErrorHandler(handler func(Context, error) error) {
	app.errorHandler = handler
}

func (app *App) Start(address string) error {
	return http.ListenAndServe(address, app)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	r = r.WithContext(context.WithValue(r.Context(), ContextKey, ctx))
	ctx.SetRequest(r)

	h, params, ok := app.router.Match(w, r)
	if !ok {
		h = app.notFoundHandler
	}

	if params != nil {
		ctx.SetPathParam(params(ctx))
	}

	err := h(ctx)
	if err != nil {
		app.errorHandler(ctx, err)
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

type Skipper func(ctx Context) bool
