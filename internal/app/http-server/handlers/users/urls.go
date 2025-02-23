package users

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func URLs(cfg *config.Config, r chi.Router, log *slog.Logger, dataBase *db.DataBase) {
	usersUrl := New(cfg, r, log, dataBase)

	r.Get("/api/user/short/", usersUrl.ShortUserInfo)
	r.Get("/api/my_user/", usersUrl.MyUserInfo)
	r.Get("/api/user/{id}/", usersUrl.UserInfo)
	r.Patch("/api/my_user/", usersUrl.UpdateUser)
	return
}
