package main

import (
	"context"
	"log"
	"recommendation-service/internal/app"
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
