package service

import (
	"car-system-go/repository"
	"car-system-go/request"
	"errors"
	"fmt"
	"time"
)

func BirthdayService(req request.UserBirthDayRequest) (string, int, bool, error) {
	// 获取当前日期
	now := time.Now()
	currentYear := now.Year()

	// 根据身份证号查询用户信息
	user, err := repository.UserFindByIDCardNumber(req.IDCardNumber)
	if err != nil {
		return "", 0, false, err
	}

	// 从身份证号提取生日（第7-14位，格式YYYYMMDD）
	if len(user.IDCardNumber) < 14 {
		return "", 0, false, errors.New("无效的身份证号码")
	}
	birthStr := user.IDCardNumber[6:14] // 提取出生日期部分

	// 解析生日字符串为时间对象
	birthDate, err := time.Parse("20060102", birthStr)
	if err != nil {
		return "", 0, false, fmt.Errorf("解析生日失败: %w", err)
	}

	// 构建今年的生日日期
	thisYearBirthday := time.Date(currentYear, birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, now.Location())

	// 处理跨年情况（如果今年生日已过，则计算明年的）
	if thisYearBirthday.Before(now) {
		thisYearBirthday = thisYearBirthday.AddDate(1, 0, 0)
	}

	// 计算天数差
	daysDiff := int(thisYearBirthday.Sub(now).Hours() / 24)

	// 判断是否在3天或以内
	return user.RealName, daysDiff, daysDiff <= 3, nil
}
