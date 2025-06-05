package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthConfig struct {
	AuthServiceURL string
}

func NewAuthConfig(authServiceURL string) *AuthConfig {
	return &AuthConfig{
		AuthServiceURL: authServiceURL,
	}
}

func (c *AuthConfig) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		log.Printf("Auth header: %s", authHeader)
		log.Printf("All request headers: %v", ctx.Request.Header)

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			ctx.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token is required"})
			ctx.Abort()
			return
		}

		log.Printf("Verifying token with auth service: %s", c.AuthServiceURL)
		// Verify token with auth service
		tokenReq := struct {
			Token string `json:"token"`
		}{
			Token: token,
		}
		reqBody, err := json.Marshal(tokenReq)
		if err != nil {
			log.Printf("Failed to marshal token request: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification request"})
			ctx.Abort()
			return
		}

		log.Printf("Sending verification request to auth service with body: %s", string(reqBody))
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/auth/verify", c.AuthServiceURL), bytes.NewBuffer(reqBody))
		if err != nil {
			log.Printf("Failed to create verification request: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create verification request"})
			ctx.Abort()
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed to verify token: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify token"})
			ctx.Abort()
			return
		}
		defer resp.Body.Close()

		log.Printf("Auth service response status: %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Auth service response body: %s", string(body))

		if resp.StatusCode != http.StatusOK {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			ctx.Abort()
			return
		}

		var userData struct {
			UserID   string `json:"user_id"`
			Username string `json:"username"`
		}
		if err := json.Unmarshal(body, &userData); err != nil {
			log.Printf("Failed to decode user data: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode user data"})
			ctx.Abort()
			return
		}

		// Convert user ID from string to int64
		userID, err := strconv.ParseInt(userData.UserID, 10, 64)
		if err != nil {
			log.Printf("Failed to parse user ID: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user ID format"})
			ctx.Abort()
			return
		}

		log.Printf("Token verified successfully for user: %s", userData.Username)
		// Set user data in context
		ctx.Set("user_id", userID)
		ctx.Set("username", userData.Username)
		ctx.Next()
	}
}
