package middleware_auth

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"net/http"
	"time"
)

// AuthMiddleware check auth_token in cookie
func AuthMiddleware(next http.Handler, cfg *config.Config, db *db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessCookie, errAccess := r.Cookie("access_token")
		refreshCookie, errRefresh := r.Cookie("refresh_token")

		validAccess := false
		if errAccess == nil && accessCookie != nil && accessCookie.Value != "" {
			validAccess = utils.IsValidToken(accessCookie.Value, cfg)
		}
		if validAccess {
			if r.URL.Path == "/login/" || r.URL.Path == "/register/" || r.URL.Path == "/api/login" || r.URL.Path == "/api/register" {
				http.Redirect(w, r, "/home/", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if errRefresh == nil && refreshCookie != nil && refreshCookie.Value != "" {
			accessToken, refreshToken := utils.RefreshAccessToken(refreshCookie.Value, cfg, db)
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    accessToken,
				Path:     "/",
				Domain:   r.Host,
				Expires:  time.Now().Add(15 * time.Minute),
				HttpOnly: true,
				Secure:   refreshCookie.Secure,
				SameSite: refreshCookie.SameSite,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    refreshToken,
				Path:     "/",
				Domain:   r.Host,
				Expires:  time.Now().Add(30 * 24 * time.Hour),
				HttpOnly: true,
				Secure:   refreshCookie.Secure,
				SameSite: refreshCookie.SameSite,
			})
			if accessToken == "" || refreshToken == "" {
			}
			next.ServeHTTP(w, r)
			return

		}

		if r.URL.Path != "/login/" && r.URL.Path != "/register/" {
			if r.URL.Path == "/api/login/" || r.URL.Path == "/api/register/" {
				next.ServeHTTP(w, r)
				return
			}
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
