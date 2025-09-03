package routes

import (
	"car-system-go/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCodeRouter(r *gin.Engine) {
	CodeGroup := r.Group("/api/code")
	{
		CodeGroup.POST("/getVerificationCode", handler.VerificationCodeGetHandler)
	}
}
