package main

import (
	"ifconfig/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/unrolled/secure"
	"ifconfig/ip"
)

func routes(s *server) *chi.Mux {
	// Middleware
	s.router.Use(middleware.RequestID)
	s.router.Use(logger.Logger(s.logger))
	s.router.Use(middleware.Recoverer)

	// Security headers
	s.router.Use(secure.New(secure.Options{
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src $NONCE",
	}).Handler)

	s.router.Use(middleware.Heartbeat("/health"))
	s.router.Use(middleware.Compress(5))
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.RedirectSlashes)

	// CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
	}))

	s.router.Get("/", ip.Get)

	return s.router
}
