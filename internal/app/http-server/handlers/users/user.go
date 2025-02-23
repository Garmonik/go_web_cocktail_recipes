package users

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

	fmt.Println(user)
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
		http.Error(w, "Error while encoding JSON", http.StatusInternalServerError)
	}
}

func (u *User) UserInfo(w http.ResponseWriter, r *http.Request) {
	myUser, err := utils.GetUserByToken(r, u.cfg, u.DataBase)
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	var user models.User
	if err = u.DataBase.Db.Preload("Avatar").
		Select("id", "name", "email", "password", "avatar_id", "bio").
		Where("id = ?", userID).
		First(&user).Error; err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	avatarData, err := os.ReadFile(user.Avatar.Path)
	if err != nil {
		http.Error(w, `{"error": "failed to read avatar"}`, http.StatusBadRequest)
		return
	}

	avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
	response := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Name,
		"avatar":     avatarBase64,
		"bio":        user.Bio,
		"email":      user.Email,
		"my_account": user.ID == myUser.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"error": "Error while encoding JSON"}`, http.StatusInternalServerError)
	}
}

func (u *User) UpdateUser(w http.ResponseWriter, r *http.Request) {
	myUser, err := utils.GetUserByToken(r, u.cfg, u.DataBase)
	if err != nil {
		http.Error(w, `{"error": "user not found"}`, http.StatusNotFound)
		return
	}

	err = r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, `{"error": "Error processing form"}`, http.StatusBadRequest)
		return
	}

	bio := r.FormValue("bio")

	var avatarURL string
	file, handler, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()

		ext := strings.ToLower(filepath.Ext(handler.Filename))
		if _, ok := utils.AllowedExtensions[ext]; !ok {
			http.Error(w, `{"error": "Invalid file format"}`, http.StatusBadRequest)
			return
		}

		uuid := utils.GenerateUUID()
		filename := fmt.Sprintf("%s%s", uuid, ext)
		filePath := filepath.Join("file/users", filename)
		avatarURL = "file/users/" + filename

		if err := os.MkdirAll("file/users", os.ModePerm); err != nil {
			http.Error(w, `{"error": "Error creating directory"}`, http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, `{"error": "Error saving file"}`, http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, `{"error": "Error copying file"}`, http.StatusInternalServerError)
			return
		}
	} else {
		avatarURL = myUser.Avatar.Path
	}

	result := u.DataBase.Db.Model(&models.User{}).
		Where("id = ?", myUser.ID).
		Updates(map[string]interface{}{
			"bio": bio,
		})

	if result.Error != nil {
		http.Error(w, `{"error": "Error updating user"}`, http.StatusInternalServerError)
		return
	}

	result = u.DataBase.Db.Model(&models.Avatar{}).
		Where("id = ?", myUser.Avatar.ID).
		Updates(map[string]interface{}{
			"path": avatarURL,
		})

	response := map[string]string{
		"bio":    bio,
		"avatar": avatarURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
