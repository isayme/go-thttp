package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/isayme/go-thttp"
)

/**
 * see https://oneuptime.com/blog/post/2026-01-23-go-graceful-shutdown/view
 */
func main() {
	app := thttp.New()

	app.Get("/", func(ctx thttp.Context) error {
		time.Sleep(10 * time.Second)
		return ctx.String(http.StatusOK, ctx.QueryString())
	})

	// receive signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":8080"
	server := http.Server{
		Addr:    addr,
		Handler: app,
	}

	// Start server
	go func() {
		slog.Info("start listen", "addr", addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()

	slog.Info("shutting down start")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to stop server", "error", err)
	}
	slog.Info("shutting down end")
}
