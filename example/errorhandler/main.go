package main

import (
	"errors"
	"log/slog"

	"github.com/isayme/go-thttp"
)

func main() {
	app := thttp.New(thttp.WithRouterType(thttp.RouterTypeGin))

	app.Get("/", func(ctx thttp.Context) error {
		return errors.New("error")
	})

	addr := ":8080"
	if err := app.Start(addr); err != nil {
		slog.Error("start fail", "err", err)
	}
}
