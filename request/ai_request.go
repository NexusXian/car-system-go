package request

type AIAnalyzeRequest struct {
	IDCardNumber string `json:"IDCardNumber"`
}

type AIAnswerRequest struct {
	Question string `json:"question"`
}
