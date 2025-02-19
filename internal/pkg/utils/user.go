package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// hashWithHMAC adds security HMAC before bcrypt
func hashWithHMAC(password string, secretKey string) []byte {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(password))
	return h.Sum(nil)
}

// HashPassword hashes password with bcrypt
func HashPassword(password string, secretKey string) (string, error) {
	hmacHash := hashWithHMAC(password, secretKey)
	hashedPassword, err := bcrypt.GenerateFromPassword(hmacHash, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the password matches the hash
func CheckPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func CheckUserByEmail(email string, dataBase *db.DataBase) (*models.User, string) {
	var user models.User
	if err := dataBase.Db.Select("id", "name", "email", "password").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return &models.User{}, "users not found"
	}
	return &user, ""
}

func CheckUserByName(name string, dataBase *db.DataBase) (*models.User, string) {
	var user models.User
	if err := dataBase.Db.Select("id", "name", "email", "password").
		Where("name = ?", name).
		First(&user).Error; err != nil {
		return &models.User{}, "users not found"
	}
	return &user, ""
}

// CheckUser check users exist in db
func CheckUser(email string, password string, dataBase *db.DataBase) (*models.User, error) {
	user, err := CheckUserByEmail(email, dataBase)
	if err != "" {
		return nil, fmt.Errorf(err)
	}

	ok := CheckPassword(user, password)

	if ok {
		return user, nil
	}
	return nil, fmt.Errorf("the password will be entered incorrectly")
}

// GenerateAccessToken creates an ES256 token
func GenerateAccessToken(user *models.User, cfg *config.Config) (string, error) {
	now := time.Now()
	expirationTime := now.Add(15 * time.Minute).Unix() // +1 month

	claims := jwt.MapClaims{
		"exp":   expirationTime,
		"sub":   user.ID,
		"email": user.Email,
		"urn":   user.Name,
		"iat":   now.Unix(),
		"site":  "cocktail website",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Checking that privateKey is not nil and of the correct type
	if cfg.JWT.PrivateKey == nil {
		return "", fmt.Errorf("private key is nil")
	}

	privateKey := cfg.JWT.PrivateKey

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateRefreshToken creates an ES256 refresh_token
func GenerateRefreshToken(user *models.User, cfg *config.Config) (string, error) {
	now := time.Now()
	expirationTime := now.Add(30 * 24 * time.Hour).Unix() // +1 month

	claims := jwt.MapClaims{
		"exp":   expirationTime,
		"sub":   user.ID,
		"scope": "refresh",
		"iat":   now.Unix(),
		"site":  "cocktail website",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Checking that privateKey is not nil and of the correct type
	if cfg.JWT.PrivateKey == nil {
		return "", fmt.Errorf("private key is nil")
	}

	privateKey := cfg.JWT.PrivateKey

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// IsValidToken check token
func IsValidToken(tokenString string, cfg *config.Config) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return cfg.JWT.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return false
	}

	return true
}

func RefreshAccessToken(tokenString string, cfg *config.Config, db *db.DataBase) (string, string) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return cfg.JWT.PublicKey, nil
	})

	if err != nil || !token.Valid {
		return "", ""
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ""
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return "", ""
	}

	scope, ok := claims["scope"].(string)
	if !ok || scope != "refresh" {
		return "", ""
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", ""
	}

	user, errUser := CheckUserByEmail(email, db)
	if errUser != "" {
		return "", ""
	}

	accessToken, errToken := GenerateAccessToken(user, cfg)
	if errToken != nil {
		return "", ""
	}
	refreshToken, errToken := GenerateRefreshToken(user, cfg)
	if errToken != nil {
		return "", ""
	}
	return accessToken, refreshToken

}
