package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gateplaza.local/facility/internal/platform/config"
	"gateplaza.local/facility/internal/platform/db"
	"gateplaza.local/facility/internal/platform/httpx"
)

func main() { os.Exit(run()) }
func run() int {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		slog.New(slog.NewJSONHandler(os.Stderr, nil)).Error("configuration_invalid", "reason", err.Error())
		return 1
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.Level}))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("database_initialization_failed")
		return 1
	}
	defer pool.Close()
	server := httpx.Server(cfg.Address, httpx.Handler(pool, cfg.DBTimeout, cfg.Docs, logger))
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		logger.Error("http_listen_failed")
		return 1
	}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	logger.Info("server_started", "address", cfg.Address, "environment", cfg.Environment)
	select {
	case err = <-done:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http_server_failed")
			return 1
		}
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdown); err != nil {
			logger.Error("shutdown_failed")
			_ = server.Close()
			return 1
		}
		if err := <-done; !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http_server_failed")
			return 1
		}
	}
	logger.Info("server_stopped")
	return 0
}
