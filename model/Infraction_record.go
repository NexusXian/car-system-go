package model

import (
	"time"
)

type InfractionRecord struct {
	ID           uint64    `json:"id" gorm:"primary_key;auto_increment"`
	RealName     string    `json:"realName" `
	LicensePlate string    `json:"licensePlate"`
	IDCardNumber string    `json:"IDCardNumber" comment:"身份证"`
	Record       string    `json:"record"`
	CreateAt     time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdateAt     time.Time `gorm:"column:update_at;autoUpdateTime"`
}

func (InfractionRecord) TableName() string {
	return "infraction_record_table"
}
