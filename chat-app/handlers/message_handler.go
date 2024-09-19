package handlers

import (
	"chat-app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostMessage(service *services.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json struct {
			Username string `json:"username"`
			Message  string `json:"message"`
		}

		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		err := service.AddMessage(json.Username, json.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Message received"})
	}
}

func GetMessages(service *services.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		messages, err := service.GetMessages()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}
		c.JSON(http.StatusOK, messages)
	}
}

func GetMessageCount(service *services.MessageService) gin.HandlerFunc {
	return func(c *gin.Context) {
		count, err := service.GetMessageCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count messages"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"count": count})
	}
}
