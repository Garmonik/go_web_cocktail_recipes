package main

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/logger"
	"os"
)

func main() {
	// setup config
	cfg := config.MustLoad()

	// setup base_logger
	log := logger.SetupLogger(cfg.Env)

	// setup database
	dataBase, err := db.New(cfg.DBPath)
	if err != nil {
		log.Error("Error opening database", "err", err)
		os.Exit(1)
	}

	//setup router
	router := app.SetupRouter(log, cfg, dataBase)

	// configuration server
	srv := app.ConfigServer(cfg, router)

	// run server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Error starting server", "err", err)
	}

	log.Info("server stopped")
}
