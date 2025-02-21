package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

type Json400 struct {
	Error string `json:"Error,omitempty"`
}

type Json200Facecontrol struct {
	Status string `json:"Status,omitempty"`
	ID     uint   `json:"ID,omitempty"`
}

func JsonResponse400(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	err := json.NewEncoder(w).Encode(Json400{Error: text})
	if err != nil {
		os.Exit(1)
	}
}

func JsonResponse200Facecontrol(UserId uint, statusCode int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(Json200Facecontrol{Status: "Success", ID: UserId})
	if err != nil {
		os.Exit(1)
	}
}

func GenerateUUID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "default-uuid"
	}
	return hex.EncodeToString(bytes)
}

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}
