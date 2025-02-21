package posts

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func URLs(cfg *config.Config, r chi.Router, log *slog.Logger, dataBase *db.DataBase) {
	usersUrl := New(cfg, r, log, dataBase)

	r.Get("/api/recipes/", usersUrl.PostsListAPI)
	r.Get("/api/user/{id}/posts/", usersUrl.PostsListByUserAPI)
	r.Get("/api/recipes/{id}/", usersUrl.PostsByIdAPI)
	r.Post("/api/recipes/", usersUrl.PostCreate)
	r.Post("/api/post/{id}/like/", usersUrl.LikeAPI)
	return
}
