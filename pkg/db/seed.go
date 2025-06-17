package db

import (
    "database/sql"
    "fmt"
    "log"
    "time"
)

func SeedMessages(db *sql.DB) error {
    log.Println("Starting database seeding...")
    for i := 0; i < 20; i++ { // 20 messages 
        content := fmt.Sprintf("Initial Message %d - %s", i+1, time.Now().Format("15:04:05"))
        phone := fmt.Sprintf("+90555%07d", 1000000+i)

        _, err := db.Exec(`
            INSERT INTO messages (content, phone, sent)
            VALUES ($1, $2, FALSE)
        `, content, phone)
        if err != nil {
            return fmt.Errorf("failed to insert message %d: %w", i, err)
        }
    }
    log.Println("Messages seeded successfully!")
    return nil
}