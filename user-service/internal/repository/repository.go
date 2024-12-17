package repository

import "database/sql"

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Функция для вставки пользователя
func CreateUser(db *sql.DB, name, email string) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2)`
	_, err := db.Exec(query, name, email)
	return err
}

// Функция для получения всех пользователей
func GetUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
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

// Функция для получения пользователя по ID
func GetUserByID(db *sql.DB, id uint64) (User, error) {
	var user User
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Функция для обновления пользователя
func UpdateUser(db *sql.DB, id uint64, name, email string) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := db.Exec(query, name, email, id)
	return err
}

//func createTable(db *sql.DB) {
//	query := `
//    CREATE TABLE IF NOT EXISTS users (
//        id SERIAL PRIMARY KEY,
//        name TEXT NOT NULL,
//        email TEXT UNIQUE NOT NULL
//    )`
//	if _, err := db.Exec(query); err != nil {
//		log.Fatalf("Ошибка создания таблицы: %v", err)
//	}
//}
