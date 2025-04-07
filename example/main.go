package main

import (
	"log/slog"

	"github.com/isayme/go-thttp"
	"github.com/isayme/go-thttp/middleware"
)

func main() {
	app := thttp.New()

	app.Use(middleware.RequestID())

	app.Get("/hello", func(ctx thttp.Context) {
		ctx.String(200, "hi")
	})

	slog.Error("start fail: %s", app.Start(":8080"))
}
