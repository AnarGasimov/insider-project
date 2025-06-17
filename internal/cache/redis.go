package cache

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"github.com/go-redis/redis/v8"
)

var Client *redis.Client

func InitRedis(addr string) error {
    Client = redis.NewClient(&redis.Options{
        Addr: addr,
    })
    _, err := Client.Ping(context.Background()).Result()
    return err
}

func StoreMessage(id int, messageId string, sentAt time.Time) error {
    data, _ := json.Marshal(map[string]interface{}{
        "messageId": messageId,
        "sent_at":   sentAt,
    })
    return Client.Set(context.Background(), "message:"+strconv.Itoa(id), data, 24*time.Hour).Err()
}