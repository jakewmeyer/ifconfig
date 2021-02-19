package main

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"github.com/joho/godotenv"
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

	s, err := newServer(listenAddress)
	if err != nil {
		panic(err)
	}

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal(err.Error())
		}
	}()
	s.logger.Info("Server started", zap.String("port", port))

	<-stop

	s.logger.Info("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)

	defer func() {
		_ = s.logger.Sync()
		cancel()
		s.logger.Info("Server shutdown gracefully")
	}()

	if shutdownErr := s.srv.Shutdown(ctx); shutdownErr != nil {
		s.logger.Fatal("Server shutdown failed:", zap.String("error", shutdownErr.Error()))
	}
}
