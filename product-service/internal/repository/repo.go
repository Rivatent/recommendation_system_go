package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

type Repo struct {
	db *sql.DB
}

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

func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
