package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type User struct {
	UID                string      `json:"UID"`
	RealName           string      `json:"realName" comment:"真实姓名"`
	HireDate           string      `json:"hireDate" comment:"入职时间"`
	DrivingExperience  int         `json:"drivingExperience" comment:"驾龄"`
	IDCardNumber       string      `json:"IDCardNumber" comment:"身份证"`
	LicensePlate       string      `json:"licensePlate" comment:"车牌"`
	BloodType          string      `json:"BloodType" comment:"血型"`
	ResidentialAddress string      `json:"ResidentialAddress" comment:"居住地址"`
	EmergencyContact   string      `json:"emergencyContact" comment:"紧急联系人"`
	Allergies          string      `json:"allergies" comment:"过敏症"`
	IsOrganDonor       bool        `json:"IsOrganDonor" comment:"器官捐赠者 (是/否)"`
	MedicalNotes       string      `json:"MedicalNotes" comment:"医疗注意事项"`
	Certificates       StringSlice `json:"certificates" comment:"技能证书展示" gorm:"type:text"`
	FamilyBrief        string      `json:"familyBrief" comment:"家庭情况简要记录"`
	Subsidy            int         `json:"subsidy" comment:"津贴" gorm:"default:1000"`
	InfractionCount    int         `json:"infractionCount" comment:"违规次数"`
	OxygenSaturation   float64     `json:"oxygenSaturation" comment:"血氧"`
	HeartRate          float64     `json:"heartRate" comment:"心率"`
	BodyTemperature    float64     `json:"bodyTemperature" comment:"体温"`
	CreateAt           time.Time   `gorm:"column:create_at;autoCreateTime"`
	UpdateAt           time.Time   `gorm:"column:update_at;autoUpdateTime"`
}

func (user *User) TableName() string {
	return "user_tb"
}

type StringSlice []string

// Value：将 StringSlice 转为 JSON 字符串，存入数据库（实现 driver.Valuer 接口）
func (s StringSlice) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan：从数据库读取 JSON 字符串，反序列化为 StringSlice（实现 sql.Scanner 接口）
func (s *StringSlice) Scan(value interface{}) error {
	// 处理空值情况
	if value == nil {
		*s = nil
		return nil
	}
	// 将数据库返回的字节流转为 JSON
	data, ok := value.([]byte)
	if !ok {
		return errors.New("invalid data type for StringSlice")
	}
	return json.Unmarshal(data, s)
}
