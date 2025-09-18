package routes

import (
	"car-system-go/handler"
	"car-system-go/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	RecordGroup := r.Group("/api/record")
	{
		RecordGroup.POST("/findByIDCard", handler.InfractionRecordFindByIDCardNumberHandler)
		RecordGroup.POST("/findAll", handler.InfractionRecordFindAllHandler).Use(jwt.AuthRequired())
		RecordGroup.POST("/findByIDCardNumber", handler.InfractionRecordFindByIDCardNumberHandler).Use(jwt.AuthRequired())
	}
}
