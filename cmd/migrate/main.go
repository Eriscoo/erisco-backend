package main

import (
	"log"

	"github.com/eriscoo/blog-backend/internal/infrastructure/config"
	"github.com/eriscoo/blog-backend/internal/infrastructure/persistence"
)

func main() {
	cfg := config.Load()

	db, err := persistence.OpenDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	persistence.RunMigrations(db)
}
