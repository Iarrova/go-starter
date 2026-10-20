package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/Iarrova/go-starter/internal/app"
	"github.com/Iarrova/go-starter/internal/config"
	"github.com/Iarrova/go-starter/internal/httpapi"
	"github.com/Iarrova/go-starter/internal/logging"
	"github.com/Iarrova/go-starter/internal/postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("loading config", "error", err)
		os.Exit(1)
	}

	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	logger := slog.New(logging.NewHandler(slog.NewJSONHandler(os.Stdout, opts)))
	slog.SetDefault(logger)

	if err := run(cfg, logger); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Config, logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, db, err := postgres.Open(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer func() {
		if err := client.Close(); err != nil {
			logger.Error("closing database", "error", err)
		}
	}()

	gin.SetMode(cfg.HTTP.GinMode)
	server := app.NewServer(client, db, logger)
	srv := server.Server(fmt.Sprintf(":%d", cfg.HTTP.Port))

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listening: %w", err)
	case <-ctx.Done():
		stop()
	}

	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpapi.ShutdownTimeout)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
