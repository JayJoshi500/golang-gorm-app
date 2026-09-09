package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/JayJoshi500/golang-gorm-app/internal/services"
	"github.com/JayJoshi500/golang-gorm-app/pkg/response"
)

// AuthHandler adapts HTTP requests to the AuthService use cases.
type AuthHandler struct {
	service services.AuthService
	logger  *slog.Logger
}

func NewAuthHandler(service services.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{service: service, logger: logger}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || len(req.Password) < 8 {
		response.Error(w, http.StatusUnprocessableEntity, "name, email, and password (min 8 chars) are required")
		return
	}

	user, err := h.service.Register(req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			response.Error(w, http.StatusConflict, err.Error())
			return
		}
		h.logger.Error("register failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "could not register user")
		return
	}

	response.Success(w, http.StatusCreated, user)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredential) {
			response.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		h.logger.Error("login failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "could not log in")
		return
	}

	response.Success(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// Logout handles POST /api/v1/auth/logout. Expects "Authorization: Bearer <token>".
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "missing bearer token")
		return
	}

	if err := h.service.Logout(token); err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func bearerToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
