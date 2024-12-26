package repository

import (
	"database/sql"
	"errors"
	"product-service/internal/model"
)

// GetProductsRepo извлекает список всех продуктов из базы данных.
// Возвращает срез продуктов и ошибку, если она возникла.
// Если при выполнении запроса возникает ошибка, функция возвращает nil и ошибку.
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

// CreateProductRepo добавляет новый продукт в базу данных.
// Принимает продукт в качестве параметра и возвращает его уникальный идентификатор и ошибку.
// Если продукт успешно добавлен, возвращается его ID; в противном случае возвращается ошибка.
func (r *Repo) CreateProductRepo(product model.Product) (string, error) {
	var productID string

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

// UpdateProductRepo обновляет информацию о существующем продукте.
// Принимает продукт в качестве параметра и возвращает обновленный продукт и ошибку.
// Если продукт найден и обновлён, возвращается обновлённый объект; в противном случае возвращается ошибка
func (r *Repo) UpdateProductRepo(product model.Product) (model.Product, error) {
	var oldSalesCount int

	err := r.db.QueryRow("SELECT sales_count FROM products WHERE id = $1", product.ID).Scan(&oldSalesCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, errors.New("Product not found")
		}
		return model.Product{}, err
	}
	newSalesCount := product.SalesCount
	salesChange := newSalesCount - oldSalesCount
	baseRating := 3.0
	salesDivisor := 10.0
	newRating := baseRating + float64(salesChange)/salesDivisor
	if newRating < 0 {
		newRating = 0
	} else if newRating > 5 {
		newRating = 5
	}

	query := `
	UPDATE products
	SET name = $1, description = $2, price = $3, sales_count = $4, rating = $5, updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	`
	_, err = r.db.Exec(query, product.Name, product.Description, product.Price, newSalesCount, newRating, product.ID)
	if err != nil {
		return model.Product{}, err
	}

	product.Rating = newRating

	return product, nil
}

// GetProductByIDRepo извлекает продукт по его уникальному идентификатору.
// Принимает ID продукта в качестве параметра и возвращает продукт и ошибку.
// Если продукт найден, возвращается его объект; в противном случае возвращается ошибка.
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

// DeleteProductByIDRepo удаляет продукт по его уникальному идентификатору.
// Принимает ID продукта в качестве параметра и возвращает ошибку, если она возникла.
// Если продукт успешно удалён, функция возвращает nil; если продукт не найден, возвращается ошибка.
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
