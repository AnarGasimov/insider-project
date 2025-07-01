package api

import (
    "github.com/gin-gonic/gin"
    "insider-project/internal/api/handlers"
    "insider-project/internal/db"
)

func SetupRoutes(r *gin.Engine, store db.MessageStore) {
    h := handlers.Handler{Store: store}

    v1 := r.Group("/api/v1")
    {
        v1.POST("/start", handlers.StartHandler)
        v1.POST("/stop", handlers.StopHandler)
        v1.GET("/sent", h.GetSentMessages)
    }
}