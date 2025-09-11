package routes

import (
	"car-system-go/handler"
	"car-system-go/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	UserGroup := r.Group("/api/user")
	{
		UserGroup.POST("/findAll", handler.UserFindAllHandler)
		UserGroup.POST("/createRecord", handler.UserInfractionCreateHandler)
		UserGroup.GET("/user-findAll", handler.UserFindAllInfoHandler).Use(jwt.AuthRequired())
		UserGroup.POST("/userFind", handler.UserFindHandler).Use(jwt.AuthRequired())
	}
}
