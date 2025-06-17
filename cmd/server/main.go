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
	cfg "insider-project/internal/config"
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