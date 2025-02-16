package render_page

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"net/http"
	"time"
)

type Renderer struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Renderer {
	return &Renderer{cfg: cfg}
}

// IsAuthenticated checks if the user is authorized
func (r *Renderer) IsAuthenticated(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("auth_token")
	if err == nil && utils.IsValidToken(cookie.Value, r.cfg) {
		return true
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	return false
}

// RenderPageWithoutAuth checks authorization and renders the page
func (r *Renderer) RenderPageWithoutAuth(w http.ResponseWriter, req *http.Request, templatePath string) {
	if r.IsAuthenticated(w, req) {
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}
	http.ServeFile(w, req, templatePath)
}

// LoginPage renders the login page
func (r *Renderer) LoginPage(w http.ResponseWriter, req *http.Request) {
	r.RenderPageWithoutAuth(w, req, "./static/templates/facecontrol/login.html")
}

// RegisterPage renders the registration page
func (r *Renderer) RegisterPage(w http.ResponseWriter, req *http.Request) {
	r.RenderPageWithoutAuth(w, req, "./static/templates/facecontrol/register.html")
}
