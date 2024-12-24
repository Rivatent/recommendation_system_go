package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"product-service/internal/model"
	"testing"
	"time"
)

func TestGetProductsRepo_Success(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "rating", "sales_count", "created_at", "updated_at"}).
		AddRow("1", "Product 1", "Description 1", 100.0, 4.5, 10, time.Now(), time.Now()).
		AddRow("2", "Product 2", "Description 2", 200.0, 4.8, 20, time.Now(), time.Now())

	mock.ExpectQuery(`
		SELECT id, name, description, price, rating, sales_count, created_at, updated_at
		FROM products
	`).WillReturnRows(rows)

	products, err := repo.GetProductsRepo()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Product 1", products[0].Name)
	assert.Equal(t, 100.0, products[0].Price)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsRepo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectQuery(`
		SELECT id, name, description, price, rating, sales_count, created_at, updated_at
		FROM products
	`).WillReturnError(sqlmock.ErrCancelled)

	products, err := repo.GetProductsRepo()

	assert.Error(t, err)
	assert.Nil(t, products)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsRepo_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "rating", "sales_count", "created_at", "updated_at"}).
		AddRow("1", "Product 1", "Description 1", "invalid_price", 4.5, 10, time.Now(), time.Now())

	mock.ExpectQuery(`
		SELECT id, name, description, price, rating, sales_count, created_at, updated_at
		FROM products
	`).WillReturnRows(rows)

	products, err := repo.GetProductsRepo()

	assert.Error(t, err)
	assert.Nil(t, products)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductsRepo_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "rating", "sales_count", "created_at", "updated_at"}).
		AddRow("1", "Product 1", "Description 1", 100.0, 4.5, 10, time.Now(), time.Now()).
		RowError(0, sqlmock.ErrCancelled)

	mock.ExpectQuery(`
		SELECT id, name, description, price, rating, sales_count, created_at, updated_at
		FROM products
	`).WillReturnRows(rows)

	products, err := repo.GetProductsRepo()

	assert.Error(t, err)
	assert.Nil(t, products)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProductRepo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Rating:      4.5,
		SalesCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery(`
	INSERT INTO products \(name, description, price, rating, sales_count, created_at, updated_at\)
	VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\)
	RETURNING id
`).
		WithArgs(product.Name, product.Description, product.Price, product.Rating, product.SalesCount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("123"))

	productID, err := repo.CreateProductRepo(product)

	assert.NoError(t, err)
	assert.Equal(t, "123", productID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProductRepo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Rating:      4.5,
		SalesCount:  10,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectQuery(`
	INSERT INTO products \(name, description, price, rating, sales_count, created_at, updated_at\)
	VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\)
	RETURNING id
`).
		WithArgs(product.Name, product.Description, product.Price, product.Rating, product.SalesCount).
		WillReturnError(fmt.Errorf("canceling query due to user request"))

	productID, err := repo.CreateProductRepo(product)

	assert.Error(t, err)
	assert.Equal(t, "", productID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProductRepo_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		Rating:      4.5,
		SalesCount:  10,
	}

	mock.ExpectQuery(`
		INSERT INTO products \(name, description, price, rating, sales_count, created_at, updated_at\)
		VALUES \(\$1, \$2, \$3, \$4, \$5, NOW\(\), NOW\(\)\)
		RETURNING id
	`).
		WithArgs(product.Name, product.Description, product.Price, product.Rating, product.SalesCount).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))

	productID, err := repo.CreateProductRepo(product)

	assert.Error(t, err)
	assert.Equal(t, "", productID)
	assert.Contains(t, err.Error(), "Scan error")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductRepo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		ID:          "123",
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       150.0,
		SalesCount:  20,
	}

	mock.ExpectQuery(`SELECT sales_count FROM products WHERE id = \$1`).
		WithArgs(product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"sales_count"}).AddRow(10))

	mock.ExpectExec(`UPDATE products SET name = \$1, description = \$2, price = \$3, sales_count = \$4, rating = \$5, updated_at = CURRENT_TIMESTAMP WHERE id = \$6`).
		WithArgs(product.Name, product.Description, product.Price, product.SalesCount, 4.0, product.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updatedProduct, err := repo.UpdateProductRepo(product)

	assert.NoError(t, err)
	assert.Equal(t, 4.0, updatedProduct.Rating)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductRepo_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		ID: "123",
	}

	mock.ExpectQuery("SELECT sales_count FROM products WHERE id = \\$1").
		WithArgs(product.ID).
		WillReturnError(sql.ErrNoRows)

	updatedProduct, err := repo.UpdateProductRepo(product)

	assert.Error(t, err)
	assert.Equal(t, "Product not found", err.Error())
	assert.Equal(t, model.Product{}, updatedProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductRepo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		ID: "123",
	}

	mock.ExpectQuery("SELECT sales_count FROM products WHERE id = \\$1").
		WithArgs(product.ID).
		WillReturnError(errors.New("query error"))

	updatedProduct, err := repo.UpdateProductRepo(product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "query error")
	assert.Equal(t, model.Product{}, updatedProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductRepo_UpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	product := model.Product{
		ID:          "123",
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       150.0,
		SalesCount:  20,
	}

	mock.ExpectQuery(`SELECT sales_count FROM products WHERE id = \$1`).
		WithArgs(product.ID).
		WillReturnRows(sqlmock.NewRows([]string{"sales_count"}).AddRow(10))

	mock.ExpectExec(`UPDATE products SET name = \$1, description = \$2, price = \$3, sales_count = \$4, rating = \$5, updated_at = CURRENT_TIMESTAMP WHERE id = \$6`).
		WithArgs(product.Name, product.Description, product.Price, product.SalesCount, 4.0, product.ID).
		WillReturnError(errors.New("update error"))

	updatedProduct, err := repo.UpdateProductRepo(product)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.Equal(t, model.Product{}, updatedProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductByIDRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("successful retrieval", func(t *testing.T) {
		productID := "123"
		expectedProduct := model.Product{
			ID:          productID,
			Name:        "Test Product",
			Description: "This is a test product.",
			Price:       19.99,
			Rating:      4.5,
			SalesCount:  100,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "rating", "sales_count", "created_at", "updated_at"}).
			AddRow(expectedProduct.ID, expectedProduct.Name, expectedProduct.Description, expectedProduct.Price, expectedProduct.Rating, expectedProduct.SalesCount, expectedProduct.CreatedAt, expectedProduct.UpdatedAt)

		mock.ExpectQuery("^SELECT (.+) FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnRows(rows)

		product, err := repo.GetProductByIDRepo(productID)

		assert.NoError(t, err)
		assert.Equal(t, expectedProduct, product)
	})

	t.Run("product not found", func(t *testing.T) {
		productID := "456"

		mock.ExpectQuery("^SELECT (.+) FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetProductByIDRepo(productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "Product not found")
		assert.Equal(t, model.Product{}, product)
	})

	t.Run("database error", func(t *testing.T) {
		productID := "789"

		mock.ExpectQuery("^SELECT (.+) FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnError(errors.New("database error"))

		product, err := repo.GetProductByIDRepo(productID)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		assert.Equal(t, model.Product{}, product)
	})

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteProductByIDRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	t.Run("successful deletion", func(t *testing.T) {
		productID := "123"

		mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteProductByIDRepo(productID)
		require.NoError(t, err)

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("product not found", func(t *testing.T) {
		productID := "456"

		mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteProductByIDRepo(productID)
		assert.EqualError(t, err, "Product not found")

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})

	t.Run("database error", func(t *testing.T) {
		productID := "789"

		mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
			WithArgs(productID).
			WillReturnError(errors.New("database error"))

		err := repo.DeleteProductByIDRepo(productID)
		assert.EqualError(t, err, "database error")

		err = mock.ExpectationsWereMet()
		require.NoError(t, err)
	})
}
