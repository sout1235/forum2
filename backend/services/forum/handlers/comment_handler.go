package handlers

import (
	"net/http"
	"strconv"

	"backend/services/forum/models"
	"backend/services/forum/repository"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Repo *repository.CommentRepository
}

func NewCommentHandler(repo *repository.CommentRepository) *CommentHandler {
	return &CommentHandler{Repo: repo}
}

func (h *CommentHandler) GetAllCommentsByTopic(c *gin.Context) {
	topicIDStr := c.Param("topic_id")
	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный topic_id"})
		return
	}
	comments, err := h.Repo.GetAllCommentsByTopic(topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения комментариев"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	comment, err := h.Repo.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Комментарий не найден"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	if err := h.Repo.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания комментария"})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	comment.ID = id
	if err := h.Repo.UpdateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления комментария"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	if err := h.Repo.DeleteComment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления комментария"})
		return
	}
	c.Status(http.StatusNoContent)
}
