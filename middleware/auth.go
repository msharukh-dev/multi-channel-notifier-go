package middleware

import (
	"net/http"
	"webhook-api/config"
	"webhook-api/models"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates API key from request header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Missing X-API-Key header",
			})
			c.Abort()
			return
		}

		// Validate API key exists and is active in database
		var key models.APIKey
		if err := config.DB.Where("key = ? AND is_active = ?", apiKey, true).First(&key).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid API key",
			})
			c.Abort()
			return
		}

		// Verify client is active
		var client models.Client
		if err := config.DB.First(&client, key.ClientID).Error; err != nil || !client.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "Client account is not active",
			})
			c.Abort()
			return
		}

		// Store client info in context for use in controllers
		c.Set("client_id", client.ID)
		c.Set("api_key", apiKey)
		c.Set("client", client)

		c.Next()
	}
}
