package config

import (
	"fmt"
	"os"
)

func GetDB_DSN() string {
	dbHost := "db"
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)
}

func GetMigrateDB_URL() string {
    dbHost := "db"
    dbUser := os.Getenv("POSTGRES_USER")
    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    return fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
        dbUser, dbPassword, dbHost, dbName)
}