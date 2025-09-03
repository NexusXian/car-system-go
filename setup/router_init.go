package setup

import (
	"car-system-go/middleware"
	"car-system-go/routes"

	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	routes.RegisterAdminRouter(r)
	routes.RegisterRecordRouter(r)
	routes.RegisterAiRouter(r)
	routes.RegisterUserRouter(r)
	routes.RegisterCodeRouter(r)
	return r
}
