package handlers

import (
	"encoding/json"
	"main_service/internal/grpcclient"
	"main_service/internal/storage"
	"net/http"
)

type AuthHandler struct {
	storage *storage.PostgresStorage
}

func NewAuthHandler(storage *storage.PostgresStorage) *AuthHandler {
	return &AuthHandler{storage: storage}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	status, err := grpcclient.RegisterUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&creds)

	token, err := grpcclient.LoginUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Ошибка входа", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
