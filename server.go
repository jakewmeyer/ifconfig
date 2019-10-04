package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"ifconfig/ip"
	"log"
	"net/http"
	"os"
)

func main() {
	router := httprouter.New()
	n := negroni.New()

	// Middleware
	n.UseHandler(router)
	n.Use(negroni.NewRecovery())
	n.Use(cors.Default())
	if os.Getenv("APP_ENV") != "production" {
		n.Use(negroni.NewLogger())
	}

	// Routes
	router.HandlerFunc("GET", "/", ip.Transform)

	// Server start
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}
