package routes

import (
	"webhook-api/controllers"
	"webhook-api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// Public endpoint - register new client
		api.POST("/register", controllers.RegisterAPIKey)

		// Protected endpoints - require API key
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Send notification
			protected.POST("/send", controllers.SendNotification)

			// Get notification status
			protected.GET("/status/:id", controllers.GetStatus)

			// Get usage statistics
			protected.GET("/usage", controllers.GetUsage)
		}
	}
}
