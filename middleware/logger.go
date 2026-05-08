package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/isayme/go-thttp"
)

// type contextHandler struct {
// 	slog.Handler
// }

// func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
// 	if kvs, ok := ctx.Value(thttp.LoggerCtxKey).(map[string]interface{}); ok {
// 		for k, v := range kvs {
// 			r.AddAttrs(slog.Any(k, v))
// 		}
// 	}
// 	return h.Handler.Handle(ctx, r)
// }

func Logger() thttp.MiddlewareFunc {
	hanlder := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	})
	logger := slog.New(hanlder)

	return func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			start := time.Now()

			// r := ctx.Request()
			// r = r.WithContext(context.WithValue(r.Context(), thttp.LoggerCtxKey, map[string]interface{}{}))
			// ctx.SetRequest(r)

			err := next(ctx)
			logger.InfoContext(ctx.Context(), "request", "method", ctx.Method(), "path", ctx.Request().RequestURI, "duration", time.Since(start))
			return err
		}
	}
}
