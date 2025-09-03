package routes

import (
	"car-system-go/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) {
	AdminGroup := r.Group("/api/admin")
	{
		AdminGroup.POST("/register", handler.AdminRegisterHandler)
		AdminGroup.POST("/login", handler.AdminLoginHandler)
		AdminGroup.POST("/findPassword", handler.AdminFindPasswordHandler)
	}
}
