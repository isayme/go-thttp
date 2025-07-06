package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isayme/go-thttp"
	"github.com/stretchr/testify/require"
)

func TestBasicAuth(t *testing.T) {
	require := require.New(t)

	app := thttp.New()
	creds := map[string]string{
		"demo": "passwd",
	}
	app.Use(BasicAuth("test", creds))
	app.Get("/hi", func(ctx thttp.Context) error {
		return ctx.String(200, "ok")
	})

	t.Run("pass", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hi", nil)
		headers := http.Header{}
		headers.Add("Authorization", "Basic ZGVtbzpwYXNzd2Q=")
		req.Header = headers

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		require.Equal(http.StatusOK, w.Code)
		require.Equal("ok", w.Body.String())
	})

	t.Run("not pass", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hi", nil)
		headers := http.Header{}
		headers.Add("Authorization", "Basic ZGVtbzphZG1pbg==")
		req.Header = headers

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		require.Equal(http.StatusUnauthorized, w.Code)
		require.Equal("", w.Body.String())
	})
}
