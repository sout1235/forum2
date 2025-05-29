package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/services/auth/models"
	"backend/services/auth/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Repo      *repository.UserRepository
	jwtSecret string
	redis     *redis.Client
}

func NewAuthHandler(repo *repository.UserRepository, jwtSecret string, redis *redis.Client) *AuthHandler {
	return &AuthHandler{
		Repo:      repo,
		jwtSecret: jwtSecret,
		redis:     redis,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ParseJWT parses and validates a JWT token
func (h *AuthHandler) ParseJWT(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	// Create user
	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     "user",
	}

	createdUser, err := h.Repo.CreateUser(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	c.JSON(201, createdUser)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	// Get user by email
	user, err := h.Repo.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(401, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := h.generateTokens(strconv.FormatInt(user.ID, 10), user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при генерации токенов"})
		return
	}

	// Store refresh token in Redis
	err = h.redis.Set(c.Request.Context(), refreshToken, user.ID, 24*time.Hour).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при сохранении токена"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	// Get user ID from Redis
	userIDStr, err := h.redis.Get(c.Request.Context(), input.RefreshToken).Result()
	if err != nil {
		c.JSON(401, gin.H{"error": "Недействительный токен обновления"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при обработке ID пользователя"})
		return
	}

	// Get user from database
	user, err := h.Repo.GetUserByID(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении пользователя"})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := h.generateTokens(strconv.FormatInt(user.ID, 10), user.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при генерации токенов"})
		return
	}

	// Delete old refresh token
	h.redis.Del(c.Request.Context(), input.RefreshToken)

	// Store new refresh token
	err = h.redis.Set(c.Request.Context(), refreshToken, user.ID, 24*time.Hour).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при сохранении токена"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверные данные"})
		return
	}

	// Delete refresh token from Redis
	err := h.redis.Del(c.Request.Context(), input.RefreshToken).Err()
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при удалении токена"})
		return
	}

	c.JSON(200, gin.H{"message": "Успешный выход"})
}

func (h *AuthHandler) generateTokens(userID string, role string) (string, string, error) {
	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (h *AuthHandler) Verify(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Некорректный формат токена"})
		return
	}

	token := tokenParts[1]
	claims, err := h.ParseJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный токен"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": claims["user_id"],
		"role":    claims["role"],
	})
}
