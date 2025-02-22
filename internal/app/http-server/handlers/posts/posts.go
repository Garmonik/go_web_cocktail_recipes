package posts

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/config"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/app/db/models"
	"github.com/Garmonik/go_web_cocktail_recipes/internal/pkg/utils"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Posts struct {
	cfg      *config.Config
	router   chi.Router
	log      *slog.Logger
	DataBase *db.DataBase
}

func New(cfg *config.Config, r chi.Router, log *slog.Logger, dataBase *db.DataBase) *Posts {
	return &Posts{cfg: cfg, router: r, log: log, DataBase: dataBase}
}

func (p Posts) PostsListAPI(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	orderBy := r.URL.Query().Get("order_by")
	var order string
	switch orderBy {
	case "created":
		order = "created_at DESC"
	case "popular":
		order = "(SELECT COUNT(*) FROM likes WHERE likes.post_id = posts.id) DESC"
	default:
		order = "created_at DESC"
	}

	var posts []models.Post
	if err := p.DataBase.Db.
		Preload("Author.Avatar").
		Preload("Image").
		Order(order).
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		http.Error(w, "Error receiving posts", http.StatusInternalServerError)
		return
	}
	likedPosts := make(map[uint]bool)
	user, _ := utils.GetUserByToken(r, p.cfg, p.DataBase)
	if user != nil {
		likedPosts = utils.GetUserLikedPosts(p.DataBase.Db, user.ID)
	}

	response := make([]map[string]interface{}, len(posts))
	for i, post := range posts {
		imageData, err := os.ReadFile(post.Image.Path)
		if err != nil {
			http.Error(w, `{"error": "Failed to read image"}`, http.StatusBadRequest)
			return
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageData)
		avatarData, err := os.ReadFile(post.Author.Avatar.Path)
		if err != nil {
			http.Error(w, `{"error": "Failed to read avatar"}`, http.StatusBadRequest)
			return
		}
		avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
		response[i] = map[string]interface{}{
			"id":          post.ID,
			"name":        post.Name,
			"description": post.Description,
			"image":       imageBase64,
			"like":        likedPosts[post.ID],
			"author": map[string]string{
				"id":       strconv.Itoa(int(post.Author.ID)),
				"username": post.Author.Name,
				"avatar":   avatarBase64,
			},
		}
	}
	responseData := map[string]interface{}{
		"count":   len(posts),
		"content": response,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Error while encoding JSON", http.StatusInternalServerError)
	}
}

