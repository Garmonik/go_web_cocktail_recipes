package render_page

import (
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type Renderer struct {
	cfg    *config.Config
	router chi.Router
	log    *slog.Logger
}

func New(cfg *config.Config, r chi.Router, log *slog.Logger) *Renderer {
	return &Renderer{cfg: cfg, router: r, log: log}
}

// RenderPage RenderPageWithoutAuth checks authorization and renders the page
func (r *Renderer) RenderPage(w http.ResponseWriter, req *http.Request, templatePath string) {
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
