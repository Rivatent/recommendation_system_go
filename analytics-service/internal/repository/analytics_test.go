package repository

import (
	"analytics-service/internal/model"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetAnalyticsRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %s", err)
	}
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectQuery("^SELECT (.+) FROM analytics$").WillReturnRows(sqlmock.NewRows([]string{
		"id", "total_users", "total_sales", "sales_progression_rate",
		"users_progression_rate", "product_average_rating", "created_at", "updated_at"}).
		AddRow("1", 10, 100, 0.10, 0.20, 4.5, time.Now(), time.Now()))

	analytics, err := repo.GetAnalyticsRepo()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := []model.Analytics{
		{
			ID:                   "1",
			TotalUsers:           10,
			TotalSales:           100,
			SalesProgressionRate: 0.10,
			UsersProgressionRate: 0.20,
			ProductAverageRating: 4.5,
			CreatedAt:            time.Time{},
			UpdatedAt:            time.Time{},
		},
	}

	if len(analytics) != len(expected) {
		t.Errorf("expected %v, got %v", expected, analytics)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateAnalyticsMsgRepo_Success(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок-объекта базы данных: %v", err)
	}
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

	mock.ExpectQuery(`SELECT SUM\(sales_count\), COALESCE\(AVG\(rating\), 0.00\) FROM products`).
		WillReturnRows(sqlmock.NewRows([]string{"sum", "avg"}).AddRow(2000, 4.5))

	mock.ExpectQuery(`SELECT id, total_users, total_sales FROM analytics ORDER BY created_at DESC LIMIT 1`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "total_users", "total_sales"}).
			AddRow("analytics_id", 50, 1500))

	mock.ExpectExec(`INSERT INTO analytics`).
		WithArgs(100, 2000, 0.3333333333333333, 1.0, 4.5, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.UpdateAnalyticsMsgRepo()

	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}

func TestUpdateAnalyticsMsgRepo_Fail_BeginTransaction(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок-объекта базы данных: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to begin transaction"))

	repo := &Repo{db: db}

	err = repo.UpdateAnalyticsMsgRepo()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to begin transaction")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}

func TestUpdateAnalyticsMsgRepo_Fail_QueryRow_Error(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок-объекта базы данных: %v", err)
	}
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
		WillReturnError(fmt.Errorf("failed to get total users"))

	err = repo.UpdateAnalyticsMsgRepo()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get total users")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Невыполненные ожидания: %v", err)
	}
}
