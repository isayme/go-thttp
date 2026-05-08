package middleware

import (
	"fmt"
	"net/http"

	"github.com/isayme/go-thttp"
)

func Recovery() thttp.MiddlewareFunc {
	return func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			defer func() {
				if r := recover(); r != nil {
					ctx.String(http.StatusInternalServerError, fmt.Sprintf("%v", r))
				}
			}()
			return next(ctx)
		}
	}
}
