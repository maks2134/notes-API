package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notes-api/internal/model"
	"notes-api/internal/service"
	"notes-api/internal/util"

	"github.com/gorilla/mux"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
}

func decodeJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		fmt.Println(err)
	}
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Создает нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.User true "Данные пользователя"
// @Success 201 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := decodeJSON(r, &user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Валидация пароля
	if user.Password == "" {
		respondError(w, http.StatusBadRequest, "Password cannot be empty")
		return
	}

	log.Printf("Registering user: %s", user.Username) // Логируем только имя
	if err := h.authService.Register(&user); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Не возвращаем пароль в ответе
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
	})
}

// Login godoc
// @Summary Аутентификация пользователя
// @Description Вход пользователя в систему
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body model.LoginRequest true "Учетные данные"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	token, err := h.authService.Login(&req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Неверные учетные данные")
		return
	}

	respondJSON(w, http.StatusOK, model.LoginResponse{Token: token})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Требуется токен авторизации")
			return
		}

		userID, err := util.ParseJWT(authHeader)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "Неверный или просроченный токен")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
