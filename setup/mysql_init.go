package setup

import (
	"car-system-go/database"
	"car-system-go/model"
)

func InitMySQL() {
	err := database.InitMySQL()
	if err != nil {
		panic(err)
	}
	err = database.DB.AutoMigrate(&model.Admin{}, &model.User{}, &model.InfractionRecord{}, &model.EmbeddedWatch{},
		&model.AlcoholModule{}, &model.CollisionModule{}, &model.SmokeModule{}, &model.CarComputer{})
	if err != nil {
		panic("automigrate error " + err.Error())
	}

}
