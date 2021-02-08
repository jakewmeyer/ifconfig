package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ContextTimeout = 5 * time.Second

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	listenAddress := ":" + port

	server, err := newServer(listenAddress)
	if err != nil {
		panic(err)
	}

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()
	log.Printf("Server started on: %v", port)

	<-stop

	log.Print("Starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)

	defer func() {
		_ = server.logger.Sync()
		cancel()
		log.Print("Server shutdown gracefully")
	}()

	if shutdownErr := server.srv.Shutdown(ctx); shutdownErr != nil {
		log.Panicf("Server shutdown failed: %+v", shutdownErr)
	}
}
