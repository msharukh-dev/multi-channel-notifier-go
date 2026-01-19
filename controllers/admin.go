package controllers

import (
	"net/http"
	"time"
	"webhook-api/config"
	"webhook-api/dto"
	"webhook-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AdminLogin handles admin authentication
// TODO: Implement JWT token generation
func AdminLogin(c *gin.Context) {
	var req dto.AdminLoginRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.AdminLoginResponse{
			Status:  "error",
			Message: "Invalid request",
		})
		return
	}

	// Find admin user
	var admin models.AdminUser
	if err := config.DB.Where("email = ?", req.Email).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.AdminLoginResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.AdminLoginResponse{
			Status:  "error",
			Message: "Invalid email or password",
		})
		return
	}

	// TODO: Generate JWT token
	token := "jwt_token_here"
	expiresAt := time.Now().Add(24 * time.Hour)

	c.JSON(http.StatusOK, dto.AdminLoginResponse{
		Status:  "success",
		Message: "Login successful",
		Data: &dto.AdminToken{
			Token:     token,
			Email:     admin.Email,
			ExpiresAt: expiresAt.Format(time.RFC3339),
		},
	})
}

// GetAllClients returns all registered clients
// TODO: Add pagination
func GetAllClients(c *gin.Context) {
	var clients []models.Client
	if err := config.DB.Find(&clients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch clients",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Clients retrieved",
		"data":    clients,
	})
}

// DeactivateClient disables a client account
func DeactivateClient(c *gin.Context) {
	clientID := c.Param("client_id")

	if err := config.DB.Model(&models.Client{}).Where("id = ?", clientID).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to deactivate client",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Client deactivated",
	})
}

// GetNotificationStats returns notification delivery statistics
func GetNotificationStats(c *gin.Context) {
	var totalNotifications int64
	var sentNotifications int64
	var failedNotifications int64

	config.DB.Model(&models.Notification{}).Count(&totalNotifications)
	config.DB.Model(&models.Notification{}).Where("status = ?", "sent").Count(&sentNotifications)
	config.DB.Model(&models.Notification{}).Where("status = ?", "failed").Count(&failedNotifications)

	successRate := 0.0
	if totalNotifications > 0 {
		successRate = (float64(sentNotifications) / float64(totalNotifications)) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"total":        totalNotifications,
			"sent":         sentNotifications,
			"failed":       failedNotifications,
			"success_rate": successRate,
			"pending":      totalNotifications - sentNotifications - failedNotifications,
		},
	})
}
