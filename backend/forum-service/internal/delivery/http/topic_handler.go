package httpDelivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/service"
)

// TopicHandler handles HTTP requests for topics
type TopicHandler struct {
	topicUseCase service.TopicService
	userRepo     repository.UserRepository
}

// Author represents a topic or comment author
// @Description Author information
type Author struct {
	ID       int64  `json:"id" example:"1"`
	Username string `json:"username" example:"johndoe"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
}

// Category represents a topic category
// @Description Category information
type Category struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"General"`
}

// Tag represents a topic tag
// @Description Tag information
type Tag struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"golang"`
}

// TopicResponse represents a topic with all its details
// @Description Topic information with comments and metadata
type TopicResponse struct {
	ID           int64     `json:"id" example:"1"`
	Title        string    `json:"title" example:"How to use Go"`
	Content      string    `json:"content" example:"This is a tutorial about Go programming language"`
	AuthorID     int64     `json:"author_id" example:"1"`
	Author       Author    `json:"author"`
	CategoryID   int64     `json:"category_id" example:"1"`
	Category     Category  `json:"category"`
	Views        int       `json:"views" example:"100"`
	Comments     []Comment `json:"comments"`
	Tags         []Tag     `json:"tags"`
	CreatedAt    string    `json:"created_at" example:"2024-03-15T10:00:00Z"`
	UpdatedAt    string    `json:"updated_at" example:"2024-03-15T10:00:00Z"`
	CommentCount int       `json:"comment_count" example:"5"`
}

// Comment represents a comment on a topic
// @Description Comment information
type Comment struct {
	ID        int64  `json:"id" example:"1"`
	Content   string `json:"content" example:"Great post!"`
	AuthorID  int64  `json:"author_id" example:"1"`
	Author    Author `json:"author"`
	TopicID   int64  `json:"topic_id" example:"1"`
	CreatedAt string `json:"created_at" example:"2024-03-15T10:00:00Z"`
	UpdatedAt string `json:"updated_at" example:"2024-03-15T10:00:00Z"`
}

func NewTopicHandler(topicUseCase service.TopicService, userRepo repository.UserRepository) *TopicHandler {
	return &TopicHandler{
		topicUseCase: topicUseCase,
		userRepo:     userRepo,
	}
}

// @Summary Create a new topic
// @Description Create a new topic with the provided information
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param topic body entity.Topic true "Topic information"
// @Success 201 {object} TopicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics [post]
func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var topic entity.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	topic.AuthorID = userID.(int64)

	err := h.topicUseCase.CreateTopic(c.Request.Context(), &topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get author information
	author, err := h.userRepo.GetUserByID(c.Request.Context(), topic.AuthorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	topic.Author = author

	c.JSON(http.StatusCreated, topic)
}

// @Summary Get a topic by ID
// @Description Get detailed information about a specific topic
// @Tags topics
// @Accept json
// @Produce json
// @Param id path int true "Topic ID"
// @Success 200 {object} TopicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id} [get]
func (h *TopicHandler) GetTopic(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	topic, err := h.topicUseCase.GetTopicByID(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get author information if not already set
	if topic.Author == nil && topic.AuthorID > 0 {
		author, err := h.userRepo.GetUserByID(c.Request.Context(), topic.AuthorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		topic.Author = author
	}

	c.JSON(http.StatusOK, topic)
}

// @Summary Get all topics
// @Description Get a list of all topics with pagination
// @Tags topics
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} TopicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics [get]
func (h *TopicHandler) GetAllTopics(c *gin.Context) {
	log.Printf("Getting all topics")
	topics, err := h.topicUseCase.GetAllTopics(c.Request.Context())
	if err != nil {
		log.Printf("Error getting topics: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get topics: %v", err)})
		return
	}
	log.Printf("Successfully retrieved %d topics", len(topics))

	// Get author information for topics that don't have it
	for i, topic := range topics {
		if topic.AuthorID > 0 && (topic.Author == nil || topic.Author.Username == "") {
			log.Printf("Getting author information for topic %d (author_id: %d)", topic.ID, topic.AuthorID)
			author, err := h.userRepo.GetUserByID(c.Request.Context(), topic.AuthorID)
			if err != nil {
				log.Printf("Error getting author for topic %d: %v", topic.ID, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get author for topic %d: %v", topic.ID, err)})
				return
			}
			topics[i].Author = author
			log.Printf("Successfully retrieved author for topic %d: %+v", topic.ID, author)
		}
	}

	c.JSON(http.StatusOK, topics)
}

// @Summary Update a topic
// @Description Update an existing topic's information
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Param topic body entity.Topic true "Updated topic information"
// @Success 200 {object} TopicResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id} [put]
func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	var topic entity.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topic.ID = topicID

	err = h.topicUseCase.UpdateTopic(c.Request.Context(), &topic)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, topic)
}

// @Summary Delete a topic
// @Description Delete a topic by ID
// @Tags topics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id} [delete]
func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	err = h.topicUseCase.DeleteTopic(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TopicHandler) UpdateCommentCount(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}

	err = h.topicUseCase.UpdateCommentCount(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
