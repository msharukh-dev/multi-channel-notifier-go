package dto

type StatusResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    *NotificationData `json:"data,omitempty"`
}

type NotificationData struct {
	ID           uint    `json:"id"`
	Type         string  `json:"type"`
	To           string  `json:"to"`
	Subject      string  `json:"subject"`
	Status       string  `json:"status"`
	ErrorMessage string  `json:"error_message,omitempty"`
	SentAt       *string `json:"sent_at,omitempty"`
	RetryCount   int     `json:"retry_count"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}
