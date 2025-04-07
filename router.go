package thttp

type Router interface {
	Use(middlewares ...MiddlewareFunc)

	Get(pattern string, h HandlerFunc)
	Head(pattern string, h HandlerFunc)
	Options(pattern string, h HandlerFunc)
	Patch(pattern string, h HandlerFunc)
	Put(pattern string, h HandlerFunc)
	Post(pattern string, h HandlerFunc)
	Del(pattern string, h HandlerFunc)
	Trace(pattern string, h HandlerFunc)
}
