package routes

import (
	"car-system-go/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAiRouter(r *gin.Engine) {
	AdminGroup := r.Group("/api/ai")
	{
		AdminGroup.POST("/analyze", handler.AiAnalyzeHandler)
		AdminGroup.POST("/stream", handler.AiAnalyzeStreamHandler)
	}
}
