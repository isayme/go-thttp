package thttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithNotFoundHandler(t *testing.T) {
	require := require.New(t)

	t.Run("custom notfound handler", func(t *testing.T) {
		errMsg := randomString()
		var handler HandlerFunc = func(Context) error {
			return errors.New(errMsg)
		}

		app := New(WithNotFoundHandler(handler))

		require.Equal(errMsg, app.notFoundHandler(nil).Error())
	})
}

func TestWithErrorHandler(t *testing.T) {
	require := require.New(t)

	t.Run("custom notfound handler", func(t *testing.T) {
		errMsg := randomString()
		var handler ErrorHandlerFunc = func(Context, error) error {
			return errors.New(errMsg)
		}

		app := New(WithErrorHandler(handler))

		require.Equal(errMsg, app.errorHandler(nil, nil).Error())
	})
}

func TestWithRouterType(t *testing.T) {
	require := require.New(t)

	t.Run("custom router type", func(t *testing.T) {
		app := New(WithRouterType(RouterTypeHttprouter))
		require.IsType(&httprouterMux{}, app.router)
	})
}

func TestWithPrefix(t *testing.T) {
	require := require.New(t)

	t.Run("prefix", func(t *testing.T) {
		app := New(WithPrefix("/v1"))
		app.Get("/hi", func(ctx Context) error {
			return ctx.String(http.StatusOK, "OK")
		})

		w := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodGet, "/v1/hi", nil)
		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("OK", w.Body.String())
	})
}
