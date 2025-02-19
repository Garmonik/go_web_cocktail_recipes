package facecontrol

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Facecontrol struct {
	cfg      *config.Config
	router   chi.Router
	log      *slog.Logger
	DataBase *db.DataBase
}

func New(cfg *config.Config, r chi.Router, log *slog.Logger, dataBase *db.DataBase) *Facecontrol {
	return &Facecontrol{cfg: cfg, router: r, log: log, DataBase: dataBase}
}

func (facecontrol *Facecontrol) LoginUser(w http.ResponseWriter, req *http.Request) {
	facecontrol.log.Info("start LoginUser")
	if err := req.ParseForm(); err != nil {
		utils.JsonResponse400("Form parsing error", w)
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")

	if email == "" || password == "" {
		facecontrol.log.Error("email or password is empty")
		utils.JsonResponse400("Email and password are required", w)
		return
	}
	user, err := utils.CheckUser(email, password, facecontrol.DataBase)
	if err != nil {
		facecontrol.log.Error(err.Error())
		utils.JsonResponse400(err.Error(), w)
		return
	}
	err = GenerateToken(facecontrol, user, w)
	if err != nil {
		utils.JsonResponse400(err.Error(), w)
	}

	utils.JsonResponse200Facecontrol(user.ID, http.StatusOK, w)
	return
}

func (facecontrol *Facecontrol) RegisterUser(w http.ResponseWriter, req *http.Request) {
	facecontrol.log.Info("start RegisterUser")
	if err := req.ParseForm(); err != nil {
		utils.JsonResponse400("Form parsing error", w)
		return
	}

	email := req.FormValue("email")
	password := req.FormValue("password")
	password2 := req.FormValue("password2")
	name := req.FormValue("username")

	if email == "" || password == "" || password2 == "" || name == "" {
		facecontrol.log.Error("Email, password, username and password again are required")
		utils.JsonResponse400("Email, password, username and password again are required", w)
		return
	}
	if len(name) > 25 || len(name) <= 5 {
		facecontrol.log.Error("Name length exceeds 20 characters", w)
		utils.JsonResponse400("Name length exceeds 20 characters", w)
	}
	user, err := utils.CheckUserByEmail(email, facecontrol.DataBase)
	if err != "users not found" {
		facecontrol.log.Error("User with this date already exists")
		utils.JsonResponse400("User with this date already exists", w)
		return
	}
	user, err = utils.CheckUserByName(name, facecontrol.DataBase)
	if err != "users not found" {
		facecontrol.log.Error("User with this name already exists")
		utils.JsonResponse400("User with this date already exists", w)
		return
	}
	if password != password2 {
		facecontrol.log.Error("Passwords do not match")
		utils.JsonResponse400("Passwords do not match", w)
		return
	}

	user, errorCreate := CreateUser(email, password, name, facecontrol)
	if errorCreate != nil {
		facecontrol.log.Error(errorCreate.Error())
		utils.JsonResponse400(errorCreate.Error(), w)
		return
	}
	errTokens := GenerateToken(facecontrol, user, w)
	if errTokens != nil {
		utils.JsonResponse400(errTokens.Error(), w)
	}

	utils.JsonResponse200Facecontrol(user.ID, http.StatusCreated, w)
	return
}

func (facecontrol *Facecontrol) LogoutUser(w http.ResponseWriter, req *http.Request) {
	facecontrol.log.Info("start LogoutUser")
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   "localhost", // Должно совпадать с установленными куками!
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Удаляем refresh_token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   "localhost", // Должно совпадать с установленными куками!
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusNoContent)
	return
}
