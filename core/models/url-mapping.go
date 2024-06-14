package models

type URLMapping struct {
	ShortCode   string `json:"shortCode"`
	OriginalURL string `json:"originalUrl"`
	CreatedAt   string `json:"createdAt"`
	UsageCount  int    `json:"usageCount"`
	IsActive    bool   `json:"isActive"`
	UserID      int    `json:"userId"`
}
