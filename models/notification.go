package models

import (
	"time"

	"gorm.io/gorm"
)

// APIKey represents a registered API key for clients
type APIKey struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Key       string         `gorm:"uniqueIndex;not null" json:"key"`
	Name      string         `gorm:"not null" json:"name"`
	ClientID  uint           `gorm:"not null" json:"client_id"`
	Client    Client         `gorm:"foreignKey:ClientID" json:"client,omitempty"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Client represents a customer/client
type Client struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null" json:"name"`
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Website       string         `json:"website"`
	WebhookURL    string         `json:"webhook_url"`
	DailyLimit    int            `gorm:"default:1000" json:"daily_limit"`
	MonthlyLimit  int            `gorm:"default:30000" json:"monthly_limit"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	APIKeys       []APIKey       `gorm:"foreignKey:ClientID" json:"api_keys,omitempty"`
	Notifications []Notification `gorm:"foreignKey:ClientID" json:"notifications,omitempty"`
}

// Notification represents a notification sent through the API
type Notification struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ClientID         uint           `gorm:"not null;index" json:"client_id"`
	Client           Client         `gorm:"foreignKey:ClientID" json:"-"`
	NotificationType string         `gorm:"not null" json:"type"` // email, sms, webhook
	To               string         `gorm:"not null" json:"to"`
	Subject          string         `json:"subject"`
	Message          string         `gorm:"type:text;not null" json:"message"`
	Status           string         `gorm:"not null;default:'pending'" json:"status"` // pending, sent, failed
	ErrorMessage     string         `gorm:"type:text" json:"error_message"`
	SentAt           *time.Time     `json:"sent_at"`
	RetryCount       int            `gorm:"default:0" json:"retry_count"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// UsageLog tracks API usage for quota management
type UsageLog struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ClientID          uint      `gorm:"not null;index:idx_client_date" json:"client_id"`
	Client            Client    `gorm:"foreignKey:ClientID" json:"-"`
	NotificationCount int       `gorm:"default:0" json:"notification_count"`
	Date              time.Time `gorm:"index:idx_client_date;not null" json:"date"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AdminUser represents admin users for the platform
type AdminUser struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
