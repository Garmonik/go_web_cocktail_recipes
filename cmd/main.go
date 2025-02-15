package main

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/handlers"
	middleware_auth "github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/auth"
	middleware_base "github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/base"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"

	"github.com/Garmonik/go_web_cocktail_recipes/internal/config"
)

func main() {
	// set app config
	cfg := config.MustLoad()

	// set app logger
	log := logger.SetupLogger(cfg.Env)

	// set app database
	dataBase, err := db.New(cfg.DBPath)
	if err != nil {
		log.Error("Error opening database", "err", err)
		os.Exit(1)
	}
	_ = dataBase

	//set app router
	router := chi.NewRouter()

	// set app middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.Use(middleware_base.TrailingSlashMiddleware)

		r.Get("/login/", func(w http.ResponseWriter, r *http.Request) {
			handlers.LoginPage(w, r, cfg)
		})

		router.Group(func(r chi.Router) {
			r.Use(middleware_auth.AuthMiddleware)

		})
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug mod are enabled")

	// configuration server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Error starting server", "err", err)
	}

	log.Info("server stopped")
}
