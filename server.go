package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type server struct {
	router *chi.Mux
	logger *zap.Logger
	srv    *http.Server
}

func newServer(listenAddr string) (*server, error) {
	env := os.Getenv("GO_ENV")
	var config zap.Config
	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	logger, err := config.Build()
	if err != nil {
		logger.Fatal(err.Error())
	}

	errorLog, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		logger.Fatal(err.Error())
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
