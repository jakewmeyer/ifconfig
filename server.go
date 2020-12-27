package main

import (
	"context"
	"errors"
	"ifconfig/ip"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const ContextTimeout = 10

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Heartbeat("/health"))
	r.Use(middleware.Logger)

	// CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	r.Use(cors.Handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Routes
	r.Get("/", ip.Get)

	// Handle graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeout)
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Run listen in goroutine to catch signal
	go func() {
		log.Printf("Starting on: %v", port)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-stop

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown Failed:%+v", err)
	}

	log.Print("Server Exited Properly")
}
