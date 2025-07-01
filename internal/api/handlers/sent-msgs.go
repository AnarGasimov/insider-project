package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "insider-project/internal/db"
)
type Handler struct {
    Store db.MessageStore
}

func (h *Handler) GetSentMessages(c *gin.Context) {
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

    messages, err := h.Store.GetSentMessages(limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, messages)
}