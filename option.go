package thttp

type optionFunc func(app *App)

func WithPrefix(prefix string) optionFunc {
	return func(app *App) {
		app.prefix = prefix
	}
}

func WithNotFoundHandler(handler HandlerFunc) optionFunc {
	return func(app *App) {
		app.notFoundHandler = handler
	}
}

func WithErrorHandler(handler func(Context, error) error) optionFunc {
	return func(app *App) {
		app.errorHandler = handler
	}
}

func WithRouterType(typ RouterType) optionFunc {
	return func(app *App) {
		app.useRouter(typ)
	}
}
