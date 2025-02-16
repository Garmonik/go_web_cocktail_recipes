package middleware_auth

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"net/http"
)

// AuthMiddleware check auth_token in cookie
func AuthMiddleware(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil || cookie.Value == "" || !utils.IsValidToken(cookie.Value, cfg) {
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
