package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type Repo struct {
	db *sql.DB
}

func New() *Repo {
	dsn := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("========= Connection to database error: %v =========", err)
		return nil
	}

	if err = db.Ping(); err != nil {
		log.Printf("========= Ping to database error: %v =========", err)
		return nil
	}

	log.Println("========= Connected to database via ping =========")

	return &Repo{db: db}
}

func (r *Repo) Close() error {
	err := r.db.Close()
	if err != nil {
		return err
	}
	return nil
}
