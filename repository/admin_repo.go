package repository

import (
	"car-system-go/database"
	"car-system-go/model"
	"errors"
)

func AdminInsert(admin *model.Admin) error {
	if err := database.DB.Create(&admin).Error; err != nil {
		return errors.New("管理员创建失败！")
	}
	return nil
}

func AdminFindByAdminID(adminID string) (*model.Admin, error) {
	var admin model.Admin
	if err := database.DB.Where("admin_id = ?", adminID).First(&admin).Error; err != nil {
		return nil, errors.New("管理员工号" + adminID + "不存在")
	}
	return &admin, nil
}
