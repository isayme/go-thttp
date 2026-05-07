package thttp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithNotFoundHandler(t *testing.T) {
	require := require.New(t)

	t.Run("custom notfound handler", func(t *testing.T) {
		errMsg := randomString()
		var handler HandlerFunc = func(Context) error {
			return fmt.Errorf(errMsg)
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
			return fmt.Errorf(errMsg)
		}

		app := New(WithErrorHandler(handler))

		require.Equal(errMsg, app.errorHandler(nil, nil).Error())
	})
}

func TestWithRouterType(t *testing.T) {
	require := require.New(t)

	t.Run("custom router type", func(t *testing.T) {
		app := New(WithRouterType(RouterTypeHttprouter))
		require.IsType(&HttprouterMux{}, app.router)
	})
}
