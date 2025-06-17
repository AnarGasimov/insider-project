package api

import (
    "github.com/gin-gonic/gin"
    "insider-project/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
    v1 := r.Group("/api/v1")
    {
        v1.POST("/start", handlers.StartHandler)
        v1.POST("/stop", handlers.StopHandler)
        v1.GET("/sent", handlers.GetSentMessages)
    }
}