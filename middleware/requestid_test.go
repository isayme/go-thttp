package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/isayme/go-thttp"
	"github.com/stretchr/testify/require"
)

func TestRequestId(t *testing.T) {
	require := require.New(t)

	app := thttp.New()
	app.Use(RequestID())
	app.Get("/hi", func(ctx thttp.Context) error {
		return ctx.String(200, "ok")
	})

	t.Run("use request specified id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hi", nil)
		headers := http.Header{}

		reqId := generator()

		headers.Add(thttp.HeaderXRequestID, reqId)
		req.Header = headers

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		require.Equal(http.StatusOK, w.Code)
		require.Equal(reqId, w.Header().Get(thttp.HeaderXRequestID))
	})

	t.Run("auto generate id", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/hi", nil)

		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		require.Equal(http.StatusOK, w.Code)
		require.NotEmpty(w.Header().Get(thttp.HeaderXRequestID))
	})
}
