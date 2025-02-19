package users

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type User struct {
	cfg      *config.Config
	router   *chi.Mux
	log      *slog.Logger
	DataBase *db.DataBase
}

func New(cfg *config.Config, r *chi.Mux, log *slog.Logger, dataBase *db.DataBase) *User {
	return &User{cfg: cfg, router: r, log: log, DataBase: dataBase}
}

func (u *User) ShortUserInfo(w http.ResponseWriter, r *http.Request) {

}
