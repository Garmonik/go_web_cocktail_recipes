package users

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func URLs(cfg *config.Config, r *chi.Mux, log *slog.Logger, dataBase *db.DataBase) {
	users_url := New(cfg, r, log, dataBase)

	r.Post("/api/user/short/", users_url.ShortUserInfo)
	return
}
