package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"user-service/internal/model"
)

func TestGetUsersRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	expectedUsers := []model.User{
		{
			ID:        "1",
			Username:  "testuser1",
			Email:     "test1@example.com",
			Password:  "hashedpassword1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "2",
			Username:  "testuser2",
			Email:     "test2@example.com",
			Password:  "hashedpassword2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
		AddRow(expectedUsers[0].ID, expectedUsers[0].Username, expectedUsers[0].Email, expectedUsers[0].Password, expectedUsers[0].CreatedAt, expectedUsers[0].UpdatedAt).
		AddRow(expectedUsers[1].ID, expectedUsers[1].Username, expectedUsers[1].Email, expectedUsers[1].Password, expectedUsers[1].CreatedAt, expectedUsers[1].UpdatedAt)

	mock.ExpectQuery("SELECT id, username, email, password, created_at, updated_at FROM users").WillReturnRows(rows)

	users, err := repo.GetUsersRepo()

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserRepo(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	expectedUserID := "123e4567-e89b-12d3-a456-426614174000"

	mock.ExpectQuery(`INSERT INTO users \(username, email, password\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(
			testUser.Username,
			testUser.Email,
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedUserID))

	userID, err := repo.CreateUserRepo(testUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserRepo_QueryError(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		Username: "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testUser.Password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	mock.ExpectQuery(`INSERT INTO users \(username, email, password\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(testUser.Username, testUser.Email, string(hashedPassword)).
		WillReturnError(errors.New("query error"))

	_, err = repo.CreateUserRepo(testUser)

	assert.Error(t, err)

}

func TestUpdateUserRepo_Success_NoPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		Username: "updateduser",
		Email:    "updateduser@example.com",
	}

	mock.ExpectExec(`UPDATE users SET username = \$1, email = \$2, updated_at = CURRENT_TIMESTAMP WHERE id = \$3`).
		WithArgs(testUser.Username, testUser.Email, testUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedUser, err := repo.UpdateUserRepo(testUser)

	assert.NoError(t, err)
	assert.Equal(t, testUser, updatedUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserRepo_Success_WithPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Password: "newpassword",
	}

	mock.ExpectExec(`UPDATE users SET username = \$1, email = \$2, password = \$3, updated_at = CURRENT_TIMESTAMP WHERE id = \$4`).
		WithArgs(testUser.Username, testUser.Email, sqlmock.AnyArg(), testUser.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedUser, err := repo.UpdateUserRepo(testUser)

	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, updatedUser.ID)
	assert.Equal(t, testUser.Username, updatedUser.Username)
	assert.Equal(t, testUser.Email, updatedUser.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserRepo_HashPasswordError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Password: string(make([]byte, bcrypt.MaxCost+1)),
	}

	_, err = repo.UpdateUserRepo(testUser)

	assert.Error(t, err)
}

func TestUpdateUserRepo_SQLExecError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		ID:       "123e4567-e89b-12d3-a456-426614174000",
		Username: "updateduser",
		Email:    "updateduser@example.com",
		Password: "newpassword",
	}

	mock.ExpectExec(`UPDATE users SET username = \$1, email = \$2, password = \$3, updated_at = CURRENT_TIMESTAMP WHERE id = \$4`).
		WithArgs(testUser.Username, testUser.Email, sqlmock.AnyArg(), testUser.ID).
		WillReturnError(errors.New("database error"))

	_, err = repo.UpdateUserRepo(testUser)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDRepo_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	testUser := model.User{
		ID:        "123e4567-e89b-12d3-a456-426614174000",
		Username:  "testuser",
		Email:     "testuser@example.com",
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs(testUser.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
			AddRow(testUser.ID, testUser.Username, testUser.Email, testUser.Password, testUser.CreatedAt, testUser.UpdatedAt))

	user, err := repo.GetUserByIDRepo(testUser.ID)

	assert.NoError(t, err)
	assert.Equal(t, testUser, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDRepo_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs("non-existent-id").
		WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUserByIDRepo("non-existent-id")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDRepo_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnError(errors.New("database error"))

	_, err = repo.GetUserByIDRepo("123e4567-e89b-12d3-a456-426614174000")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDRepo_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &Repo{db: db}

	mock.ExpectQuery(`SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = \$1`).
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at"}).
			AddRow("123e4567-e89b-12d3-a456-426614174000", "testuser", "testuser@example.com", "hashedpassword", "2024-01-01T00:00:00Z"))

	_, err = repo.GetUserByIDRepo("123e4567-e89b-12d3-a456-426614174000")

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
