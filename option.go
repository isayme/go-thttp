package thttp

type OptionFunc func(app *App)

func WithPrefix(prefix string) OptionFunc {
	return func(app *App) {
		app.prefix = prefix
	}
}

func WithNotFoundHandler(handler HandlerFunc) OptionFunc {
	return func(app *App) {
		app.notFoundHandler = handler
	}
}

func WithErrorHandler(handler ErrorHandlerFunc) OptionFunc {
	return func(app *App) {
		app.errorHandler = handler
	}
}

func WithRouterType(typ RouterType) OptionFunc {
	return func(app *App) {
		app.useRouter(typ)
	}
}
