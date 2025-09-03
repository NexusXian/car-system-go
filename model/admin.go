package model

import (
	"time"
)

type Admin struct {
	ID       uint      `gorm:"primary_key"`
	AdminID  string    `json:"adminID"`
	Password string    `json:"password"`
	CreateAt time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdateAt time.Time `gorm:"column:update_at;autoUpdateTime"`
}

func (Admin) TableName() string {
	return "admin_tb"
}
