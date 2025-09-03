package service

import (
	"car-system-go/database"
	"car-system-go/model"
	"car-system-go/repository"
	"car-system-go/request"
)

func InfractionRecordFindByIDCardNumberService(req request.InfractionRecordFindByIDCardNumber) (record *[]model.InfractionRecord, err error) {
	result, err := repository.InfractionRecordFindByIDCardNumber(req.IDCardNumber)
	if err != nil {
		return nil, err
	}
	return result, err
}
func InfractionRecordFindByRealNameService(req request.InfractionRecordFindByRealName) (record *[]model.InfractionRecord, err error) {
	result, err := repository.InfractionRecordFindByIDCardNumber(req.RealName)
	if err != nil {
		return nil, err
	}
	return result, err
}

func InfractionRecordFindAllService() (record *[]model.InfractionRecord, err error) {
	result, err := repository.InfractionRecordFindAll()
	if err != nil {
		return nil, err
	}
	return result, err
}

func GetLatestThreeInfractionRecordsService(analyzeRequest request.AIAnalyzeRequest) ([]model.InfractionRecord, error) {
	var records []model.InfractionRecord

	result := database.DB.Where("id_card_number = ?", analyzeRequest.IDCardNumber).
		Order("create_at DESC").
		Limit(3).
		Find(&records)

	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}
