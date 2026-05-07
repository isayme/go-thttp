package thttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// chi 静态优先
func TestChiDuplicaPattern(t *testing.T) {
	require := require.New(t)

	t.Run("allow duplicate route 1", func(t *testing.T) {
		key := randomString()

		app := New(WithRouterType(RouterTypeChi))

		app.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		app.Get("/task/VERSION", func(ctx Context) error {
			return ctx.String(http.StatusOK, "-VERSION-")
		})

		{
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", key), nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(key, w.Body.String())
		}

		{
			req := httptest.NewRequest(http.MethodGet, "/task/VERSION", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal("-VERSION-", w.Body.String())
		}
	})

	t.Run("allow duplicate route 2", func(t *testing.T) {
		key := randomString()

		app := New(WithRouterType(RouterTypeChi))

		app.Get("/task/VERSION", func(ctx Context) error {
			return ctx.String(http.StatusOK, "-VERSION-")
		})

		app.Get("/task/{tid}", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.PathParam("tid"))
		})

		{
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/task/%s", key), nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal(key, w.Body.String())
		}

		{
			req := httptest.NewRequest(http.MethodGet, "/task/VERSION", nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			require.Equal(http.StatusOK, w.Code)
			require.Equal("-VERSION-", w.Body.String())
		}
	})
}
