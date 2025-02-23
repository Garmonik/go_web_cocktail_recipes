package render_page

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

func URLs(cfg *config.Config, r chi.Router, log *slog.Logger) {
	renderer := New(cfg, r, log)

	r.Get("/login/", renderer.LoginPage)
	r.Get("/register/", renderer.RegisterPage)
	r.Get("/home/", renderer.HomePage)
	r.Get("/my_user/", renderer.MyUserPage)
	r.Get("/my_user/update/", renderer.UpdatePage)
	r.Get("/recipes/", renderer.PostsList)
	r.Get("/user/{id}/", renderer.SomeUserPage)
	r.Get("/favorites/", renderer.FavoritesPostsList)
	r.Get("/recipes/create/", renderer.CreatePost)
	return
}
