package utils

import (
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//{
//	"exp": time
//	"sub": id
//	"email": email
//	"urn": cript(username)
//  "iat": time.Now()
//}

// IsValidToken check token
func IsValidToken(tokenString string, cfg *config.Config) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return cfg.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false
		}
		if time.Now().Unix() > int64(exp) {
			return false
		}
		return true
	}
	return false
}
