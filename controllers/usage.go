package controllers

import (
	"net/http"
	"time"
	"webhook-api/config"
	"webhook-api/dto"
	"webhook-api/models"

	"github.com/gin-gonic/gin"
)

// GetUsage retrieves usage statistics for the authenticated client
func GetUsage(c *gin.Context) {
	clientID := c.GetUint("client_id")

	// Fetch client details
	var client models.Client
	if err := config.DB.First(&client, clientID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.UsageResponse{
			Status:  "error",
			Message: "Client not found",
		})
		return
	}

	// Get today's usage
	today := time.Now().Truncate(24 * time.Hour)
	var todayUsage int64
	if err := config.DB.Model(&models.Notification{}).
		Where("client_id = ? AND created_at >= ? AND status = ?", clientID, today, "sent").
		Count(&todayUsage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.UsageResponse{
			Status:  "error",
			Message: "Failed to fetch today's usage",
		})
		return
	}

	// Get this month's usage
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var monthlyUsage int64
	if err := config.DB.Model(&models.Notification{}).
		Where("client_id = ? AND created_at >= ? AND status = ?", clientID, monthStart, "sent").
		Count(&monthlyUsage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.UsageResponse{
			Status:  "error",
			Message: "Failed to fetch monthly usage",
		})
		return
	}

	// Calculate remaining
	remainingToday := int(int64(client.DailyLimit) - todayUsage)
	if remainingToday < 0 {
		remainingToday = 0
	}

	remainingMonth := int(int64(client.MonthlyLimit) - monthlyUsage)
	if remainingMonth < 0 {
		remainingMonth = 0
	}

	// Calculate percentages
	percentageToday := 0.0
	if client.DailyLimit > 0 {
		percentageToday = (float64(todayUsage) / float64(client.DailyLimit)) * 100
	}

	percentageMonth := 0.0
	if client.MonthlyLimit > 0 {
		percentageMonth = (float64(monthlyUsage) / float64(client.MonthlyLimit)) * 100
	}

	// Get last reset time (start of today)
	lastReset := today.Format("2006-01-02T15:04:05Z07:00")

	c.JSON(http.StatusOK, dto.UsageResponse{
		Status:  "success",
		Message: "Usage retrieved successfully",
		Data: &dto.UsageDataInfo{
			TodayUsage:         int(todayUsage),
			MonthlyUsage:       int(monthlyUsage),
			DailyLimit:         client.DailyLimit,
			MonthlyLimit:       client.MonthlyLimit,
			RemainingToday:     remainingToday,
			RemainingThisMonth: remainingMonth,
			PercentageToday:    percentageToday,
			PercentageMonth:    percentageMonth,
			LastReset:          lastReset,
		},
	})
}
