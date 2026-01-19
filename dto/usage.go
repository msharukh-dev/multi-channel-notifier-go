package dto

type UsageResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *UsageDataInfo `json:"data,omitempty"`
}

type UsageDataInfo struct {
	TodayUsage         int     `json:"today_usage"`
	MonthlyUsage       int     `json:"monthly_usage"`
	DailyLimit         int     `json:"daily_limit"`
	MonthlyLimit       int     `json:"monthly_limit"`
	RemainingToday     int     `json:"remaining_today"`
	RemainingThisMonth int     `json:"remaining_this_month"`
	PercentageToday    float64 `json:"percentage_today"`
	PercentageMonth    float64 `json:"percentage_month"`
	LastReset          string  `json:"last_reset"`
}
