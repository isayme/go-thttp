package thttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	require := require.New(t)

	t.Run("Group Middleware: add md before add route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/method", nil)

		w := httptest.NewRecorder()

		app := New()

		key := randomString()
		expected := randomString()
		app.Use(createHeaderMiddleware(key, expected))
		app.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})

	t.Run("Middleware: add md after add route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/method", nil)

		w := httptest.NewRecorder()

		app := New()
		key := randomString()
		expected := randomString()
		app.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})
		app.Use(createHeaderMiddleware(key, expected))

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})

	t.Run("Middleware: add md when add route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/method", nil)

		w := httptest.NewRecorder()

		app := New()

		key := randomString()
		expected := randomString()
		app.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		}, createHeaderMiddleware(key, expected))

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(http.MethodGet, w.Body.String())
		require.Equal(expected, w.Header().Get(key))
	})
}

func TestMiddlewareOrder(t *testing.T) {
	require := require.New(t)

	t.Run("execute order", func(t *testing.T) {
		app := New()

		respBody := ""

		app.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				respBody += "1"
				return next(ctx)
			}
		})

		app.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				respBody += "2"
				return next(ctx)
			}
		})

		app.Get("/hi", func(ctx Context) error {
			respBody += "3"
			return ctx.String(http.StatusOK, respBody)
		})

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/hi", nil)
		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("123", w.Body.String())
	})

	t.Run("middleware error", func(t *testing.T) {
		app := New()

		app.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Response().Write([]byte("1"))
				return next(ctx)
			}
		})

		app.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Response().Write([]byte("2"))
				return fmt.Errorf("mderror")
			}
		})

		app.Use(func(next HandlerFunc) HandlerFunc {
			return func(ctx Context) error {
				ctx.Response().Write([]byte("3"))
				return next(ctx)
			}
		})

		app.Get("/hi", func(ctx Context) error {
			return ctx.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/hi", nil)
		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("12mderror", w.Body.String())
	})
}
