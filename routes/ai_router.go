package routes

import (
	"car-system-go/handler"
	"car-system-go/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterAiRouter(r *gin.Engine) {
	AdminGroup := r.Group("/api/ai")
	AdminGroup.Use(jwt.AuthRequired())
	{
		AdminGroup.POST("/analyze", handler.AiAnalyzeHandler)
	}
}
