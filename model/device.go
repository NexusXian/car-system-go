package model

// 嵌入式手表模块
type EmbeddedWatch struct {
	Status       string  `json:"status"`       // 连接状态："正常"或"未连接"
	Temperature  float64 `json:"temperature"`  // 温度（单位：°C）
	Age          int     `json:"age"`          // 寿命预估（单位：年）
	IDCardNumber string  `json:"idCardNumber"` // 用户身份证号
	DeviceID     string  `json:"deviceID"`     // 设备唯一标识符
}

// 酒精模块
type AlcoholModule struct {
	Status       string  `json:"status"`       // 连接状态："正常"或"未连接"
	Temperature  float64 `json:"temperature"`  // 温度（单位：°C）
	Age          int     `json:"age"`          // 寿命预估（单位：年）
	Sensitivity  float64 `json:"sensitivity"`  //灵敏状态（单位：%）;
	IDCardNumber string  `json:"idCardNumber"` // 用户身份证号
	DeviceID     string  `json:"deviceID"`     // 设备唯一标识符
}

// 烟雾模块
type SmokeModule struct {
	Status       string  `json:"status"`       // 连接状态："正常"或"未连接"
	Temperature  float64 `json:"temperature"`  // 温度（单位：°C）
	Age          int     `json:"age"`          // 寿命预估（单位：年）
	Sensitivity  float64 `json:"sensitivity"`  //灵敏状态（单位：%）;
	IDCardNumber string  `json:"idCardNumber"` // 用户身份证号
	DeviceID     string  `json:"deviceID"`     // 设备唯一标识符
}

// 车载电脑模块
type CarComputer struct {
	Status       string  `json:"status"`       // 连接状态："正常"或"未连接"
	Temperature  float64 `json:"temperature"`  // 温度（单位：°C）
	Age          int     `json:"age"`          // 寿命预估（单位：年）
	IDCardNumber string  `json:"idCardNumber"` // 用户身份证号
	DeviceID     string  `json:"deviceID"`     // 设备唯一标识符
}

// 碰撞模块
type CollisionModule struct {
	Status       string  `json:"status"`       // 连接状态："正常"或"未连接"
	Temperature  float64 `json:"temperature"`  // 温度（单位：°C）
	Age          int     `json:"age"`          // 寿命预估（单位：年）
	Sensitivity  float64 `json:"sensitivity"`  //灵敏状态（单位：%）;
	IDCardNumber string  `json:"idCardNumber"` // 用户身份证号
	DeviceID     string  `json:"deviceID"`     // 设备唯一标识符
}

func (EmbeddedWatch) TableName() string {
	return "embedded_watch"
}
func (AlcoholModule) TableName() string {
	return "alcohol_module"
}
func (SmokeModule) TableName() string {
	return "smoke_module"
}
func (CarComputer) TableName() string {
	return "car_computer"
}
func (CollisionModule) TableName() string {
	return "collision_module"
}
