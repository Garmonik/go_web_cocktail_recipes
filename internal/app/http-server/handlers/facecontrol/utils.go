package facecontrol

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
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

func CheckUserByEmail(email string, facecontrol *Facecontrol) (*models.User, string) {
	var user models.User
	if err := facecontrol.dataBase.Db.Select("id", "name", "email", "password").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return &models.User{}, "user not found"
	}
	return &user, ""
}

func CheckUserByName(name string, facecontrol *Facecontrol) (*models.User, string) {
	var user models.User
	if err := facecontrol.dataBase.Db.Select("id", "name", "email", "password").
		Where("name = ?", name).
		First(&user).Error; err != nil {
		return &models.User{}, "user not found"
	}
	return &user, ""
}

// CheckUser check user exist in db
func CheckUser(email string, password string, facecontrol *Facecontrol) (*models.User, error) {
	user, err := CheckUserByEmail(email, facecontrol)
	if err != "" {
		return nil, fmt.Errorf(err)
	}

	ok := CheckPassword(user, password)

	if ok {
		return user, nil
	}
	return nil, fmt.Errorf("the password will be entered incorrectly")
}

// GenerateToken creates an ES256 token
func GenerateToken(user *models.User, facecontrol *Facecontrol) (string, error) {
	now := time.Now()
	expirationTime := now.Add(30 * 24 * time.Hour).Unix() // +1 month

	claims := jwt.MapClaims{
		"exp":   expirationTime,
		"sub":   user.ID,
		"email": user.Email,
		"urn":   user.Name,
		"iat":   now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// Checking that privateKey is not nil and of the correct type
	if facecontrol.cfg.JWT.PrivateKey == nil {
		return "", fmt.Errorf("private key is nil")
	}

	privateKey := facecontrol.cfg.JWT.PrivateKey

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
	if err := facecontrol.dataBase.Db.Create(&avatar).Error; err != nil {
		return nil, fmt.Errorf("failed to create avatar")
	}
	// Create a user in the database
	user := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Bio:      "",
		AvatarID: avatar.ID,
		Avatar:   avatar,
	}

	if err := facecontrol.dataBase.Db.Create(&user).Error; err != nil {
		facecontrol.dataBase.Db.Delete(&avatar)
		return nil, fmt.Errorf("failed to create user")
	}

	return &user, nil
}
