package handler

import (
	"encoding/json"
	"net/http"
	"notes-api/internal/model"
	"notes-api/internal/service"
	"notes-api/internal/util"

	"github.com/gorilla/mux"
)

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
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
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
		respondError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	if err := h.authService.Register(&user); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, user)
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
		token := r.Header.Get("Authorization")
		if token == "" {
			respondError(w, http.StatusUnauthorized, "Требуется токен авторизации")
			return
		}

		if _, err := util.ParseJWT(token); err != nil {
			respondError(w, http.StatusUnauthorized, "Неверный токен")
			return
		}

		next.ServeHTTP(w, r)
	})
}
