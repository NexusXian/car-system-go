package main

import (
	"car-system-go/service"
	"car-system-go/setup"
)

func init() {
	setup.InitAvatar()
	setup.InitViper()
	setup.InitMySQL()

}

func main() {
	r := setup.RouterInit()

	go service.AlcoholSmokeService()
	if err := r.Run(":8200"); err != nil {
		panic(err)
	}
}
