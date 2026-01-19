package controllers

import (
	"net/http"
	"strconv"
	"webhook-api/config"
	"webhook-api/dto"
	"webhook-api/models"

	"github.com/gin-gonic/gin"
)

// GetStatus retrieves the status of a notification by ID
func GetStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.StatusResponse{
			Status:  "error",
			Message: "Invalid notification ID",
		})
		return
	}

	// Get API key from context (set by middleware)
	clientID := c.GetUint("client_id")

	// Fetch notification and verify it belongs to the client
	var notification models.Notification
	if err := config.DB.Where("id = ? AND client_id = ?", uint(id), clientID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.StatusResponse{
			Status:  "error",
			Message: "Notification not found",
		})
		return
	}

	var sentAtStr *string
	if notification.SentAt != nil {
		sentAt := notification.SentAt.Format("2006-01-02T15:04:05Z07:00")
		sentAtStr = &sentAt
	}

	c.JSON(http.StatusOK, dto.StatusResponse{
		Status:  "success",
		Message: "Notification status retrieved",
		Data: &dto.NotificationData{
			ID:           notification.ID,
			Type:         notification.NotificationType,
			To:           notification.To,
			Subject:      notification.Subject,
			Status:       notification.Status,
			ErrorMessage: notification.ErrorMessage,
			SentAt:       sentAtStr,
			RetryCount:   notification.RetryCount,
			CreatedAt:    notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}
