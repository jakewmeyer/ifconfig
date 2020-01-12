package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"ifconfig/ip"
	"log"
	"net/http"
	"os"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	r.Use(cors.Handler)

	// Routes
	r.Get("/", ip.Get)

	// Server start
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
