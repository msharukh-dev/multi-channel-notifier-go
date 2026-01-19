package dto

type SendRequest struct {
	Type    string `json:"type" binding:"required"`
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject"`
	Message string `json:"message" binding:"required"`
}

type SendResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    *SendDataInfo `json:"data,omitempty"`
}

type SendDataInfo struct {
	NotificationID uint   `json:"notification_id"`
	Type           string `json:"type"`
	To             string `json:"to"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
}
