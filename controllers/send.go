package controllers

import (
	"net/http"
	"time"
	"webhook-api/config"
	"webhook-api/dto"
	"webhook-api/models"
	"webhook-api/utils"

	"github.com/gin-gonic/gin"
)

// SendNotification sends a notification and stores it in the database
func SendNotification(c *gin.Context) {
	var req dto.SendRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.SendResponse{
			Status:  "error",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate notification type
	if req.Type != "email" && req.Type != "sms" && req.Type != "webhook" {
		c.JSON(http.StatusBadRequest, dto.SendResponse{
			Status:  "error",
			Message: "Invalid notification type. Supported: email, sms, webhook",
		})
		return
	}

	// Get client info from context (set by middleware)
	clientID := c.GetUint("client_id")

	// Get client details
	var client models.Client
	if err := config.DB.First(&client, clientID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.SendResponse{
			Status:  "error",
			Message: "Invalid API key or client not found",
		})
		return
	}

	// Check if client is active
	if !client.IsActive {
		c.JSON(http.StatusForbidden, dto.SendResponse{
			Status:  "error",
			Message: "Client account is inactive",
		})
		return
	}

	// Check daily limit
	today := time.Now().Truncate(24 * time.Hour)
	var todayCount int64
	config.DB.Model(&models.Notification{}).
		Where("client_id = ? AND created_at >= ? AND status = ?", clientID, today, "sent").
		Count(&todayCount)

	if int(todayCount) >= client.DailyLimit {
		c.JSON(http.StatusTooManyRequests, dto.SendResponse{
			Status:  "error",
			Message: "Daily limit reached. Please try again tomorrow.",
		})
		return
	}

	// Create notification record with pending status
	notification := models.Notification{
		ClientID:         clientID,
		NotificationType: req.Type,
		To:               req.To,
		Subject:          req.Subject,
		Message:          req.Message,
		Status:           "pending",
		RetryCount:       0,
	}

	// Save to database
	if err := config.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.SendResponse{
			Status:  "error",
			Message: "Failed to save notification: " + err.Error(),
		})
		return
	}

	// Send notification asynchronously
	go func() {
		err := utils.Send(req.Type, req.To, req.Message, client.WebhookURL)

		// Update notification status
		updateData := map[string]interface{}{
			"status": "sent",
		}

		if err != nil {
			updateData["status"] = "failed"
			updateData["error_message"] = err.Error()
			updateData["retry_count"] = 0
		} else {
			now := time.Now()
			updateData["sent_at"] = now
		}

		config.DB.Model(&notification).Updates(updateData)
	}()

	c.JSON(http.StatusAccepted, dto.SendResponse{
		Status:  "success",
		Message: "Notification queued for delivery",
		Data: &dto.SendDataInfo{
			NotificationID: notification.ID,
			Type:           notification.NotificationType,
			To:             notification.To,
			Status:         notification.Status,
			CreatedAt:      notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	})
}
