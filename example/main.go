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
	// app.Use(middleware.BasicAuth("dev", map[string]string{"admin": "123456"}))

	app.Get("/hello", func(ctx thttp.Context) error {
		ctx.String(200, "hi")
		return nil
	})

	app.Get("/abc/{key}/{value}", func(ctx thttp.Context) error {
		// ctx.String(200, "hi")
		ctx.String(200, fmt.Sprintf("k: %s, v: %s, k1: %s, k2: %s",
			ctx.PathParam("key"), ctx.PathParam("value"),
			ctx.Request().PathValue("key"), ctx.Request().PathValue("value"),
		))
		return nil
	})

	app.Get("/error", func(ctx thttp.Context) error {
		return errors.New("got error")
	})

	g1 := app.Group("/v1")
	g1.Use(func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			slog.Info("group v1 md 1")
			return next(ctx)
		}
	})
	g1.Get("/tasks/{tid}", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"id": ctx.PathParam("tid"),
		})
	})

	g1.Use(func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			slog.Info("group v1 md 2")
			return next(ctx)
		}
	})

	g2 := g1.Group("/msg")
	g2.Use(func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			slog.Info("group v2 md 1")
			return next(ctx)
		}
	})
	g2.Get("/tasks/{tid}", func(ctx thttp.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"id": ctx.PathParam("tid"),
		})
	})

	g2.Use(func(next thttp.HandlerFunc) thttp.HandlerFunc {
		return func(ctx thttp.Context) error {
			slog.Info("group v2 md 2")
			return next(ctx)
		}
	})

	slog.Error("start fail", "err", app.Start(":1323"))
}
