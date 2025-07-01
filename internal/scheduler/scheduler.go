package scheduler

import (
	"insider-project/internal/cache"
	"insider-project/internal/db"
	"insider-project/internal/sender"
	"log"
	"sync"
	"time"
) 



var (
	running  bool
	mu       sync.Mutex
	stopChan = make(chan struct{})
	store db.MessageStore
	StartFunc func()error
	ProcessMessagesFunc func()
)

func init(){
	StartFunc = Start
	ProcessMessagesFunc = ProcessMessages
}

func SetStore(s db.MessageStore) {
	store = s
}

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
				log.Println("Scheduler: 2-minute interval passed. Processing messages...")
				ProcessMessagesFunc()
			case <-stopChan:
				log.Println("Scheduler received stop signal. Stopping...")
				return
			}
		}
	}()
	return nil
}

func ProcessMessages() {

	messages, err := store.GetUnsentMessages(2)
	if err != nil {
		log.Println("Error fetching messages:", err)
		return
	}
	for _, msg := range messages {
		log.Println("msg: ", msg)
		messageID, err := sender.Service.SendMessage(msg.Phone, msg.Content)
		if err != nil {
			log.Println("Error sending message:", err)
			continue
		}
		if err := store.MarkMessageSent(msg.ID); err != nil {
			log.Println("Error marking message sent:", err)
			continue
		}
		if err := cache.StoreMessage(msg.ID, messageID, time.Now()); err != nil {
			log.Println("Error caching message:", err)
		}
	}
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