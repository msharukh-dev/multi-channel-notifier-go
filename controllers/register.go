package controllers

import (
	"net/http"
	"regexp"
	"webhook-api/config"
	"webhook-api/dto"
	"webhook-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterAPIKey creates a new client and generates an API key
func RegisterAPIKey(c *gin.Context) {
	var req dto.RegisterRequest

	// Validate request
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Status:  "error",
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Status:  "error",
			Message: "Invalid email format",
		})
		return
	}

	// Set defaults
	if req.DailyLimit == 0 {
		req.DailyLimit = 1000
	}
	if req.MonthlyLimit == 0 {
		req.MonthlyLimit = 30000
	}

	// Validate limits
	if req.DailyLimit <= 0 || req.MonthlyLimit <= 0 {
		c.JSON(http.StatusBadRequest, dto.RegisterResponse{
			Status:  "error",
			Message: "Daily and monthly limits must be positive numbers",
		})
		return
	}

	// Check if client already exists
	var existingClient models.Client
	if err := config.DB.Where("email = ?", req.Email).First(&existingClient).Error; err == nil {
		c.JSON(http.StatusConflict, dto.RegisterResponse{
			Status:  "error",
			Message: "Client with this email already exists",
		})
		return
	}

	// Create new client
	client := models.Client{
		Name:         req.ClientName,
		Email:        req.Email,
		Website:      req.Website,
		WebhookURL:   req.WebhookURL,
		DailyLimit:   req.DailyLimit,
		MonthlyLimit: req.MonthlyLimit,
		IsActive:     true,
	}

	if err := config.DB.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Status:  "error",
			Message: "Failed to create client: " + err.Error(),
		})
		return
	}

	// Generate API key
	apiKey := uuid.New().String()
	key := models.APIKey{
		Key:      apiKey,
		Name:     req.ClientName + " API Key",
		ClientID: client.ID,
		IsActive: true,
	}

	if err := config.DB.Create(&key).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.RegisterResponse{
			Status:  "error",
			Message: "Failed to generate API key: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		Status:  "success",
		Message: "Client registered successfully",
		Data: &dto.ApiKeyData{
			ClientID:     client.ID,
			ClientName:   client.Name,
			Email:        client.Email,
			APIKey:       apiKey,
			DailyLimit:   client.DailyLimit,
			MonthlyLimit: client.MonthlyLimit,
		},
	})
}

// Helper function to validate email
func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
