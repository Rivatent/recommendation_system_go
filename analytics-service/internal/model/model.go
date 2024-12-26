package model

import "time"

// Analytics - структура, представляющая аналитическую информацию о пользователях и продажах.
// Эта структура используется для хранения и передачи данных о метриках аналитики.
type Analytics struct {
	ID                   string    `json:"id"`
	TotalUsers           int       `json:"total_users"`
	TotalSales           int       `json:"total_sales"`
	SalesProgressionRate float64   `json:"sales_progression_rate"`
	UsersProgressionRate float64   `json:"users_progression_rate"`
	ProductAverageRating float64   `json:"product_average_rating"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
