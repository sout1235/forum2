package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/services/forum/config"
	"backend/services/forum/grpc"
	"backend/services/forum/handlers"
	"backend/services/forum/middleware"
	"backend/services/forum/repository"

	"github.com/sout1235/forumski/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Load configuration
	cfg := config.NewConfig()

	// Initialize database
	db, err := repository.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize gRPC client
	authClient, err := grpc.NewAuthClient(cfg.AuthGRPCURL)
	if err != nil {
		logger.Fatal("Failed to connect to auth service", zap.Error(err))
	}
	defer authClient.Close()

	// Initialize repositories
	topicRepo := repository.NewTopicRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)

	// Initialize handlers
	topicHandler := handlers.NewTopicHandler(topicRepo)
	commentHandler := handlers.NewCommentHandler(commentRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryRepo)
	tagHandler := handlers.NewTagHandler(tagRepo)

	// Initialize Gin router
	r := gin.Default()

	// Add health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Public routes
	public := r.Group("/api")
	{
		// Topic routes
		topics := public.Group("/topics")
		{
			topics.GET("", topicHandler.GetAllTopics)
			topics.GET("/:id", topicHandler.GetTopicByID)
		}

		// Comment routes
		comments := public.Group("/comments")
		{
			comments.GET("/topic/:topic_id", commentHandler.GetAllCommentsByTopic)
		}

		// Category routes
		categories := public.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllCategories)
			categories.GET("/:id", categoryHandler.GetCategoryByID)
		}

		// Tag routes
		tags := public.Group("/tags")
		{
			tags.GET("", tagHandler.GetAllTags)
			tags.GET("/:id", tagHandler.GetTagByID)
		}
	}

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(authClient))
	{
		// Topic routes
		topics := protected.Group("/topics")
		{
			topics.POST("", topicHandler.CreateTopic)
			topics.PUT("/:id", topicHandler.UpdateTopic)
			topics.DELETE("/:id", topicHandler.DeleteTopic)
		}

		// Comment routes
		comments := protected.Group("/comments")
		{
			comments.POST("", commentHandler.CreateComment)
			comments.PUT("/:id", commentHandler.UpdateComment)
			comments.DELETE("/:id", commentHandler.DeleteComment)
		}

		// Category routes
		categories := protected.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.PUT("/:id", categoryHandler.UpdateCategory)
			categories.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		// Tag routes
		tags := protected.Group("/tags")
		{
			tags.POST("", tagHandler.CreateTag)
			tags.PUT("/:id", tagHandler.UpdateTag)
			tags.DELETE("/:id", tagHandler.DeleteTag)
		}
	}

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		logger.Info("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
