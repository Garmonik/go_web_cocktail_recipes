package render_page

import (
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"strings"
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

func (r *Renderer) MyUserPage(w http.ResponseWriter, req *http.Request) {
	r.RenderPage(w, req, "./static/templates/users/my_user.html")
}

func (r *Renderer) PostsList(w http.ResponseWriter, req *http.Request) {
	r.RenderPage(w, req, "./static/templates/posts/post.html")
}

func (r *Renderer) SomeUserPage(w http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "id")
	html, err := os.ReadFile("./static/templates/users/user.html")
	if err != nil {
		http.Error(w, "Error with render page", http.StatusInternalServerError)
		return
	}
	modifiedHTML := strings.Replace(string(html), "<body>", fmt.Sprintf("<body data-user-id=\"%s\">", userID), 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(modifiedHTML))
}
