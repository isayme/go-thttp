package main

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/isayme/go-thttp"
	"github.com/isayme/go-thttp/middleware"
)

func main() {
	app := thttp.New()

	app.Use(middleware.RequestID())
	app.Use(middleware.BasicAuth("dev", map[string]string{"admin": "123456"}))

	app.Get("/hello", func(ctx thttp.Context) error {
		ctx.String(200, "hi")
		return nil
	})

	app.Get("/abc/{key}/{value}", func(ctx thttp.Context) error {
		// ctx.String(200, "hi")
		ctx.String(200, fmt.Sprintf("k: %s, v: %s", ctx.PathParam("key"), ctx.PathParam("value")))
		return nil
	})

	app.Get("/error", func(ctx thttp.Context) error {
		return errors.New("got error")
	})

	slog.Error("start fail: %s", app.Start(":8080"))
}
