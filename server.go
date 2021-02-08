package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type server struct {
	router *chi.Mux
	logger *zap.Logger
	srv    *http.Server
}

func newServer(listenAddr string) (*server, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Panicf("Can't initialize zap logger: %v", err)
	}

	errorLog, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		log.Panicf("Can't initialize zap error logger: %v", err)
	}

	server := &server{
		router: chi.NewRouter(),
		logger: logger,
	}

	server.srv = &http.Server{
		Addr:         listenAddr,
		Handler:      routes(server),
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return server, nil
}
