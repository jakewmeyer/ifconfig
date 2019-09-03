package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"github.com/rs/cors"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// IPAddress representation
type IPAddress struct {
	IP string `json:"ip"`
}

func handleIP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip := net.ParseIP(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	// Plaintext by default, JSON with query param
	if _, ok := r.URL.Query()["json"]; ok {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IPAddress{ip.String()})
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, ip.String())
}

func main() {
	router := httprouter.New()
	n := negroni.New()

	// Middleware
	n.UseHandler(router)
	n.Use(cors.Default())
	if os.Getenv("APP_ENV") != "production" {
		n.Use(negroni.NewLogger())
	}
	
	// Routes
	router.GET("/", handleIP)

	// Server start
	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	log.Println("Starting on port: ", port)
	log.Fatal(http.ListenAndServe(":"+port, n))
}
