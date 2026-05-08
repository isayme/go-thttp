package thttp

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type contextKey int

const (
	ContextKey contextKey = iota
	PathParamsCtxKey
	CatchAllPathParamCtxKey
	PathRawParamsCtxKey
	HandlerFoundKey
	RequestIDKey
	HandlerKey
	LoggerCtxKey
)

type App struct {
	pool sync.Pool

	prefix      string
	router      Router
	middlewares []MiddlewareFunc

	notFoundHandler HandlerFunc
	errorHandler    ErrorHandlerFunc

	logger *slog.Logger
}

// New creates a new thttp application.
// Options can be passed to customize the app behavior.
// Default router is net/http, can be changed via THTTP_ROUTER_TYPE env or WithRouterType option.
func New(options ...OptionFunc) *App {
	app := &App{
		middlewares: make([]MiddlewareFunc, 0),
	}

	app.defaultRouter()
	app.notFoundHandler = app.defaultNotFoundHandler
	app.errorHandler = app.defaultErrorHandler

	// apply options
	if len(options) > 0 {
		for _, option := range options {
			option(app)
		}
	}

	app.pool.New = func() any {
		return NewContext(nil, nil)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})
	app.logger = slog.New(handler)

	return app
}

// Logger returns the app's logger.
func (app *App) Logger() *slog.Logger {
	return app.logger
}

func (app *App) getHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx Context) error {
		h := applyMiddleware(handler, app.middlewares...)
		return h(ctx)
	}
}

func (app *App) formatPattern(pattern string) (string, string) {
	catchAllKey := ""
	segs := ParsePath(pattern)
	parts := make([]string, 0, len(segs))
	for _, s := range segs {
		part := app.router.FormatSegment(s)
		parts = append(parts, part)
		if s.Type == CatchAll {
			catchAllKey = s.Name
		}
	}

	return "/" + strings.Join(parts, "/"), catchAllKey
}

// Handle registers a handler for the given method and pattern.
// Pattern supports: static (/users), param (/users/:id), catch-all (/users/*path).
func (app *App) Handle(method, pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	pattern, catchAllKey := app.formatPattern(app.prefix + pattern)
	// slog.Info("handle", "method", method, "pattern", pattern)

	if catchAllKey != "" {
		md := func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Set(CatchAllPathParamCtxKey, catchAllKey)
				return next(ctx)
			}
		}
		middleware = append([]MiddlewareFunc{md}, middleware...)
	}

	app.router.Handle(method, pattern, applyMiddleware(app.getHandler(handler), middleware...))
}

// Get registers a GET handler for the given pattern.
func (app *App) Get(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodGet, pattern, handler, middleware...)
}

// Post registers a POST handler for the given pattern.
func (app *App) Post(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPost, pattern, handler, middleware...)
}

// Put registers a PUT handler for the given pattern.
func (app *App) Put(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPut, pattern, handler, middleware...)
}

// Patch registers a PATCH handler for the given pattern.
func (app *App) Patch(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodPatch, pattern, handler, middleware...)
}

// Delete registers a DELETE handler for the given pattern.
func (app *App) Delete(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodDelete, pattern, handler, middleware...)
}

// Head registers a HEAD handler for the given pattern.
func (app *App) Head(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodHead, pattern, handler, middleware...)
}

// Options registers an OPTIONS handler for the given pattern.
func (app *App) Options(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	app.Handle(http.MethodOptions, pattern, handler, middleware...)
}

// Any registers a handler for all HTTP methods for the given pattern.
func (app *App) Any(pattern string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	for _, method := range allowedHttpMethods {
		app.Handle(method, pattern, handler, middleware...)
	}
}

// Use registers global middlewares that apply to all routes.
func (app *App) Use(middleware ...MiddlewareFunc) {
	app.middlewares = append(app.middlewares, middleware...)
}

// Static serves static files from the given root directory.
// Pattern should end with *path to capture the file path.
func (app *App) Static(pattern string, root string) {
	app.Get(pattern+"*path", func(ctx Context) error {
		path := ctx.PathParam("path")
		path = filepath.Join(root, filepath.FromSlash(path))
		http.ServeFile(ctx.Response(), ctx.Request(), path)
		return nil
	})
}

// Group creates a route group with a common prefix and optional middlewares.
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

func (app *App) useRouter(typ RouterType) {
	fn, ok := routerTypeMap[typ]
	if !ok {
		panic("invalid router type")
	}

	app.router = fn()
}

func (app *App) defaultRouter() {
	typ := os.Getenv("THTTP_ROUTER_TYPE")
	if typ == "" {
		app.useRouter(RouterTypeStd)
	} else {
		app.useRouter(RouterType(typ))
	}
}

func (app *App) defaultNotFoundHandler(ctx Context) error {
	return ctx.String(http.StatusNotFound, "404 page not found")
}

func (app *App) defaultErrorHandler(ctx Context, err error) error {
	return ctx.String(http.StatusInternalServerError, err.Error())
}

// Start starts the HTTP server and listens on the given address.
func (app *App) Start(address string) error {
	return http.ListenAndServe(address, app)
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := app.pool.Get().(Context)
	ctx.Reset(r, w, app.logger)
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

// WrapHandler wraps an http.Handler into a thttp HandlerFunc.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// WrapHandlerFunc wraps an http.HandlerFunc into a thttp HandlerFunc.
func WrapHandlerFunc(h http.HandlerFunc) HandlerFunc {
	return func(c Context) error {
		h(c.Response(), c.Request())
		return nil
	}
}

// WrapMiddleware wraps a standard library middleware into a thttp MiddlewareFunc.
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

// HandlerFunc is the function signature for HTTP request handlers.
type HandlerFunc func(ctx Context) error

// ErrorHandlerFunc is the function signature for error handlers.
type ErrorHandlerFunc func(ctx Context, err error) error

// Skipper is a function that determines whether to skip the middleware.
type Skipper func(ctx Context) bool
