package main

import (
	"log/slog"

	"github.com/isayme/go-thttp"
	"github.com/isayme/go-thttp/middleware"
)

func main() {
	app := thttp.New()
	app.Use(middleware.Recovery())

	app.Get("/", func(ctx thttp.Context) error {
		panic(ctx.QueryString())
	})

	addr := ":8080"
	if err := app.Start(addr); err != nil {
		slog.Error("start fail", "err", err)
	}
}
