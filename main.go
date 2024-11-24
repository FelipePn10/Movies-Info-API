package main

import (
	"context"
	"fmt"
	"log/slog"
	"movies-info-api/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to execute code", "error", err)
		os.Exit(1)
	}
	slog.Info("All systems offline")
}

func run() error {
	apiKey := os.Getenv("OMDB_KEY")
	if apiKey == "" {
		return fmt.Errorf("OMDB_KEY environment variable is not set")
	}

	handler := api.NewHandler(apiKey)

	server := &http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Error during server shutdown", "error", err)
		return err
	}

	slog.Info("Server gracefully stopped")
	return nil
}
