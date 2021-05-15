package main

import (
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

type server struct {
	Router *chi.Mux
	Logger *zap.Logger
	Srv    *http.Server
}

// New creates a new server
func new(listenAddr string) (*server, error) {
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
		return nil, err
	}

	errorLog, err := zap.NewStdLogAt(logger, zap.ErrorLevel)
	if err != nil {
		return nil, err
	}

	server := &server{
		Router: chi.NewRouter(),
		Logger: logger,
	}

	server.Srv = &http.Server{
		Addr:         listenAddr,
		Handler:      routes(server),
		ErrorLog:     errorLog,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return server, nil
}
