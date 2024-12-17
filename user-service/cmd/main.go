package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
	"user-service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err = a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

//handler := &handlers.Handler{}
//handler.SetDB(db)
//
//r := chi.NewRouter()
//r.Get("/", handler.MainPage)
//r.Get("/users", handler.GerUsers)
//r.Get("/user/{id}", handler.GetUserByID)
//r.Post("/user", handler.CreateUser)
//r.Put("/user/{id}", handler.UpdateUserByID)
//
//log.Println("Started server at :8081")
//if err := http.ListenAndServe(":8081", r); err != nil {
//	log.Fatalf("Starting server error: %v", err)
//}
//
//fmt.Println("========= starting server at :8081 =========")
