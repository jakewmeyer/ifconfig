package main

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ContextTimeout = 5 * time.Second

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	listenAddress := ":" + port

	s, err := new(listenAddress)
	if err != nil {
		panic(err)
	}

	go func() {
		s.Logger.Info("Server started", zap.String("port", port))
		if err := s.Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Logger.Fatal(err.Error())
		}
	}()

	// Graceful Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	s.Logger.Info("Starting graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	if shutdownErr := s.Srv.Shutdown(ctx); shutdownErr != nil {
		s.Logger.Fatal("Server shutdown failed:", zap.String("error", shutdownErr.Error()))
	}
	defer func() {
		_ = s.Logger.Sync()
		cancel()
		s.Logger.Info("Server shutdown gracefully")
	}()
}
