package main

import (
	"context"
	"log"

	"store-service/internal/application"
)

func main() {
	ctx := context.Background()

	app, err := application.New(ctx)
	if err != nil {
		log.Fatalf("failed to init application: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}
}
