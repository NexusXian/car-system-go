package main

import (
	"car-system-go/setup"
)

func init() {
	setup.InitViper()
	setup.InitMySQL()
}

func main() {
	r := setup.RouterInit()
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
