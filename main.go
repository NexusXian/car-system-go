package main

import (
	"car-system-go/setup"
)

func init() {
	setup.InitAvatar()
	setup.InitViper()
	setup.InitMySQL()
}

func main() {
	r := setup.RouterInit()

	if err := r.Run(":8100"); err != nil {
		panic(err)
	}
}
