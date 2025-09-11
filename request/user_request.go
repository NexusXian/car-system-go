package request

type UserRegisterRequest struct {
	RealName           string   `json:"realName" comment:"真实姓名"`
	HireDate           string   `json:"hireDate" comment:"入职时间"`
	DrivingExperience  int      `json:"drivingExperience" comment:"驾龄"`
	IDCardNumber       string   `json:"IDCardNumber" comment:"身份证"`
	LicensePlate       string   `json:"licensePlate" comment:"车牌"`
	BloodType          string   `json:"BloodType" comment:"血型"`
	ResidentialAddress string   `json:"ResidentialAddress" comment:"居住地址"`
	EmergencyContact   string   `json:"emergencyContact" comment:"紧急联系人"`
	Allergies          string   `json:"allergies" comment:"过敏症"`
	IsOrganDonor       bool     `json:"IsOrganDonor" comment:"器官捐赠者 (是/否)"`
	MedicalNotes       string   `json:"MedicalNotes" comment:"医疗注意事项"`
	Certificates       []string `json:"certificates" comment:"技能证书展示"`
	FamilyBrief        string   `json:"familyBrief" comment:"家庭情况简要记录"`
}
type UserInfractionCreateRequest struct {
	IDCardNumber string `json:"idCardNumber"`
	RealName     string `json:"realName"`
	LicensePlate string `json:"licensePlate"`
	Record       string `json:"record"`
}

type UserBirthDayRequest struct {
	IDCardNumber string `json:"idCardNumber"`
}

type UserFindRequest struct {
	IDCardNumber string `json:"idCardNumber"`
}
