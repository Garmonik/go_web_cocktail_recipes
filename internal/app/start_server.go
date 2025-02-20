package app

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/handlers/facecontrol"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/handlers/render_page"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/handlers/users"
	middleware_auth "github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/auth"
	middlewarebase "github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/base"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"mime"
	"net/http"
	"os"
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

func StandardMiddleware(r *chi.Mux, log *slog.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	if log != nil {
		r.Use(logger.New(log))
	}
}

func SetupRouter(log *slog.Logger, cfg *config.Config, dataBase *db.DataBase) *chi.Mux {
	router := chi.NewRouter()

	// middleware
	StandardMiddleware(router, log)
	err := mime.AddExtensionType(".css", "text/css")
	if err != nil {
		log.Error("Failed to add css", "error", err)
		os.Exit(1)
	}

	// static
	router.Group(func(r chi.Router) {
		r.Use(middleware.StripSlashes)
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
		r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "static/images/general/logo.ico")
		})
	})

	// URLs list without auth
	router.Group(func(r chi.Router) {
		r.Use(middleware.URLFormat)
		r.Use(middlewarebase.TrailingSlashMiddleware)
		r.Use(func(next http.Handler) http.Handler {
			return middleware_auth.AuthMiddleware(next, cfg, dataBase)
		})
		render_page.URLs(cfg, r, log)
		facecontrol.URLs(cfg, r, log, dataBase)
		users.URLs(cfg, r, log, dataBase)
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	log.Info("starting application", slog.String("env", cfg.Env))
	log.Debug("debug mod are enabled")
	return router
}
