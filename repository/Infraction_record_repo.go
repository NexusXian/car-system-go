package repository

import (
	"car-system-go/database"
	"car-system-go/model"
	"errors"
)

func InfractionRecordInsert(record *model.InfractionRecord) error {
	if err := database.DB.Create(record).Error; err != nil {
		return errors.New("违规记录创建失败")
	}
	return nil
}

func InfractionRecordFindByIDCardNumber(IDCarNumber string) (records *[]model.InfractionRecord, err error) {
	var result []model.InfractionRecord
	// 根据身份证号查询所有相关记录
	if err = database.DB.Where("id_card_number = ?", IDCarNumber).Find(&result).Error; err != nil {
		return nil, errors.New("查询违规记录失败：" + err.Error())
	}
	// 即使没有查询到记录，也返回空切片（而非nil），方便调用方处理
	return &result, nil
}

func InfractionRecordFindAll() (records *[]model.InfractionRecord, err error) {
	var result []model.InfractionRecord
	// 去掉Where条件，查询所有记录
	if err = database.DB.Find(&result).Error; err != nil {
		return nil, errors.New("查询违规记录失败：" + err.Error())
	}
	// 返回所有记录的指针，即使为空也返回空切片而非nil
	return &result, nil
}
