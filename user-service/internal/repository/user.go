package repository

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user-service/internal/model"
)

func (r *Repo) GetUsersRepo() ([]model.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, password, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) CreateUserRepo(user model.User) (string, error) {
	var userID string

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return userID, errors.New("failed to hash password")
	}

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = r.db.QueryRow(query, user.Username, user.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (r *Repo) UpdateUserRepo(user model.User) (model.User, error) {
	if user.Username == "" || user.Email == "" {
		return model.User{}, errors.New("username and email are required")
	}

	var query string
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return model.User{}, errors.New("failed to hash password")
		}

		query = `UPDATE users SET username = $1, email = $2, password = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
		_, err = r.db.Exec(query, user.Username, user.Email, string(hashedPassword), user.ID)
		if err != nil {
			return model.User{}, err
		}
	} else {
		query = `UPDATE users SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
		_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
		if err != nil {
			return model.User{}, err
		}
	}

	return user, nil
}

func (r *Repo) GetUserByIDRepo(id string) (model.User, error) {
	var user model.User
	row := r.db.QueryRow("SELECT id, username, email, password, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.User{}, errors.New("user not found")
		}
		return model.User{}, err
	}

	return user, nil
}
