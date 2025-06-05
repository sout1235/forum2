package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/auth-service/internal/delivery/http/middleware"
	"github.com/sout1235/forum2/backend/auth-service/internal/usecase"
)

type Router struct {
	authHandler *AuthHandler
}

func NewRouter(authUseCase *usecase.AuthUseCase, internalToken string) *Router {
	return &Router{
		authHandler: NewAuthHandler(authUseCase, internalToken),
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Публичные маршруты
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/verify", r.authHandler.VerifyToken)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}

	// Защищенные маршруты
	protected := router.Group("/api/auth")
	protected.Use(middleware.AuthMiddleware(r.authHandler.authUseCase, r.authHandler.internalToken))
	{
		protected.GET("/profile", r.authHandler.GetProfile)
		protected.GET("/user/:id", r.authHandler.GetUserByID)
	}

	return router
}
