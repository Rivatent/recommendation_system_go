package repository

import (
	"database/sql"
	"errors"
	"math"
	"product-service/internal/model"
)

func (r *Repo) GetProductsRepo() ([]model.Product, error) {
	rows, err := r.db.Query(`
        SELECT id, name, description, price, rating, sales_count, created_at, updated_at
        FROM products
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Rating, &product.SalesCount, &product.CreatedAt, &product.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repo) CreateProductRepo(product model.Product) (string, error) {
	var productID string

	if product.Name == "" {
		return productID, errors.New("Product name is required")
	}
	if product.Price <= 0 {
		return productID, errors.New("Product price must be greater than zero")
	}

	query := `
        INSERT INTO products (name, description, price, rating, sales_count, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) 
        RETURNING id
    `
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Rating, product.SalesCount).Scan(&productID)
	if err != nil {
		return productID, err
	}

	return productID, nil
}

//func (r *Repo) UpdateProductRepo(product model.Product) (model.Product, error) {
//	if product.Name == "" {
//		return model.Product{}, errors.New("Product name is required")
//	}
//
//	query := `
//        UPDATE products
//        SET name = $1, description = $2, price = $3, rating = $4, sales_count = $5, updated_at = CURRENT_TIMESTAMP
//        WHERE id = $6
//    `
//	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.Rating, product.SalesCount, product.ID)
//	if err != nil {
//		return model.Product{}, err
//	}
//
//	return product, nil
//}

func (r *Repo) UpdateProductRepo(product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("Product name is required")
	}

	// Логика пересчета рейтинга, если он не указан
	var newRating float64

	// Получаем текущий SalesCount из базы
	var salesCount int
	err := r.db.QueryRow("SELECT sales_count FROM products WHERE id = $1", product.ID).Scan(&salesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, errors.New("Product not found")
		}
		return model.Product{}, err
	}

	// Рассчитываем новый рейтинг
	baseRating := 3.0
	salesDivisor := 10.0
	newRating = math.Min(5.0, baseRating+float64(salesCount)/salesDivisor)

	// Обновляем продукт в базе данных
	query := `
        UPDATE products 
        SET name = $1, description = $2, price = $3, rating = $4, updated_at = CURRENT_TIMESTAMP 
        WHERE id = $5
    `
	_, err = r.db.Exec(query, product.Name, product.Description, product.Price, newRating, product.ID)
	if err != nil {
		return model.Product{}, err
	}

	// Возвращаем обновленный продукт
	product.Rating = newRating
	return product, nil
}

func (r *Repo) GetProductByIDRepo(id string) (model.Product, error) {
	var product model.Product
	query := `
        SELECT id, name, description, price, rating, sales_count, created_at, updated_at 
        FROM products 
        WHERE id = $1
    `
	err := r.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.Rating, &product.SalesCount, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, errors.New("Product not found")
		}
		return model.Product{}, err
	}

	return product, nil
}

func (r *Repo) DeleteProductByIDRepo(id string) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
