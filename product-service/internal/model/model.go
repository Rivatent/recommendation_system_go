package model

import "time"

// Product представляет собой модель продукта в системе.
// Эта структура используется для хранения информации о продукте и передачи данных
// через API. Каждый продукт имеет уникальный идентификатор, имя, описание, цену,
// рейтинг, количество продаж и временные метки для отслеживания создания и обновления.
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       float64   `json:"price" validate:"required"`
	Rating      float64   `json:"rating"`
	SalesCount  int       `json:"sales_count" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
