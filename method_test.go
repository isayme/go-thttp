package thttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpMethod(t *testing.T) {
	require := require.New(t)

	t.Run("GET", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/method", nil)

		w := httptest.NewRecorder()

		app := New()
		app.Get("/method", func(ctx Context) error {
			return ctx.String(http.StatusOK, ctx.Method())
		})

		app.ServeHTTP(w, req)

		require.Equal(http.StatusOK, w.Code)
		require.Equal("GET", w.Body.String())
	})

	// methods := []string{
	// 	"GET",
	// 	"PUT",
	// 	"POST",
	// 	"DELETE",
	// 	"HEAD",
	// }

	// for _, method := range methods {
	// 	t.Run(method, func(t *testing.T) {

	// 		req := httptest.NewRequest(method, "/method", nil)

	// 		w := httptest.NewRecorder()

	// 		app := New()
	// 		app.Get()

	// 		app.ServeHTTP(w, req)

	// 		require.Equal(http.StatusOK, w.Code)
	// 		require.Equal("")
	// 	})
	// }

}
