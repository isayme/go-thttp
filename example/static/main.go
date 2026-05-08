package main

import (
	"log/slog"
	"net/http"

	"github.com/isayme/go-thttp"
	"github.com/isayme/go-thttp/middleware"
)

func main() {
	app := thttp.New()
	app.Use(middleware.Recovery())

	app.Get("/hi", func(ctx thttp.Context) error {
		return ctx.String(http.StatusOK, ctx.QueryString())
	})

	app.Static("/", "./public")

	addr := ":8080"
	if err := app.Start(addr); err != nil {
		slog.Error("start fail", "err", err)
	}
}
