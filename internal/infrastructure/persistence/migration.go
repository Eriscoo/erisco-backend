package persistence

import (
	"database/sql"
	"log"

	"github.com/eriscoo/blog-backend/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func RunMigrations(db *sql.DB) {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		log.Fatalf("migration source: %v", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("migration driver: %v", err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate setup: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate up: %v", err)
	}
	log.Println("migrations applied successfully")
}
