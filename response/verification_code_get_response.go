package response

type VerificationCodeGetResponse struct {
	AdminID          string `json:"adminID"`
	VerificationCode string `json:"verificationCode"`
}
