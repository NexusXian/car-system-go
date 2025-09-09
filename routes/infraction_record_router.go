package routes

import (
	"car-system-go/handler"
	"car-system-go/jwt"

	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	RecordGroup := r.Group("/api/record")
	RecordGroup.Use(jwt.AuthRequired())
	{
		RecordGroup.POST("/findAll", handler.InfractionRecordFindAllHandler)
		RecordGroup.POST("/findByIDCardNumber", handler.InfractionRecordFindByIDCardNumberHandler)
	}
}
