package thttp

type MiddlewareFunc func(next HandlerFunc) HandlerFunc
