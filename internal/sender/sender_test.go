package sender

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSendMessageSuccess(t *testing.T){
 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message":"ok","messageId":"123"}`))
    }))
    defer ts.Close()

    os.Setenv("WEBHOOK_URL", ts.URL)

    client := &http.Client{Timeout: 5 * time.Second}
    msg, err := SendMessageWithClient("12345", "Test message", client)
    if err != nil {
        t.Errorf("Expected success, got error: %v", err)
    }
    if !strings.Contains(msg, "123") {
        t.Errorf("Expected message to contain 'ok', got: %s", msg)
    }
}

func TestSendMessageFailureStatus(t *testing.T){
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusBadRequest)
    }))
    defer ts.Close()

    os.Setenv("WEBHOOK_URL", ts.URL)

    client := &http.Client{Timeout: 5 * time.Second}
    _, err := SendMessageWithClient("12345", "Test message", client)
    if err == nil || !strings.Contains(err.Error(), "webhook request failed") {
        t.Errorf("Expected webhook failure error, got: %v", err)
    }
}

func TestSendMessageContent(t *testing.T){
	 _, err := SendMessageWithClient("12345", strings.Repeat("A", 161), http.DefaultClient)
    if err == nil || err.Error() != "content exceeds 160 characters" {
        t.Errorf("Expected error for content too long, got: %v", err)
    }
}