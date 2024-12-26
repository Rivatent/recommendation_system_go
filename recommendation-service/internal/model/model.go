package model

import (
	"time"
)

// Recommendation - структура, представляющая рекомендацию продукта для пользователя.
// Она содержит информацию о том, какой продукт рекомендован конкретному пользователю,
// а также включает метрики для оценки и даты создания и обновления данной рекомендации.
type Recommendation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Score     float64   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
