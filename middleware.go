package thttp

import "net/http"

type MiddlewareFunc func(next HandlerFunc) HandlerFunc

type Middlewares []func(http.Handler) http.Handler
