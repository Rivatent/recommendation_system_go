package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

type User struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("========= Connection to database error: %v =========", err)
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("========= Ping to database error: %v =========", err)
		panic(err)
	}
	log.Println("========= Connected to database via ping =========")
	//createTable(db)

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("====== it works ========"))
	})
	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := getUsers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	})
	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}
		user, err := getUserByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})
	r.Post("/user", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Name == "" || user.Email == "" {
			http.Error(w, "Name and Email are required", http.StatusBadRequest)
			return
		}

		if err := createUser(db, user.Name, user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, `{"status":"user created"}`)
	})

	r.Put("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid ID format", http.StatusBadRequest)
			return
		}

		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if user.Name == "" || user.Email == "" {
			http.Error(w, "Name and Email are required", http.StatusBadRequest)
			return
		}

		if err := updateUser(db, id, user.Name, user.Email); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"user updated"}`)
	})

	log.Println("Started server at :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Starting server error: %v", err)
	}

	fmt.Println("========= starting server at :8081 =========")

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "It works, URL:", r.URL.String())
}

// Функция для создания таблицы
func createTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL
    )`
	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}
}

// Функция для вставки пользователя
func createUser(db *sql.DB, name, email string) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2)`
	_, err := db.Exec(query, name, email)
	return err
}

// Функция для получения всех пользователей
func getUsers(db *sql.DB) ([]User, error) {
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
func getUserByID(db *sql.DB, id uint64) (User, error) {
	var user User
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Функция для обновления пользователя
func updateUser(db *sql.DB, id uint64, name, email string) error {
	query := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err := db.Exec(query, name, email, id)
	return err
}
