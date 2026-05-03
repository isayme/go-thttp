package thttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	require := require.New(t)

	t.Run("Any", func(t *testing.T) {
		for _, method := range allowedHttpMethods {
			t.Run(method, func(t *testing.T) {
				req := httptest.NewRequest(method, "/v1/method", nil)

				w := httptest.NewRecorder()

				app := New()
				g := app.Group("/v1")
				g.Any("/method", func(ctx Context) error {
					return ctx.String(http.StatusOK, ctx.Method())
				})

				app.ServeHTTP(w, req)

				require.Equal(http.StatusOK, w.Code)
				require.Equal(method, w.Body.String())
			})
		}
	})

	t.Run("Get", func(t *testing.T) {
		method := http.MethodGet

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Post", func(t *testing.T) {
		method := http.MethodPost

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Post("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Put", func(t *testing.T) {
		method := http.MethodPut

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Put("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Patch", func(t *testing.T) {
		method := http.MethodPatch

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Patch("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Delete", func(t *testing.T) {
		method := http.MethodDelete

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Delete("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Options", func(t *testing.T) {
		method := http.MethodOptions

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Options("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Head", func(t *testing.T) {
		method := http.MethodHead

		req := httptest.NewRequest(method, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")
		g.Head("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(method, w.Body.String())
	})

	t.Run("Nested Group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/v2/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g1 := app.Group("/v1")
		g2 := g1.Group("/v2")

		g2.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
	})

	t.Run("Group Middleware: add md when create group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1", createHeaderMiddleware("h1", "hv1"))
		g.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal("hv1", w.Header().Get("h1"))
	})

	t.Run("Group Middleware: add md after create group", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")

		key := randomString()
		expected := randomString()
		g.Use(createHeaderMiddleware(key, expected))
		g.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})

	t.Run("Group Middleware: add md after add group route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")

		key := randomString()
		expected := randomString()
		g.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})
		g.Use(createHeaderMiddleware(key, expected))

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})

	t.Run("Group Middleware: add md when add group route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/v1/method", nil)

		w := httptest.NewRecorder()

		app := New()
		g := app.Group("/v1")

		key := randomString()
		expected := randomString()
		g.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		}, createHeaderMiddleware(key, expected))

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})
}
