package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isayme/go-thttp"
	"github.com/stretchr/testify/require"
)

func TestRecovery(t *testing.T) {
	require := require.New(t)

	t.Run("recovery", func(t *testing.T) {
		app := thttp.New()
		app.Use(Recovery())
		app.Get("/", func(ctx thttp.Context) error {
			panic("test")
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		require.Equal(http.StatusInternalServerError, w.Code)
		require.Equal("test", w.Body.String())
	})
}
