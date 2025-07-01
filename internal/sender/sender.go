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

var Service SenderService = defaultSender{}

type SenderService interface{
    SendMessage(to, content string) (string,error)
}

type defaultSender struct {}

func (s defaultSender) SendMessage(to, content string) (string, error) {
   client := &http.Client{Timeout: 5 * time.Second}
   return SendMessageWithClient(to,content,client)
}

func SendMessageWithClient(to, content string, client *http.Client)(string,error){
    if len(content) > 160 {
        return "", errors.New("content exceeds 160 characters")
    }

    payload, _ := json.Marshal(map[string]string{
        "to":      to,
        "content": content,
    })
    
    req, _ := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")
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