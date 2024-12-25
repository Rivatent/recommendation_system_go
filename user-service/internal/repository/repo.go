package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

// Repo представляет собой репозиторий для работы с базой данных.
type Repo struct {
	db *sql.DB
}

// New создает новый экземпляр Repo, устанавливает соединение с базой данных
// и возвращает указатель на репозиторий.
// Соединение к базе данных строится с использованием URL,
// полученного из переменной окружения DATABASE_URL.
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
// Возвращает ошибку, если не удалось закрыть соединение.
func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
