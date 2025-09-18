package main

import (
	"car-system-go/service"
	"car-system-go/setup"

	"github.com/gin-gonic/gin"
)

func init() {
	setup.InitAvatar()
	setup.InitViper()
	setup.InitMySQL()

}

func main() {
	r := setup.RouterInit()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	go service.AlcoholSmokeService()
	if err := r.Run(":8100"); err != nil {
		panic(err)
	}
}
