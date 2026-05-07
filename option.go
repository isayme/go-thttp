package thttp

type optionFunc func(app *App)

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
