package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

// Repo представляет собой репозиторий, который управляет подключением к базе данных.
// Он содержит ссылку на объект базы данных и предоставляет методы для работы с ним.
type Repo struct {
	db *sql.DB
}

// New создает новое подключение к базе данных PostgreSQL.
// Он использует переменную окружения DATABASE_URL для получения строки подключения.
// В случае ошибки подключения, приложение завершает работу.
func New() *Repo {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	return &Repo{db: db}
}

// Close закрывает соединение с базой данных.
// Возвращает ошибку, если закрытие подключения не удалось.
func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
