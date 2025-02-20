package users

import (
	"encoding/base64"
	"encoding/json"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
)

type User struct {
	cfg      *config.Config
	router   chi.Router
	log      *slog.Logger
	DataBase *db.DataBase
}

func New(cfg *config.Config, r chi.Router, log *slog.Logger, dataBase *db.DataBase) *User {
	return &User{cfg: cfg, router: r, log: log, DataBase: dataBase}
}

func (u *User) ShortUserInfo(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserByToken(r, u.cfg, u.DataBase)
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	u.log.Info("User info", "user", user)

	avatarData, err := os.ReadFile(user.Avatar.Path)
	if err != nil {
		http.Error(w, `{"error": "Failed to read avatar"}`, http.StatusBadRequest)
		return
	}

	avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
	response := map[string]string{
		"username": user.Name,
		"avatar":   avatarBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
