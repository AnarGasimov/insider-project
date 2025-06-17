package main

import (
	"log"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" 
	_ "github.com/golang-migrate/migrate/v4/source/file"   
	cfg "insider-project/internal/config"    
)

func main() {
	 
	dbURL := cfg.GetMigrateDB_URL()

	migrationsPath := "file:///app/insider-project/migrations"
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	if err == nil {
		log.Println("Database migrations applied successfully!")
	} else if err == migrate.ErrNoChange {
		log.Println("No new database migrations to apply.")
	}
}
