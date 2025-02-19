package facecontrol

import (
	"errors"
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

func CreateUser(email string, password string, name string, facecontrol *Facecontrol) (*models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	// Select a random avatar
	//Creating a new random number generator
	src := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(src)
	randomAvatar := rng.Intn(7) + 1

	avatarPath := fmt.Sprintf("static/images/base_avatar/avatar%d.png", randomAvatar)

	// Create an avatar in the database
	avatar := models.Avatar{Path: avatarPath}
	if err := facecontrol.DataBase.Db.Create(&avatar).Error; err != nil {
		return nil, fmt.Errorf("failed to create avatar")
	}
	// Create a users in the database
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Bio:      "",
		AvatarID: avatar.ID,
		Avatar:   avatar,
	}

	if err := facecontrol.DataBase.Db.Create(&user).Error; err != nil {
		facecontrol.DataBase.Db.Delete(&avatar)
		return nil, fmt.Errorf("failed to create users")
	}

	return &user, nil
}

func GenerateToken(facecontrol *Facecontrol, user *models.User, w http.ResponseWriter) error {
	accessToken, errAccessToken := utils.GenerateAccessToken(user, facecontrol.cfg)
	refreshToken, errRefreshToken := utils.GenerateRefreshToken(user, facecontrol.cfg)
	if errAccessToken != nil || errRefreshToken != nil {
		facecontrol.log.Error("Error generating access token or refresh token")
		return errors.New("error generating access token or refresh token")
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(24 * 30 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil

}
