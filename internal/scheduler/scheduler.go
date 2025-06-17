package scheduler

import (
    "log"
    "sync"
    "time"
    "insider-project/internal/cache"
    "insider-project/internal/db"
    "insider-project/internal/sender"
)

var (
    running  bool
    mu       sync.Mutex
    stopChan = make(chan struct{})
)

func Start() error {
    mu.Lock()
    if running {
        mu.Unlock()
        return nil
    }
    running = true
    mu.Unlock()

    go func() {
        ticker := time.NewTicker(2 * time.Minute)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                messages, err := db.GetUnsentMessages(2)
                if err != nil {
                    log.Println("Error fetching messages:", err)
                    continue
                }
                for _, msg := range messages {
                    messageID, err := sender.SendMessage(msg.Phone, msg.Content)
                    if err != nil {
                        log.Println("Error sending message:", err)
                        continue
                    }
                    if err := db.MarkMessageSent(msg.ID); err != nil {
                        log.Println("Error marking message sent:", err)
                        continue
                    }
                    if err := cache.StoreMessage(msg.ID, messageID, time.Now()); err != nil {
                        log.Println("Error caching message:", err)
                    }
                }
            case <-stopChan:
                return
            }
        }
    }()
    return nil
}

func Stop() {
    mu.Lock()
    if !running {
        mu.Unlock()
        return
    }
    running = false
    mu.Unlock()
    stopChan <- struct{}{}
}