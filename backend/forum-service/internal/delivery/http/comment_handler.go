package httpDelivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sout1235/forum2/backend/forum-service/internal/entity"
	"github.com/sout1235/forum2/backend/forum-service/internal/repository"
	"github.com/sout1235/forum2/backend/forum-service/internal/usecase"
)

// CommentHandler handles HTTP requests for comments
type CommentHandler struct {
	commentUseCase usecase.CommentUseCase
	userRepo       repository.UserRepository
}

// CommentRequest represents a request to create a comment
// @Description Request body for creating a comment
type CommentRequest struct {
	Content string `json:"content" binding:"required" example:"This is a great post!"`
}

func NewCommentHandler(commentUseCase usecase.CommentUseCase, userRepo repository.UserRepository) *CommentHandler {
	return &CommentHandler{
		commentUseCase: commentUseCase,
		userRepo:       userRepo,
	}
}

// @Summary Get all comments for a topic
// @Description Get all comments associated with a specific topic
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Topic ID"
// @Success 200 {array} Comment
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id}/comments [get]
func (h *CommentHandler) GetAllCommentsByTopic(c *gin.Context) {
	topicID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID"})
		return
	}
	comments, err := h.commentUseCase.GetCommentsByTopicID(c.Request.Context(), topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// @Summary Get a specific comment
// @Description Get detailed information about a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} Comment
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id} [get]
func (h *CommentHandler) GetComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}
	comment, err := h.commentUseCase.GetCommentByID(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment)
}

// @Summary Create a new comment
// @Description Create a new comment on a topic
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Param comment body CommentRequest true "Comment information"
// @Success 201 {object} Comment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment entity.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.commentUseCase.CreateComment(c.Request.Context(), &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary Delete a comment
// @Description Delete a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Topic ID"
// @Param commentId path int true "Comment ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /topics/{id}/comments/{commentId} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	err = h.commentUseCase.DeleteComment(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Like a comment
// @Description Add a like to a specific comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 "OK"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /comments/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	err = h.commentUseCase.LikeComment(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
