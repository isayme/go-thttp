package thttp

import "math/rand/v2"

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

func randomString() string {
	b := make([]byte, 10)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

func createHeaderMiddleware(headerKey, headerValue string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			ctx.SetHeader(headerKey, headerValue)
			return next(ctx)
		}
	}
}
