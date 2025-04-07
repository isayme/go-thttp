package thttp

import (
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type App struct {
	pool sync.Pool

	router      *mux.Router
	middlewares []MiddlewareFunc
}

func New() *App {
	app := &App{
		router:      mux.NewRouter(),
		middlewares: make([]MiddlewareFunc, 0),
	}

	app.pool.New = func() any {
		return NewContext(nil, nil)
	}

	// app.router.Match()

	return app
}

func (app *App) Get(pattern string, handler Handler) {
	app.router.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

	})
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
	// ctx := NewContext(w, r)

	// handler := app.router
	// for _, middleware := range app.middlewares {
	// 	middleware()
	// }
}

// WrapHandler wraps `http.Handler` into `echo.HandlerFunc`.
func WrapHandler(h http.Handler) HandlerFunc {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
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

type Handler func(ctx Context)
