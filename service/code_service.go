package service

import (
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/utils"
)

func VerificationCodeGetService(req request.VerificationCodeGetRequest) error {
	if _, err := repository.AdminFindByAdminID(req.AdminID); err != nil {
		return err
	}
	code := utils.GenerateCaptcha()
	utils.CaptchaStore[req.AdminID] = code
	return nil
}
