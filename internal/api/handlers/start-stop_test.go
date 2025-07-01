package handlers_test

import (
	"database/sql"
	"errors"
	"insider-project/internal/api/handlers"
	"insider-project/internal/cache"
	"insider-project/internal/db"
	"insider-project/internal/scheduler"
	"insider-project/internal/sender"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	startCalled    bool
	processCalled  bool
)

func init() {
    cache.StoreMessage = func(id int, messageId string, sentAt time.Time) error {
        return nil // test zamanı heç bir əməliyyat aparmır
    }
}

type testStore struct{}

func (s testStore) GetUnsentMessages(limit int)  ([]db.Message, error) {
	return []db.Message{{ID: 1, Phone: "123456", Content: "test"}}, nil
}

func (s testStore) MarkMessageSent(id int) error {
	return nil
}
func (s testStore) GetSentMessages(limit, offset int) ([]db.Message, error) {
	return nil, nil
}
func fakeStartSuccess() error {
	startCalled = true
	return nil
}

func fakeStartError() error {
	return errors.New("start failed")
}

func fakeProcess() {
	processCalled = true
}

func (s testStore) GetDB() *sql.DB {
	return nil
}

type Sender struct{}

func (s Sender) SendMessage(to, content string) (string, error) {
	return "mock-id-123", nil
}

func TestStartHandlerSuccess(t *testing.T) {
	
	scheduler.SetStore(testStore{})
	scheduler.StartFunc = fakeStartSuccess
	scheduler.ProcessMessagesFunc = fakeProcess
	sender.Service = Sender{}

	startCalled = false
	processCalled = false

	router := gin.Default()
	router.POST("/start", handlers.StartHandler)

	req, _ := http.NewRequest("POST", "/start", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.True(t, startCalled)
	assert.True(t, processCalled)
}

func TestStartHandlerError(t *testing.T) {
	scheduler.StartFunc = fakeStartError
	router := gin.Default()
	router.POST("/start", handlers.StartHandler)

	req, _ := http.NewRequest("POST", "/start", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestStopHandler(t *testing.T) {
	scheduler.SetStore(testStore{})

	router := gin.Default()
	router.POST("/stop", handlers.StopHandler)

	req, _ := http.NewRequest("POST", "/stop", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
