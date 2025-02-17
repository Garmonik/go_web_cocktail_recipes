package render_page

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/http-server/handlers/facecontrol"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"time"
)

type Renderer struct {
	cfg    *config.Config
	router *chi.Mux
	log    *slog.Logger
}

func New(cfg *config.Config, r *chi.Mux, log *slog.Logger) *Renderer {
	return &Renderer{cfg: cfg, router: r, log: log}
}

// IsAuthenticated checks if the user is authorized
func (r *Renderer) IsAuthenticated(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		r.log.Info("There is no auth_token cookie, we redirect to /login/")
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})
		return false
	}
	if !facecontrol.IsValidToken(cookie.Value, r.cfg) {
		r.log.Info("Invalid token, delete the cookie and redirect to /login/")
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
		})
		return false
	}

	r.log.Info("User is authenticated")
	return true
}

// RenderPage RenderPageWithoutAuth checks authorization and renders the page
func (r *Renderer) RenderPage(w http.ResponseWriter, req *http.Request, templatePath string) {
	isAuth := r.IsAuthenticated(w, req)

	// If the user is already authorized and tries to log in to /login/ or /register/, we send to /home/
	if isAuth && (req.URL.Path == "/login/" || req.URL.Path == "/register/") {
		r.log.Info("The user is already authorized, redirect to /home/")
		http.Redirect(w, req, "/home/", http.StatusFound)
		return
	}

	// If the user is not authorized and NOT on /login/, redirect to /login/
	if !isAuth && req.URL.Path != "/login/" {
		r.log.Info("User is not authenticated, redirect to /login/")
		http.Redirect(w, req, "/login/", http.StatusFound)
		return
	}

	r.log.Info("Рендеринг страницы: " + templatePath)
	http.ServeFile(w, req, templatePath)
}

// LoginPage renders the login page
func (r *Renderer) LoginPage(w http.ResponseWriter, req *http.Request) {
	r.RenderPage(w, req, "./static/templates/facecontrol/login.html")
}

// RegisterPage renders the registration page
func (r *Renderer) RegisterPage(w http.ResponseWriter, req *http.Request) {
	r.RenderPage(w, req, "./static/templates/facecontrol/register.html")
}

func (r *Renderer) HomePage(w http.ResponseWriter, req *http.Request) {
	r.RenderPage(w, req, "./static/templates/base/home.html")
}
