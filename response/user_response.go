package response

type UserBasicResponse struct {
	RealName        string `form:"real_name" json:"real_name" binding:"required"`
	IDCardNumber    string `form:"id_card_number" json:"idCardNumber" binding:"required"`
	InfractionCount int    `form:"infraction_count" json:"infractionCount" binding:"required"`
	OxygenSaturation   float64     `json:"oxygenSaturation" comment:"血氧"`
	HeartRate          float64     `json:"heartRate" comment:"心率"`
	BodyTemperature    float64     `json:"bodyTemperature" comment:"体温"`
}
