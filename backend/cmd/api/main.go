package main

import (
	"context"
	"log"

	"github.com/imamsyafii-dev/ledgerflow/backend/internal/config"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/database"
	"github.com/imamsyafii-dev/ledgerflow/backend/internal/http/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("connected to PostgreSQL")

	engine := router.New(db, cfg)

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
