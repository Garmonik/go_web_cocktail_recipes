package handlers

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/utils"
	"net/http"
)

// LoginPage templates/home.html
func LoginPage(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	cookie, err := r.Cookie("auth_token")
	if err == nil {
		if utils.IsValidToken(cookie.Value, cfg) {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}

	http.ServeFile(w, r, "./static/templates/facecontrol/login.html")
}
