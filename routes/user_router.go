package routes

import (
	"car-system-go/handler"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	UserGroup := r.Group("/api/user")
	{
		UserGroup.POST("/findAll", handler.UserFindAllHandler)
		UserGroup.POST("/createRecord", handler.UserInfractionCreateHandler)
	}
}
