package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "insider-project/internal/scheduler"
)

func StartHandler(c *gin.Context) {
    if err := scheduler.StartFunc(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    scheduler.ProcessMessagesFunc()
    c.JSON(http.StatusOK, gin.H{"message": "Scheduler started"})
}

func StopHandler(c *gin.Context) {
    scheduler.Stop()
    c.JSON(http.StatusOK, gin.H{"message": "Scheduler stopped"})
}