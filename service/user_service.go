package service

import (
	"car-system-go/model"
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/utils"
	"errors"
)

func UserRegisterService(req request.UserRegisterRequest) error {
	if _, err := repository.UserFindByIDCardNumber(req.IDCardNumber); err == nil {
		return errors.New("身份证已经存在！")
	}

	user := &model.User{
		RealName:           req.RealName,
		HireDate:           req.HireDate,
		DrivingExperience:  req.DrivingExperience,
		IDCardNumber:       req.IDCardNumber,
		LicensePlate:       req.LicensePlate,
		BloodType:          req.BloodType,
		ResidentialAddress: req.ResidentialAddress,
		EmergencyContact:   req.EmergencyContact,
		Allergies:          req.Allergies,
		IsOrganDonor:       req.IsOrganDonor,
		MedicalNotes:       req.MedicalNotes,
		Certificates:       req.Certificates,
		FamilyBrief:        req.FamilyBrief,
	}
	if err := repository.UserInsert(user); err != nil {
		return err
	}
	return nil
}

func UserFindAllService() ([]*model.User, error) {
	users, err := repository.UserFindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 用户违规了
func UserInfractionCreateService(req request.UserInfractionCreateRequest) error {
	// 查询用户是否存在
	user, err := repository.UserFindByIDCardNumber(req.IDCardNumber)
	if err != nil {
		// 如果查询出错（通常是未找到用户），返回错误
		return errors.New("当前用户 " + req.IDCardNumber + " 不存在")
	}

	// 创建违规记录
	record := &model.InfractionRecord{
		IDCardNumber: req.IDCardNumber,
		RealName:     req.RealName,
		Record:       req.Record,
		LicensePlate: req.LicensePlate,
	}

	if err := repository.InfractionRecordInsert(record); err != nil {
		return err
	}

	// 更新用户的违规次数和补贴
	user.InfractionCount++
	user.Subsidy = utils.CalculateSubsidy(user.InfractionCount)
	if err := repository.UserUpdate(user); err != nil {
		return err
	}

	return nil
}
