package service

import (
	"car-system-go/model"
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/response"
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

	if req.Record == "driving-tired" {
		req.Record = "疲劳驾驶"
	} else if req.Record == "driving-sleep" {
		req.Record = "开车睡觉"
	} else if req.Record == "driving-hand without wheel" {
		req.Record = "开车时未手握方向盘"
	} else if req.Record == "driving-call" {
		req.Record = "开车接打电话"
	} else if req.Record == "driving-smok" {
		req.Record = "开车抽烟"
	}
	//if (recordRequest.getRecord().equals("driving-tired")){
	//	record.setRecord("疲劳驾驶");
	//} else if (recordRequest.getRecord().equals("driving-sleep")) {
	//	record.setRecord("开车睡觉");
	//} else if (recordRequest.getRecord().equals("driving-hand without wheel")) {
	//	record.setRecord("开车时未手握方向盘");
	//} else if(recordRequest.getRecord().equals("driving-call")) {
	//	record.setRecord("开车接打电话");
	//} else if (recordRequest.getRecord().equals("driving-smok")) {
	//	record.setRecord("开车抽烟");
	//} else {
	//	record.setRecord("未知违规类型");
	//}
	//// 创建违规记录
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

func UserFindAllInfoService() ([]*response.UserBasicResponse, error) {
	users, err := repository.UserFindAllASC()

	var result []*response.UserBasicResponse
	for _, user := range users {
		result = append(result, &response.UserBasicResponse{
			RealName:        user.RealName,
			IDCardNumber:    user.IDCardNumber,
			InfractionCount: user.InfractionCount,
		})
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UserFindService(req request.UserFindRequest) (*model.User, error) {
	user, err := repository.UserFindByIDCardNumber(req.IDCardNumber)
	if err != nil {
		return nil, err
	}
	return user, nil
}
