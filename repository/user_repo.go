package repository

import (
	"car-system-go/database"
	"car-system-go/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// mysql gorm
func UserInsert(user *model.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		return errors.New("用户创建失败！")
	}
	return nil
}

func UserFindByRealName(realName string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("real_name = ?", realName).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

func UserFindByID(uid string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

func UserUpdate(user *model.User) error {
	if err := database.DB.Where("id_card_number = ?", user.IDCardNumber).Updates(user).Error; err != nil {
		return errors.New("用户信息更新失败: " + err.Error())
	}

	return nil
}

func UserFindByIDCardNumber(idCarNumber string) (*model.User, error) {
	var user model.User
	if err := database.DB.Where("id_card_number = ?", idCarNumber).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return &user, nil
}

func UserFindAll() ([]*model.User, error) {
	var users []*model.User
	if err := database.DB.Find(&users).Error; err != nil {
		return users, errors.New("查询失败，请稍后重试！")
	}
	return users, nil
}

func UserFindAllASC() ([]*model.User, error) {
	var users []*model.User
	// 使用Order指定按infraction_count升序排序（ASC可省略，默认即为升序）
	if err := database.DB.Order("infraction_count ASC").Find(&users).Error; err != nil {
		return users, errors.New("查询失败，请稍后重试！")
	}
	return users, nil
}

//redis

// 从缓存获取用户
func GetUserFromCache(id string) (*model.User, error) {
	key := fmt.Sprintf("user:%s", id)
	val, err := database.RDB.Get(database.CTX, key).Result()
	if err != nil {
		return nil, err
	}

	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 将用户存入缓存
func SetUserToCache(user *model.User) error {
	key := fmt.Sprintf("user:%d", user.UID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// 设置缓存，过期时间1小时
	return database.RDB.Set(database.CTX, key, data, time.Hour).Err()
}

// 从数据库获取用户
func GetUserFromDB(id string) (*model.User, error) {
	var user model.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// 清除用户缓存
func DeleteUserCache(id string) error {
	key := fmt.Sprintf("user:%s", id)
	return database.RDB.Del(database.CTX, key).Err()
}
