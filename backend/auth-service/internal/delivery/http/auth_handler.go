package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/auth-service/internal/entity"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authUseCase   *usecase.AuthUseCase
	internalToken string
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase, internalToken string) *AuthHandler {
	return &AuthHandler{
		authUseCase:   authUseCase,
		internalToken: internalToken,
	}
}

// RegisterRequest represents the request body for user registration
// @Description Registration request with username, password and email
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
}

// LoginRequest represents the request body for user login
// @Description Login request with username and password
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// AuthResponse represents the authentication response with tokens and user info
// @Description Authentication response containing access token, refresh token and user information
type AuthResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         *entity.User `json:"user"`
}

// RefreshTokenRequest represents the request body for token refresh
// @Description Request to refresh the access token using a refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse represents an error response
// @Description Error response containing error message
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}

// RegisterResponse represents the response for successful registration
// @Description Response after successful user registration
type RegisterResponse struct {
	Message string `json:"message" example:"User registered successfully"`
}

// LoginResponse represents the response for successful login
// @Description Response after successful login
type LoginResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User         *entity.User `json:"user"`
}

// ProfileResponse represents the user profile information
// @Description User profile information
type ProfileResponse struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"johndoe"`
	Email     string `json:"email" example:"john@example.com"`
	Role      string `json:"role" example:"user"`
	CreatedAt string `json:"created_at" example:"2024-03-15T10:00:00Z"`
}

// VerifyTokenResponse represents the response for token verification
// @Description Response after token verification
type VerifyTokenResponse struct {
	Valid    bool   `json:"valid" example:"true"`
	UserID   int64  `json:"user_id,omitempty" example:"1"`
	Username string `json:"username,omitempty" example:"johndoe"`
}

// VerifyTokenRequest represents the request body for token verification
// @Description Request to verify JWT token
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse represents the response for token refresh
// @Description Response after successful token refresh
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// @Summary Register new user
// @Description Register a new user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register credentials"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := entity.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Email:        req.Email,
		Role:         "user", // По умолчанию роль "user"
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.authUseCase.Register(c.Request.Context(), &user); err != nil {
		log.Printf("Register error: %v", err)
		errMsg := err.Error()
		if errMsg == "username already exists" || errMsg == "email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		return
	}

	// После успешной регистрации сразу возвращаем токены и данные пользователя
	tokens, err := h.authUseCase.GenerateTokenPair(&user)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         &user,
	})
}

// @Summary Login user
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, user, err := h.authUseCase.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	})
}

// @Summary Get user profile
// @Description Get current user profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// @Summary Verify token
// @Description Verify JWT token validity
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyTokenRequest true "Token to verify"
// @Success 200 {object} VerifyTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/verify [post]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	log.Printf("Received token verification request")
	log.Printf("Request headers: %v", c.Request.Header)

	// Читаем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read request body"})
		return
	}
	log.Printf("Request body: %s", string(body))

	var tokenReq struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(body, &tokenReq); err != nil {
		log.Printf("Failed to unmarshal token: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid token format"})
		return
	}

	log.Printf("Token from request: %s", tokenReq.Token)
	log.Printf("Validating token")
	userID, err := h.authUseCase.ValidateToken(tokenReq.Token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	log.Printf("Token validated for user ID: %s", userID)
	user, err := h.authUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	log.Printf("Token verification successful for user: %s", user.Username)
	c.JSON(http.StatusOK, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
	})
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.authUseCase.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// @Summary Refresh token
// @Description Refresh JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token"
// @Success 200 {object} RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.authUseCase.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}
