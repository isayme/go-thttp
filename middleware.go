package thttp

// MiddlewareFunc is a function that wraps a handler.
// It takes the next handler and returns a new handler that can execute
// code before and after the next handler is called.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc
