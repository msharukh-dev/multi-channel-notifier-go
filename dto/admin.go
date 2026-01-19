package dto

type AdminLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    *AdminToken `json:"data,omitempty"`
}

type AdminToken struct {
	Token     string `json:"token"`
	Email     string `json:"email"`
	ExpiresAt string `json:"expires_at"`
}
