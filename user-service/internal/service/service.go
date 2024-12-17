package service

import (
	"database/sql"
	"log"
	"os"
)

func ConnectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Printf("========= Connection to database error: %v =========", err)
		return nil, err
	}
	//defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("========= Ping to database error: %v =========", err)
		return nil, err
	}
	log.Println("========= Connected to database via ping =========")
	//createTable(db)
	return db, nil
}
