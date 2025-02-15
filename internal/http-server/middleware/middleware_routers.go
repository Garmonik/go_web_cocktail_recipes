package middleware_routers

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/handlers"
	middleware_auth "github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/auth"
	middleware_base "github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/base"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
)

func SetupRouter(log *slog.Logger, cfg *config.Config, dataBase *db.DataBase) *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)

	// Static
	router.Group(func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	})

	// URLS list without auth
	router.Group(func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.Use(middleware_base.TrailingSlashMiddleware)

		r.Get("/login/", func(w http.ResponseWriter, r *http.Request) {
			handlers.LoginPage(w, r, cfg)
		})
		r.Get("/register/", func(w http.ResponseWriter, r *http.Request) {
			handlers.RegisterPage(w, r, cfg)
		})
	})

	// URLS list with auth
	router.Group(func(r chi.Router) {
		r.Use(middleware_auth.AuthMiddleware)

	})

	return router
}
