package middleware_auth

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"net/http"
	"time"
)

func isValidAccessToken(r *http.Request, cfg *config.Config) bool {
	accessCookie, err := r.Cookie("access_token")
	return err == nil && accessCookie != nil && accessCookie.Value != "" && utils.IsValidToken(accessCookie.Value, cfg)
}

func handleRefreshToken(w http.ResponseWriter, r *http.Request, cfg *config.Config, db *db.DataBase) (string, bool) {
	refreshCookie, err := r.Cookie("refresh_token")
	if err != nil || refreshCookie == nil || refreshCookie.Value == "" {
		return "", false
	}

	accessToken, refreshToken := utils.RefreshAccessToken(refreshCookie.Value, cfg, db)
	if accessToken == "" || refreshToken == "" {
		clearCookies(w, r)
		return "", false
	}

	setAuthCookies(w, r, accessToken, refreshToken, refreshCookie)
	return accessToken, true
}

func clearCookies(w http.ResponseWriter, r *http.Request) {
	for _, name := range []string{"access_token", "refresh_token"} {
		http.SetCookie(w, &http.Cookie{
			Name:    name,
			Value:   "",
			Path:    "/",
			Domain:  r.Host,
			Expires: time.Unix(0, 0),
			MaxAge:  -1,
		})
	}
}

func setAuthCookies(w http.ResponseWriter, r *http.Request, accessToken, refreshToken string, refreshCookie *http.Cookie) {
	domain := getDomain(r.Host)

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   refreshCookie.Secure,
		SameSite: refreshCookie.SameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		Domain:   domain,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   refreshCookie.Secure,
		SameSite: refreshCookie.SameSite,
	})
}

func getDomain(host string) string {
	if host == "localhost:8000" {
		return ""
	}
	return host
}

func isAuthPage(path string) bool {
	return path == "/login/" || path == "/register/" || path == "/api/login/" || path == "/api/register/"
}
