package facecontrol

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func URLs(cfg *config.Config, r *chi.Mux, log *slog.Logger, dataBase *db.DataBase) {
	renderer := New(cfg, r, log, dataBase)

	r.Post("/api/login/", renderer.LoginUser)
	r.Post("/api/register/", renderer.RegisterUser)
	r.Post("/api/logout/", renderer.LogoutUser)
	return
}
