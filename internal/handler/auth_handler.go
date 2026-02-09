package handler

import (
	"encoding/json"
	"net/http"

	"smartbooking/internal/logger"
	"smartbooking/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register handles user registration
// @Summary Register a new user
// @Description Create a new user account with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Register: Invalid request body - %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("Register: Attempting to register user with email: %s", req.Email)

	user, err := h.authService.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		logger.LogAuth("register", req.Email, false)
		logger.Error("Register: Failed to register user %s - %v", req.Email, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.LogAuth("register", req.Email, true)
	logger.Info("Register: Successfully registered user %s with ID %d", user.Email, user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles user login
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Login: Invalid request body - %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("Login: Attempting login for user: %s", req.Email)

	user, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		logger.LogAuth("login", req.Email, false)
		logger.Error("Login: Failed authentication for user %s - %v", req.Email, err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	logger.LogAuth("login", req.Email, true)
	logger.Info("Login: Successful login for user %s (ID: %d)", user.Email, user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
