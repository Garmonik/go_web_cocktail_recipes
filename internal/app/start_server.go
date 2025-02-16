package app

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func ConfigServer(cfg *config.Config, router chi.Router) *http.Server {
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	return srv
}
