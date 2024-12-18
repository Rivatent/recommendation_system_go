package repository

import (
	"errors"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r *Repo) GetUsersRepo() ([]User, error) {
	rows, err := r.db.Query("SELECT id, username, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) CreateUserRepo(user User) (User, error) {
	if user.Username == "" || user.Email == "" {
		return User{}, errors.New("username and email are required")
	}

	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	var userID int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRow(query, user.Username, user.Email).Scan(&userID, &createdAt, &updatedAt)
	if err != nil {
		return User{}, err
	}
	user.ID = userID
	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	return user, nil
}

func (r *Repo) UpdateUserRepo(user User) (User, error) {
	query := `UPDATE users SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`
	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *Repo) GetUserByIDRepo(id int) (User, error) {
	var user User
	row := r.db.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
