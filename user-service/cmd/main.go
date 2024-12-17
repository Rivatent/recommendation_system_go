package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"user-service/internal/handlers"
	"user-service/internal/service"
)

func main() {

	db, err := service.ConnectDB()
	defer db.Close()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	handler := &handlers.Handler{}
	handler.SetDB(db)

	r := chi.NewRouter()
	r.Get("/", handler.MainPage)
	r.Get("/users", handler.GerUsers)
	r.Get("/user/{id}", handler.GetUserByID)
	r.Post("/user", handler.CreateUser)
	r.Put("/user/{id}", handler.UpdateUserByID)

	log.Println("Started server at :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("Starting server error: %v", err)
	}

	fmt.Println("========= starting server at :8081 =========")

}
