package handlers_test

import (
	"database/sql"
	"insider-project/internal/api/handlers"
	"insider-project/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type mockStore struct{}

func (m *mockStore) GetSentMessages(limit, offset int) ([]db.Message, error) {
	return []db.Message{
		{ID: 1, Content: "Test content", Phone: "+123456789"},
	}, nil
}
func (m *mockStore) GetUnsentMessages(limit int) ([]db.Message, error) { 
	return nil, nil
}

func (m *mockStore) MarkMessageSent(id int) error { 
	return nil 
}

func (m *mockStore) GetDB() *sql.DB {
	return nil
}

func TestGetSentMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := &handlers.Handler{Store: &mockStore{}}
	router := gin.Default()
	router.GET("/sent-messages", handler.GetSentMessages)

	req, _ := http.NewRequest("GET", "/sent-messages?limit=10&offset=0", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.Code)
	}

	var messages []db.Message
	if err := json.Unmarshal(resp.Body.Bytes(), &messages); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(messages) != 1 || messages[0].Content != "Test content" {
		t.Errorf("unexpected messages: %+v", messages)
	}
}
