package sender

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

type WebhookResponse struct {
    Message   string `json:"message"`
    MessageID string `json:"messageId"`
}

func SendMessage(to, content string) (string, error) {
    if len(content) > 160 {
        return "", errors.New("content exceeds 160 characters")
    }

    payload, _ := json.Marshal(map[string]string{
        "to":      to,
        "content": content,
    })

    client := &http.Client{Timeout: 10 * time.Second}
    req, _ := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-ins-auth-key", os.Getenv("WEBHOOK_AUTH_KEY"))
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return "", errors.New("webhook request failed")
    }

    var result WebhookResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return result.MessageID, nil
}