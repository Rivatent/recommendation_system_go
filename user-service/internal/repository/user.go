package repository

import "errors"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *Repo) GetUsersRepo() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) CreateUserRepo(user User) (User, error) {
	if user.Name == "" || user.Email == "" {
		return User{}, errors.New("name and email are required")
	}

	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
	var userID int
	err := r.db.QueryRow(query, user.Name, user.Email).Scan(&userID)
	if err != nil {
		return User{}, err
	}
	user.ID = userID

	return user, nil
}

func (r *Repo) UpdateUserRepo(user User) (User, error) {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := r.db.Exec(query, user.Name, user.Email, user.ID)
	if err != nil {
		return User{}, err
	}
	return user, err
}

func (r *Repo) GetUserByIDRepo(id int) (User, error) {
	var user User
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
