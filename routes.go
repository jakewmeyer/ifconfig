package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
)

func routes(s *server) *chi.Mux {
	// Middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(logger(s.Logger))
	s.Router.Use(middleware.Recoverer)

	// Security headers
	s.Router.Use(secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src $NONCE",
	}).Handler)

	s.Router.Use(middleware.Heartbeat("/health"))
	s.Router.Use(middleware.Compress(5))
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.RedirectSlashes)

	// CORS
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
	}))

	s.Router.Get("/", parseIP(s))

	return s.Router
}
