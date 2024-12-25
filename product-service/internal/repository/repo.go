package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

// Repo представляет собой структуру для управления подключением к базе данных.
// Она содержит указатель на объект sql.DB, который предоставляет методы
// для взаимодействия с базой данных.
type Repo struct {
	db *sql.DB
}

// New создает новое подключение к базе данных.
// Он считывает строку подключения из переменной окружения DATABASE_URL,
// открывает соединение с PostgreSQL, и проверяет его доступность
// с помощью метода Ping().
// Если возникают ошибки при открытии соединения или его проверке,
// функция вызывает panic и завершает выполнение программы.
// Возвращает указатель на экземпляр Repo.
func New() *Repo {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
		return nil
	}
	if err = db.Ping(); err != nil {
		panic(err)
		return nil
	}

	return &Repo{db: db}
}

// Close закрывает соединение с базой данных.
// Возвращает ошибку, если она возникла при закрытии соединения.
func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}

	return nil
}