func (p Posts) PostsByIdAPI(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := p.DataBase.Db.
		Preload("Author.Avatar").
		Preload("Image").
		Where("id = ?", id).
		First(&post).Error; err != nil {
		http.Error(w, "Error receiving mail\n", http.StatusNotFound)
		return
	}

	user, _ := utils.GetUserByToken(r, p.cfg, p.DataBase)
	likedPost := false
	if user != nil {
		likedPost = utils.GetUserLikedPost(p.DataBase.Db, user.ID, post.ID)
	}
	imageData, err := os.ReadFile(post.Image.Path)
	if err != nil {
		http.Error(w, `{"error": "Failed to read image"}`, http.StatusBadRequest)
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	avatarData, err := os.ReadFile(post.Author.Avatar.Path)
	if err != nil {
		http.Error(w, `{"error": "Failed to read avatar"}`, http.StatusBadRequest)
		return
	}
	avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
	response := map[string]interface{}{
		"id":          post.ID,
		"description": post.Description,
		"image":       imageBase64,
		"like":        likedPost,
		"author": map[string]string{
			"id":       strconv.Itoa(int(post.Author.ID)),
			"username": post.Author.Name,
			"avatar":   avatarBase64,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON\n", http.StatusInternalServerError)
	}
}

func (p Posts) PostCreate(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserByToken(r, p.cfg, p.DataBase)
	if user == nil {
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}

	err := r.ParseMultipartForm(20 << 20) // 20 MB
	if err != nil {
		http.Error(w, "Error processing form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("Image")
	if err != nil {
		http.Error(w, "Failed to receive file", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		if err != nil {
			p.log.Error(err.Error())
		}
	}(file)

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if _, ok := utils.AllowedExtensions[ext]; !ok {
		http.Error(w, "Invalid file format", http.StatusBadRequest)
		return
	}

	uuid := utils.GenerateUUID()
	filename := fmt.Sprintf("%s%s", uuid, ext)
	filePath := filepath.Join("file/post", filename)
	fileURL := "file/post/" + filename

	if err := os.MkdirAll("file/post", os.ModePerm); err != nil {
		http.Error(w, "Error creating directory", http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			p.log.Error(err.Error())
		}
	}(dst)

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	name := r.FormValue("Name")
	description := r.FormValue("Description")

	image := models.Image{Path: fileURL}
	if err := p.DataBase.Db.Create(&image).Error; err != nil {
		http.Error(w, "Failed to create image record", http.StatusInternalServerError)
		return
	}

	post := models.Post{
		Name:        name,
		Description: description,
		AuthorID:    user.ID,
		ImageID:     image.ID,
	}
	if err := p.DataBase.Db.Create(&post).Error; err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	imageData, err := os.ReadFile(fileURL)
	if err != nil {
		http.Error(w, `{"error": "Failed to read image"}`, http.StatusBadRequest)
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imageData)
	avatarData, err := os.ReadFile(user.Avatar.Path)
	if err != nil {
		http.Error(w, `{"error": "Failed to read avatar"}`, http.StatusBadRequest)
		return
	}
	avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
	response := map[string]interface{}{
		"id":          post.ID,
		"description": post.Description,
		"image":       imageBase64,
		"like":        false,
		"author": map[string]string{
			"id":       strconv.Itoa(int(user.ID)),
			"username": user.Name,
			"avatar":   avatarBase64,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error while encoding JSON", http.StatusInternalServerError)
	}
}

func (p Posts) LikeAPI(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserByToken(r, p.cfg, p.DataBase)
	if user == nil {
		http.Error(w, "user not found", http.StatusForbidden)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := p.DataBase.Db.
		Select("id").
		Where("id = ?", id).
		First(&post).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var like models.Like
	if err := p.DataBase.Db.
		Where("post_id = ? AND author_id = ?", id, user.ID).
		First(&like).Error; err == nil {
		p.DataBase.Db.Delete(&like)
		w.WriteHeader(http.StatusNoContent)
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	newLike := models.Like{
		PostID:   uint(id),
		AuthorID: user.ID,
	}
	if err := p.DataBase.Db.Create(&newLike).Error; err != nil {
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (p Posts) PostsListByUserAPI(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	var posts []models.Post
	if err := p.DataBase.Db.
		Preload("Author.Avatar").
		Preload("Image").
		Where("author_id = ?", id).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&posts).Error; err != nil {
		http.Error(w, "Error receiving posts", http.StatusInternalServerError)
		return
	}
	likedPosts := make(map[uint]bool)
	user, _ := utils.GetUserByToken(r, p.cfg, p.DataBase)
	if user != nil {
		likedPosts = utils.GetUserLikedPosts(p.DataBase.Db, user.ID)
	}

	response := make([]map[string]interface{}, len(posts))
	for i, post := range posts {
		imageData, err := os.ReadFile(post.Image.Path)
		if err != nil {
			http.Error(w, `{"error": "Failed to read image"}`, http.StatusBadRequest)
			return
		}
		imageBase64 := base64.StdEncoding.EncodeToString(imageData)
		avatarData, err := os.ReadFile(post.Author.Avatar.Path)
		if err != nil {
			http.Error(w, `{"error": "Failed to read avatar"}`, http.StatusBadRequest)
			return
		}
		avatarBase64 := base64.StdEncoding.EncodeToString(avatarData)
		response[i] = map[string]interface{}{
			"id":          post.ID,
			"name":        post.Name,
			"description": post.Description,
			"image":       imageBase64,
			"like":        likedPosts[post.ID],
			"author": map[string]string{
				"id":       strconv.Itoa(int(post.Author.ID)),
				"username": post.Author.Name,
				"avatar":   avatarBase64,
			},
		}
	}
	responseData := map[string]interface{}{
		"count":   len(posts),
		"content": response,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Error while encoding JSON", http.StatusInternalServerError)
	}
}
