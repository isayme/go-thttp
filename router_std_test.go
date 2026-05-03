package thttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStdRouter(t *testing.T) {
	require := require.New(t)

	t.Run("/task/{tid}", func(t *testing.T) {
		key := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", key), nil)
		w := httptest.NewRecorder()

		app := New()

		app.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(key, w.Body.String())
	})

	t.Run("/task/{k1}/{k2}", func(t *testing.T) {
		key1 := randomString()
		key2 := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s/%s", key1, key2), nil)
		w := httptest.NewRecorder()

		app := New()

		app.Get("/task/{k1}/{k2}", func(ctx Context) error {
			return ctx.String(http.StatusOK, fmt.Sprintf("%s\n%s", ctx.PathParam("k1"), ctx.PathParam("k2")))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(fmt.Sprintf("%s\n%s", key1, key2), w.Body.String())
	})

	t.Run("group /g1/task/{tid}", func(t *testing.T) {
		key := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/g1/task/%s", key), nil)
		w := httptest.NewRecorder()

		app := New()

		g1 := app.Group("/g1")

		g1.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(key, w.Body.String())
	})

	t.Run("nested group /g1/g2/task/{tid}", func(t *testing.T) {
		key := randomString()

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/g1/g2/task/%s", key), nil)
		w := httptest.NewRecorder()

		app := New()

		g1 := app.Group("/g1")
		g2 := g1.Group("/g2")

		g2.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal(key, w.Body.String())
	})
}
