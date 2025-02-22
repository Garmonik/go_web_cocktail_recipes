package middleware_auth

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"net/http"
)

func AuthMiddleware(next http.Handler, cfg *config.Config, db *db.DataBase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isValidAccessToken(r, cfg) {
			if isAuthPage(r.URL.Path) {
				http.Redirect(w, r, "/home/", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if _, valid := handleRefreshToken(w, r, cfg, db); valid {
			next.ServeHTTP(w, r)
			return
		}

		if !isAuthPage(r.URL.Path) {
			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
