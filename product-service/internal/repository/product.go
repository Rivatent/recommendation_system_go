package repository

import (
	"database/sql"
	"errors"
	"product-service/internal/model"
)

func (r *Repo) GetProductsRepo() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, description, price, created_at, updated_at FROM Products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt); err != nil {
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

	query := `INSERT INTO Products (name, description, price, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id`
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price).Scan(&productID)
	if err != nil {
		return productID, err
	}

	return productID, nil
}

func (r *Repo) UpdateProductRepo(product model.Product) (model.Product, error) {
	if product.Name == "" {
		return model.Product{}, errors.New("Product name is required")
	}

	query := `UPDATE Products SET name = $1, description = $2, price = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := r.db.Exec(query, product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *Repo) GetProductByIDRepo(id string) (model.Product, error) {
	var product model.Product
	row := r.db.QueryRow("SELECT id, name, description, price, created_at, updated_at FROM Products WHERE id = $1", id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, errors.New("Product not found")
		}
		return model.Product{}, err
	}

	return product, nil
}

func (r *Repo) DeleteProductByIDRepo(id string) error {
	query := "DELETE FROM Products WHERE id = $1"
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
