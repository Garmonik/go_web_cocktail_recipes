package middleware_routers

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/handlers"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/base"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func SetupRouter(log *slog.Logger, cfg *config.Config, dataBase *db.DataBase) *chi.Mux {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	// static
	router.Group(func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	})

	// URLs list without auth
	router.Group(func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.Use(middleware_base.TrailingSlashMiddleware)

		renderer := render_page.New(cfg)

		r.Get("/login/", renderer.LoginPage)
		r.Get("/register/", renderer.RegisterPage)
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug mod are enabled")
	return router
}
