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
	"strconv"
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
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		os.Exit(1)
	}
}

func (u *User) MyUserInfo(w http.ResponseWriter, r *http.Request) {
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
		"id":       strconv.Itoa(int(user.ID)),
		"username": user.Name,
		"avatar":   avatarBase64,
		"bio":      user.Bio,
		"email":    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
	}
}
