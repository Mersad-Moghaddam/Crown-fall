package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crownfall/backend/internal/game/content"
	"crownfall/backend/internal/platform/config"
	"crownfall/backend/internal/platform/logging"
	"crownfall/backend/internal/realtime/rooms"
)

func main() {
	configuration := config.Load()
	logger := logging.New(configuration.LogLevel)
	if err := content.Validate(); err != nil {
		logger.Error("invalid game content", "error", err)
		os.Exit(1)
	}
	registry := rooms.NewRegistry()
	server := &http.Server{Addr: configuration.HTTPAddress, Handler: newHandler(logger, registry), ReadHeaderTimeout: 5 * time.Second, IdleTimeout: 60 * time.Second}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Info("server started", "address", configuration.HTTPAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdown, cancel := context.WithTimeout(context.Background(), configuration.ShutdownTimeout)
	defer cancel()
	registry.Close()
	if err := server.Shutdown(shutdown); err != nil {
		logger.Error("shutdown failed", "error", err)
	}
	logger.Info("server stopped")
}
