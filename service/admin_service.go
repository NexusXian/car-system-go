package service

import (
	"car-system-go/database"
	"car-system-go/jwt"
	"car-system-go/model"
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/response"
	"car-system-go/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

// 管理员注册逻辑
func AdminRegisterService(req request.AdminRegisterRequest) error {
	if _, err := repository.AdminFindByAdminID(req.AdminID); err == nil {
		return errors.New("管理员工号" + req.AdminID + "已存在！")
	}
	hashPWS, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.New("密码加密失败")
	}
	admin := &model.Admin{
		AdminID:  req.AdminID,
		Password: hashPWS,
	}
	if err := repository.AdminInsert(admin); err != nil {
		return err
	}
	return nil
}

//管理员登录

func AdminLoginService(req request.AdminLoginRequest) (*response.AdminLoginResponse, error) {
	admin, err := repository.AdminFindByAdminID(req.AdminID)
	if err != nil {
		return nil, errors.New("账号或密码错误")
	}
	ok := utils.CheckPassword(req.Password, admin.Password)
	if !ok {
		return nil, errors.New("账号或密码错误")
	}

	expirationTime, token, err := jwt.GenerateToken(admin.AdminID)
	if err != nil {
		return nil, fmt.Errorf("token生成失败: %w", err)
	}

	resp := &response.AdminLoginResponse{
		Token:      token,
		ExpireTime: expirationTime,
		AdminInfo: response.AdminInfo{
			WorkID:    admin.AdminID,
			AvatarUrl: admin.AvatarUrl,
		},
		TimeStamp: time.Now().Unix(),
	}
	return resp, nil
}

//管理员找回密码

func AdminFindPasswordService(req request.AdminFindPasswordRequest) error {
	admin, err := repository.AdminFindByAdminID(req.AdminID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(req.VerificationCode, utils.CaptchaStore[req.AdminID]) {
		return errors.New("验证码错误")
	}
	hashPwd, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	admin.Password = hashPwd
	if err = database.DB.Save(&admin).Error; err != nil {
		return err
	}
	delete(utils.CaptchaStore, req.AdminID)
	return nil
}
