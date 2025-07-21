package api

import (
	"encoding/json"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"gocbt/internal/utils"
	"net/http"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	userService models.UserService
	jwtManager  *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService models.UserService, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	Password  string          `json:"password"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Role      models.UserRole `json:"role"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Register handles user registration
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeString(req.Email)
	req.FirstName = utils.SanitizeString(req.FirstName)
	req.LastName = utils.SanitizeString(req.LastName)

	// Validate required fields
	if utils.IsEmpty(req.Username) || utils.IsEmpty(req.Email) ||
		utils.IsEmpty(req.FirstName) || utils.IsEmpty(req.LastName) ||
		utils.IsEmpty(req.Password) {
		writeErrorResponse(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate email format
	if !utils.ValidateEmail(req.Email) {
		writeErrorResponse(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Validate username format
	if !utils.ValidateUsername(req.Username) {
		writeErrorResponse(w, "Username must be 3-50 characters, alphanumeric and underscores only", http.StatusBadRequest)
		return
	}

	// Check for SQL injection patterns
	if !utils.ValidateNoSQLInjection(req.Username) || !utils.ValidateNoSQLInjection(req.Email) ||
		!utils.ValidateNoSQLInjection(req.FirstName) || !utils.ValidateNoSQLInjection(req.LastName) {
		writeErrorResponse(w, "Invalid characters detected", http.StatusBadRequest)
		return
	}

	// Default role to student if not specified
	if req.Role == "" {
		req.Role = models.RoleStudent
	}

	// Validate role
	if !utils.ValidateRole(string(req.Role)) {
		writeErrorResponse(w, "Invalid role specified", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(req.Username, req.Email, req.Password,
		req.FirstName, req.LastName, req.Role)
	if err != nil {
		switch err {
		case auth.ErrUsernameExists:
			writeErrorResponse(w, "Username already exists", http.StatusConflict)
		case auth.ErrEmailExists:
			writeErrorResponse(w, "Email already exists", http.StatusConflict)
		case auth.ErrPasswordTooShort, auth.ErrPasswordTooLong:
			writeErrorResponse(w, err.Error(), http.StatusBadRequest)
		default:
			writeErrorResponse(w, "Registration failed", http.StatusInternalServerError)
		}
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		writeErrorResponse(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.Username = utils.SanitizeString(req.Username)

	// Validate required fields
	if utils.IsEmpty(req.Username) || utils.IsEmpty(req.Password) {
		writeErrorResponse(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Check for SQL injection patterns
	if !utils.ValidateNoSQLInjection(req.Username) {
		writeErrorResponse(w, "Invalid characters detected", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			writeErrorResponse(w, "Invalid username or password", http.StatusUnauthorized)
		case auth.ErrUserNotActive:
			writeErrorResponse(w, "Account is not active", http.StatusForbidden)
		default:
			writeErrorResponse(w, "Login failed", http.StatusInternalServerError)
		}
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user)
	if err != nil {
		writeErrorResponse(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Profile handles getting user profile
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetProfile(userID)
	if err != nil {
		switch err {
		case auth.ErrUserNotFound:
			writeErrorResponse(w, "User not found", http.StatusNotFound)
		default:
			writeErrorResponse(w, "Failed to get profile", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 {
		writeErrorResponse(w, "Authorization header required", http.StatusUnauthorized)
		return
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix

	newToken, err := h.jwtManager.RefreshToken(tokenString)
	if err != nil {
		writeErrorResponse(w, "Token refresh failed", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"token": newToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// writeErrorResponse writes an error response
func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}
