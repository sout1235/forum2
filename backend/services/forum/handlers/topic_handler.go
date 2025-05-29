package handlers

import (
	"net/http"
	"strconv"

	"backend/services/forum/models"
	"backend/services/forum/repository"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	Repo *repository.TopicRepository
}

func NewTopicHandler(repo *repository.TopicRepository) *TopicHandler {
	return &TopicHandler{Repo: repo}
}

func (h *TopicHandler) GetAllTopics(c *gin.Context) {
	topics, err := h.Repo.GetAllTopics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения тем"})
		return
	}
	c.JSON(http.StatusOK, topics)
}

func (h *TopicHandler) GetTopicByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	topic, err := h.Repo.GetTopicByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Тема не найдена"})
		return
	}
	c.JSON(http.StatusOK, topic)
}

func (h *TopicHandler) CreateTopic(c *gin.Context) {
	var t models.Topic
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	if err := h.Repo.CreateTopic(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания темы"})
		return
	}
	c.JSON(http.StatusCreated, t)
}

func (h *TopicHandler) UpdateTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	var t models.Topic
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}
	t.ID = id
	if err := h.Repo.UpdateTopic(&t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления темы"})
		return
	}
	c.JSON(http.StatusOK, t)
}

func (h *TopicHandler) DeleteTopic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
		return
	}
	if err := h.Repo.DeleteTopic(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления темы"})
		return
	}
	c.Status(http.StatusNoContent)
}
