package httpDelivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sout1235/forum2/backend/forum-service/internal/delivery/http/middleware"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/service"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// ErrorResponse represents an error response
// @Description Error response containing error message
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// SuccessResponse represents a success response
// @Description Success response containing message
type SuccessResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// @title Forum Service API
// @version 1.0
// @description API for forum service with topics, comments and chat functionality
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

type Router struct {
	engine         *gin.Engine
	topicUseCase   service.TopicService
	commentUseCase usecase.CommentUseCase
	userRepo       repository.UserRepository
	chatRepo       repository.ChatRepository
	port           string
	upgrader       websocket.Upgrader
	clients        map[*websocket.Conn]string // map[connection]username
	authConfig     *middleware.AuthConfig
	authURL        string
	logger         *zap.Logger
}

type WSMessage struct {
	Type                 string          `json:"type"`
	Token                string          `json:"token,omitempty"`
	Content              string          `json:"content,omitempty"`
	Author               string          `json:"author,omitempty"`
	Data                 json.RawMessage `json:"data,omitempty"`
	ID                   string          `json:"id,omitempty"`
	Timestamp            int64           `json:"timestamp,omitempty"`
	LastMessageTimestamp int64           `json:"lastMessageTimestamp,omitempty"`
}

func NewRouter(
	topicUseCase service.TopicService,
	commentUseCase usecase.CommentUseCase,
	userRepo repository.UserRepository,
	chatRepo repository.ChatRepository,
	port string,
	authConfig *middleware.AuthConfig,
) *Router {
	// Инициализируем логгер
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Sec-WebSocket-Protocol"},
		ExposeHeaders:    []string{"Content-Length", "Upgrade", "Connection", "Sec-WebSocket-Protocol"},
		AllowCredentials: true,
	}))

	// Создаем обработчики
	topicHandler := NewTopicHandler(topicUseCase, userRepo)
	commentHandler := NewCommentHandler(commentUseCase, userRepo)

	// Создаем middleware для аутентификации
	authMiddleware := authConfig

	// Инициализируем WebSocket
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			allowedOrigins := []string{
				"http://localhost:3000",
				"http://localhost:5173",
				"http://127.0.0.1:5173",
			}
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		},
		EnableCompression: true,
	}

	r := &Router{
		engine:         router,
		topicUseCase:   topicUseCase,
		commentUseCase: commentUseCase,
		userRepo:       userRepo,
		chatRepo:       chatRepo,
		port:           port,
		upgrader:       upgrader,
		clients:        make(map[*websocket.Conn]string),
		authConfig:     authConfig,
		authURL:        authConfig.AuthServiceURL,
		logger:         logger,
	}

	// Группа маршрутов API v1
	v1 := router.Group("/api/v1")
	{
		// Маршруты для тем
		topics := v1.Group("/topics")
		{
			topics.GET("", topicHandler.GetAllTopics)
			topics.GET("/:id", topicHandler.GetTopic)
			topics.GET("/:id/comments", commentHandler.GetAllCommentsByTopic)
			topics.POST("", authMiddleware.AuthMiddleware(), topicHandler.CreateTopic)
			topics.PUT("/:id", authMiddleware.AuthMiddleware(), topicHandler.UpdateTopic)
			topics.DELETE("/:id", authMiddleware.AuthMiddleware(), topicHandler.DeleteTopic)
		}

		// Маршруты для комментариев
		comments := v1.Group("/comments")
		{
			comments.GET("/:id", commentHandler.GetComment)
		}

		// Маршруты для комментариев к теме
		topicComments := v1.Group("/topics/:id/comments")
		{
			topicComments.POST("", authMiddleware.AuthMiddleware(), commentHandler.CreateComment)
			topicComments.DELETE("/:commentId", authMiddleware.AuthMiddleware(), commentHandler.DeleteComment)
		}

		// Маршруты для чата
		chat := v1.Group("/chat")
		{
			chat.GET("/messages", authMiddleware.AuthMiddleware(), func(c *gin.Context) {
				messages, err := r.chatRepo.GetRecentMessages(c.Request.Context(), 50)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, messages)
			})
		}
	}

	// WebSocket маршрут
	router.GET("/ws", r.handleWebSocket)

	// Инициализируем Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (r *Router) handleWebSocket(c *gin.Context) {
	// Upgrade HTTP connection to WebSocket
	conn, err := r.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		r.logger.Error("Error upgrading to WebSocket",
			zap.Error(err),
			zap.String("remote_addr", c.Request.RemoteAddr))
		return
	}
	defer conn.Close()

	// Set read deadline
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Initialize lastMessageTimestamp
	var lastMessageTimestamp int64

	// Load recent messages when client connects
	recentMessages, err := r.chatRepo.GetRecentMessages(c.Request.Context(), 50)
	if err != nil {
		r.logger.Error("Error loading recent messages",
			zap.Error(err))
	} else {
		// Send only new messages
		for i := len(recentMessages) - 1; i >= 0; i-- {
			msg := recentMessages[i]
			msgTime := msg.CreatedAt.Unix()

			// Skip old messages
			if msgTime <= lastMessageTimestamp {
				continue
			}

			response := WSMessage{
				Type:      "message",
				Content:   msg.Content,
				Author:    msg.AuthorUsername,
				ID:        fmt.Sprintf("%d:%d", msg.AuthorID, msgTime),
				Timestamp: msgTime,
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				r.logger.Error("Error marshaling message",
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				r.logger.Error("Error sending message",
					zap.Error(err))
			}
		}
	}

	// Message handling loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				r.logger.Error("Error reading message",
					zap.Error(err),
					zap.String("remote_addr", c.Request.RemoteAddr))
			}
			break
		}

		r.logger.Debug("Received raw message",
			zap.String("remote_addr", c.Request.RemoteAddr),
			zap.String("message", string(message)))

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			r.logger.Error("Error unmarshaling message",
				zap.Error(err),
				zap.String("remote_addr", c.Request.RemoteAddr))
			continue
		}

		// Update lastMessageTimestamp if provided
		if wsMsg.LastMessageTimestamp > 0 {
			lastMessageTimestamp = wsMsg.LastMessageTimestamp
		}

		r.logger.Debug("Parsed message",
			zap.String("type", wsMsg.Type),
			zap.String("content", wsMsg.Content))

		// Обработка пинг-сообщений
		if wsMsg.Type == "ping" {
			response := WSMessage{
				Type: "pong",
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				r.logger.Error("Error marshaling pong response",
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
				r.logger.Error("Error sending pong response",
					zap.Error(err))
			}
			continue
		}

		// Обработка авторизации
		if wsMsg.Type == "auth" {
			r.logger.Info("Processing auth message",
				zap.String("remote_addr", c.Request.RemoteAddr))

			// Проверяем токен
			tokenReq := struct {
				Token string `json:"token"`
			}{
				Token: wsMsg.Token,
			}
			reqBody, err := json.Marshal(tokenReq)
			if err != nil {
				r.logger.Error("Error marshaling token request",
					zap.Error(err))
				continue
			}

			req, err := http.NewRequest("POST", r.authURL+"/api/auth/verify", bytes.NewBuffer(reqBody))
			if err != nil {
				r.logger.Error("Error creating verification request",
					zap.Error(err))
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				r.logger.Error("Error verifying token",
					zap.Error(err))
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				r.logger.Error("Invalid token",
					zap.String("remote_addr", c.Request.RemoteAddr))
				errorResponse := WSMessage{
					Type:    "error",
					Content: "Invalid token",
				}
				errorBytes, err := json.Marshal(errorResponse)
				if err != nil {
					r.logger.Error("Error marshaling error response",
						zap.Error(err))
					continue
				}
				if err := conn.WriteMessage(messageType, errorBytes); err != nil {
					r.logger.Error("Error sending error response",
						zap.Error(err))
				}
				continue
			}

			// Получаем информацию о пользователе
			var userData struct {
				UserID   string `json:"user_id"`
				Username string `json:"username"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
				r.logger.Error("Error decoding user data",
					zap.Error(err))
				continue
			}

			// Парсим ID пользователя
			userID, err := strconv.ParseInt(userData.UserID, 10, 64)
			if err != nil {
				r.logger.Error("Error parsing user ID",
					zap.Error(err))
				continue
			}

			// Сохраняем информацию о пользователе
			r.clients[conn] = fmt.Sprintf("%d:%s", userID, userData.Username)
			r.logger.Info("User connected",
				zap.String("username", userData.Username),
				zap.String("remote_addr", c.Request.RemoteAddr))

			// Отправляем подтверждение авторизации
			response := WSMessage{
				Type: "auth_success",
				Data: json.RawMessage(fmt.Sprintf(`{"username": "%s"}`, userData.Username)),
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				r.logger.Error("Error marshaling auth success response",
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(messageType, responseBytes); err != nil {
				r.logger.Error("Error sending auth success response",
					zap.Error(err))
				continue
			}

			// Загружаем и отправляем новые сообщения
			recentMessages, err := r.chatRepo.GetRecentMessages(c.Request.Context(), 50)
			if err != nil {
				r.logger.Error("Error loading recent messages",
					zap.Error(err))
			} else {
				// Отправляем только новые сообщения
				for i := len(recentMessages) - 1; i >= 0; i-- {
					msg := recentMessages[i]
					msgTime := msg.CreatedAt.Unix()

					// Пропускаем старые сообщения
					if msgTime <= lastMessageTimestamp {
						continue
					}

					response := WSMessage{
						Type:      "message",
						Content:   msg.Content,
						Author:    msg.AuthorUsername,
						ID:        fmt.Sprintf("%d:%d", msg.AuthorID, msgTime),
						Timestamp: msgTime,
					}
					responseBytes, err := json.Marshal(response)
					if err != nil {
						r.logger.Error("Error marshaling message",
							zap.Error(err))
						continue
					}
					if err := conn.WriteMessage(websocket.TextMessage, responseBytes); err != nil {
						r.logger.Error("Error sending message",
							zap.Error(err))
					}
				}
			}
			continue
		}

		// Проверяем, авторизован ли клиент
		clientInfo, exists := r.clients[conn]
		if !exists {
			errorResponse := WSMessage{
				Type:    "error",
				Content: "You must authenticate first",
			}
			errorBytes, err := json.Marshal(errorResponse)
			if err != nil {
				r.logger.Error("Error marshaling error response",
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(messageType, errorBytes); err != nil {
				r.logger.Error("Error sending error response",
					zap.Error(err))
			}
			continue
		}

		// Обработка сообщений чата
		if wsMsg.Type == "message" {
			// Извлекаем username из clientInfo
			parts := strings.Split(clientInfo, ":")
			if len(parts) != 2 {
				r.logger.Error("Invalid client info format",
					zap.String("client_info", clientInfo))
				continue
			}
			userID, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				r.logger.Error("Error parsing user ID",
					zap.Error(err))
				continue
			}
			username := parts[1]

			// Генерируем уникальный ID сообщения
			messageID := fmt.Sprintf("%d:%d", userID, time.Now().UnixNano())

			// Сохраняем сообщение в базу данных
			message := &entity.ChatMessage{
				Content:        wsMsg.Content,
				AuthorID:       userID,
				AuthorUsername: username,
				CreatedAt:      time.Now(),
				ExpiresAt:      time.Now().Add(15 * time.Minute),
			}
			if err := r.chatRepo.SaveMessage(c.Request.Context(), message); err != nil {
				r.logger.Error("Error saving message",
					zap.Error(err))
				continue
			}

			// Отправляем сообщение всем клиентам, кроме отправителя
			response := WSMessage{
				Type:      "message",
				Content:   wsMsg.Content,
				Author:    username,
				ID:        messageID,
				Timestamp: time.Now().Unix(),
			}
			responseBytes, err := json.Marshal(response)
			if err != nil {
				r.logger.Error("Error marshaling message response",
					zap.Error(err))
				continue
			}

			for client := range r.clients {
				// Не отправляем сообщение отправителю
				if client == conn {
					continue
				}
				if err := client.WriteMessage(messageType, responseBytes); err != nil {
					r.logger.Error("Error broadcasting message",
						zap.Error(err))
					client.Close()
					delete(r.clients, client)
				}
			}

			// Отправляем подтверждение отправителю
			confirmation := WSMessage{
				Type:      "message_sent",
				ID:        messageID,
				Timestamp: time.Now().Unix(),
			}
			confirmationBytes, err := json.Marshal(confirmation)
			if err != nil {
				r.logger.Error("Error marshaling confirmation",
					zap.Error(err))
				continue
			}
			if err := conn.WriteMessage(messageType, confirmationBytes); err != nil {
				r.logger.Error("Error sending confirmation",
					zap.Error(err))
			}
		}
	}

	// Очищаем информацию о клиенте при отключении
	delete(r.clients, conn)
	r.logger.Info("Client disconnected",
		zap.String("remote_addr", c.Request.RemoteAddr))
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

// Engine returns the underlying gin.Engine instance
func (r *Router) Engine() *gin.Engine {
	return r.engine
}
