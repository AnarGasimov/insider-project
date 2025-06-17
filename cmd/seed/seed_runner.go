package main

import (
	"log"
	"insider-project/internal/db"
	pkgdb "insider-project/pkg/db"  
	cfg "insider-project/internal/config"  
)

func main() {
	 
	dsn := cfg.GetDB_DSN()
	 
	if err := db.InitDB(dsn); err != nil {
		log.Fatalf("Failed to initialize database connection for seeding: %v", err)
	}
 
	defer func() {
		if err := db.DB.Close(); err != nil {
			log.Printf("Error closing database connection after seeding: %v", err)
		}
	}()

	if err := pkgdb.SeedMessages(db.DB); err != nil {
		log.Fatalf("Failed to seed messages: %v", err)
	}
	log.Println("Database seeding completed.")
}
