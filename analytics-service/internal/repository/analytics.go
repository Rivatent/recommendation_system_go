package repository

import (
	"analytics-service/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	queryProductNewMsg = ``
)

func (r *Repo) GetAnalyticsRepo() ([]model.Analytics, error) {
	var analytics []model.Analytics

	rows, err := r.db.Query("SELECT * FROM analytics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var a model.Analytics
		if err := rows.Scan(&a.ID, &a.TotalUsers, &a.TotalSales, &a.SalesProgressionRate, &a.UsersProgressionRate, &a.ProductAverageRating, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		analytics = append(analytics, a)
	}

	return analytics, nil
}

func (r *Repo) ProductUpdateMsgRepo(updatedProduct map[string]interface{}) error {

	return nil
}

func (r *Repo) UserNewMsgRepo() error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var totalUsers int
	err = tx.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalUsers)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get total users: %w", err)
	}

	var totalSales int
	var avgRating float64
	err = tx.QueryRow(`
		SELECT SUM(sales_count), COALESCE(AVG(rating), 0.00)
		FROM products
	`).Scan(&totalSales, &avgRating)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get sales and average rating: %w", err)
	}

	var lastAnalyticsID string
	var previousTotalUsers, previousTotalSales int

	err = tx.QueryRow(`
		SELECT id, total_users, total_sales
		FROM analytics
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&lastAnalyticsID, &previousTotalUsers, &previousTotalSales)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return fmt.Errorf("failed to fetch previous analytics: %w", err)
	}

	usersProgressionRate := 0.00
	salesProgressionRate := 0.00

	if previousTotalUsers > 0 {
		usersProgressionRate = float64(totalUsers-previousTotalUsers) / float64(previousTotalUsers) * 100
	}
	if previousTotalSales > 0 {
		salesProgressionRate = float64(totalSales-previousTotalSales) / float64(previousTotalSales) * 100
	}

	_, err = tx.Exec(`
		INSERT INTO analytics (
			total_users, total_sales, sales_progression_rate, users_progression_rate, product_average_rating, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $6)
	`, totalUsers, totalSales, salesProgressionRate, usersProgressionRate, avgRating, time.Now())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert analytics record: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *Repo) ProductNewMsgRepo(newProduct map[string]interface{}) error {

	return nil
}
