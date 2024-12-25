package repository

import (
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user-service/internal/model"
)

// GetUsersRepo возвращает список пользователей из репозитория.
//
// Выполняет SQL-запрос к базе данных, чтобы извлечь
// все записи из таблицы пользователей. Каждый пользователь
// представляется структурой model.User, и данные заполняются
// полями: ID, Username, Email, Password, CreatedAt и UpdatedAt.
//
// Если запрос к базе данных не удался, функция вернет ошибку.
// Если во время сканирования строк произойдет ошибка, она также будет возвращена.
//
// Возвращаемые значения:
//   - []model.User: срез, содержащий пользователей,
//     или nil в случае сбоя.
//   - error: ошибка запроса, если таковая произошла, или nil,
//     если операция прошла успешно.
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

// CreateUserRepo Создает нового пользователя в базе данных.
// Принимает структуру model.User, содержащую данные пользователя,
// включая имя пользователя, электронную почту и пароль.
// Возвращает уникальный идентификатор созданного пользователя и ошибку, если произошла ошибка.
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

// UpdateUserRepo Обновляет данные пользователя в базе данных.
// Принимает структуру model.User, содержащую обновленные данные пользователя.
// Если новый пароль указан, он будет предварительно хэширован.
// Возвращает обновленную структуру model.User и ошибку, если произошла ошибка.
func (r *Repo) UpdateUserRepo(user model.User) (model.User, error) {

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

// GetUserByIDRepo Получает пользователя из базы данных по идентификатору.
// Принимает строковый идентификатор пользователя.
// Возвращает структуру model.User с данными пользователя и ошибку,
// если пользователь не найден или произошла другая ошибка.
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
