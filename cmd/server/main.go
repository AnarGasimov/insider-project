package main

import (

	"log"
	"os"
	"time" 
	"github.com/gin-gonic/gin"
	"insider-project/internal/api"
	"insider-project/internal/cache"
	"insider-project/internal/db"
	"insider-project/internal/scheduler"
	"github.com/golang-migrate/migrate/v4"
	cfg "insider-project/internal/config"
	pkgdb "insider-project/pkg/db"  
	_ "github.com/golang-migrate/migrate/v4/database/postgres"  
	_ "github.com/golang-migrate/migrate/v4/source/file" 

)

func main() {

	dsn := cfg.GetDB_DSN()

	const maxRetries = 10
	const retryDelay = 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		err := db.InitDB(dsn)
		if err == nil {
			log.Println("Successfully connected to the database.")
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v. Retrying in %s...", i+1, maxRetries, err, retryDelay)
		time.Sleep(retryDelay)
		if i == maxRetries-1 {
			log.Fatal("Exceeded maximum database connection retries.")
		}
	}

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


	log.Println("Seeding initial messages during server startup...")
	if err := pkgdb.SeedMessages(db.DB); err != nil { 
		log.Fatalf("Failed to seed messages during server startup: %v", err)
	}
	log.Println("Database seeding completed during server startup!")
	 


	redisAddr := os.Getenv("REDIS_ADDR")
	if err := cache.InitRedis(redisAddr); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Successfully connected to Redis.")

	
	go scheduler.Start()
	log.Println("Scheduler started.")

	
	r := gin.Default()
	api.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}