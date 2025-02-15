package main

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/db"
	middleware_routers "github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/logger"
	"log/slog"
	"net/http"
	"os"

	"github.com/Garmonik/go_web_cocktail_recipes/internal/config"
)

func main() {
	// setup config
	cfg := config.MustLoad()
	// setup logger
	log := logger.SetupLogger(cfg.Env)

	// setup database
	dataBase, err := db.New(cfg.DBPath)
	if err != nil {
		log.Error("Error opening database", "err", err)
		os.Exit(1)
	}

	//setup router
	router := middleware_routers.SetupRouter(log, cfg, dataBase)

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
