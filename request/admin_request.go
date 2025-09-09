package request

type AdminFindPasswordRequest struct {
	AdminID          string `json:"adminID"`
	NewPassword      string `json:"newPassword"`
	VerificationCode string `json:"verificationCode"`
}

type AdminLoginRequest struct {
	AdminID  string `json:"adminID"`
	Password string `json:"password"`
}
type AdminRegisterRequest struct {
	AdminID  string `json:"adminID"`
	Password string `json:"password"`
}
type AdminAvatarUploadRequest struct {
	AdminID string `form:"adminID" binding:"required,adminID"` // 从 form-data 中获取 email 字段
}
