package dto

type RegisterRequest struct {
	ClientName   string `json:"client_name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Website      string `json:"website"`
	WebhookURL   string `json:"webhook_url"`
	DailyLimit   int    `json:"daily_limit"`
	MonthlyLimit int    `json:"monthly_limit"`
}

type RegisterResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    *ApiKeyData `json:"data,omitempty"`
}

type ApiKeyData struct {
	ClientID     uint   `json:"client_id"`
	ClientName   string `json:"client_name"`
	Email        string `json:"email"`
	APIKey       string `json:"api_key"`
	DailyLimit   int    `json:"daily_limit"`
	MonthlyLimit int    `json:"monthly_limit"`
}
