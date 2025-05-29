package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация роутера
	r := gin.Default()

	// Настройка CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// Группа API маршрутов
	api := r.Group("/api")
	{
		// Аутентификация
		auth := api.Group("/auth")
		{
			auth.POST("/register", handleRegister)
			auth.POST("/login", handleLogin)
			auth.POST("/logout", handleLogout)
		}

		// Темы форума
		topics := api.Group("/topics")
		{
			topics.GET("/", getTopics)
			topics.POST("/", createTopic)
			topics.GET("/:id", getTopic)
			topics.PUT("/:id", updateTopic)
			topics.DELETE("/:id", deleteTopic)
		}

		// Комментарии
		comments := api.Group("/comments")
		{
			comments.GET("/topic/:topicId", getComments)
			comments.POST("/", createComment)
			comments.PUT("/:id", updateComment)
			comments.DELETE("/:id", deleteComment)
		}

		// Пользователи
		users := api.Group("/users")
		{
			users.GET("/profile", getProfile)
			users.PUT("/profile", updateProfile)
		}
	}

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Заглушки для обработчиков
func handleRegister(c *gin.Context) { c.JSON(200, gin.H{"message": "Register endpoint"}) }
func handleLogin(c *gin.Context)    { c.JSON(200, gin.H{"message": "Login endpoint"}) }
func handleLogout(c *gin.Context)   { c.JSON(200, gin.H{"message": "Logout endpoint"}) }
func getTopics(c *gin.Context)      { c.JSON(200, gin.H{"message": "Get topics endpoint"}) }
func createTopic(c *gin.Context)    { c.JSON(200, gin.H{"message": "Create topic endpoint"}) }
func getTopic(c *gin.Context)       { c.JSON(200, gin.H{"message": "Get topic endpoint"}) }
func updateTopic(c *gin.Context)    { c.JSON(200, gin.H{"message": "Update topic endpoint"}) }
func deleteTopic(c *gin.Context)    { c.JSON(200, gin.H{"message": "Delete topic endpoint"}) }
func getComments(c *gin.Context)    { c.JSON(200, gin.H{"message": "Get comments endpoint"}) }
func createComment(c *gin.Context)  { c.JSON(200, gin.H{"message": "Create comment endpoint"}) }
func updateComment(c *gin.Context)  { c.JSON(200, gin.H{"message": "Update comment endpoint"}) }
func deleteComment(c *gin.Context)  { c.JSON(200, gin.H{"message": "Delete comment endpoint"}) }
func getProfile(c *gin.Context)     { c.JSON(200, gin.H{"message": "Get profile endpoint"}) }
func updateProfile(c *gin.Context)  { c.JSON(200, gin.H{"message": "Update profile endpoint"}) }
