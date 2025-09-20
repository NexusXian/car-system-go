package service

import (
	"car-system-go/llm"
	"car-system-go/model"
)

func AiAnalyzeService(user model.User, records []model.InfractionRecord) (string, error) {
	message, err := llm.AIModel(user, records)
	if err != nil {
		return "", err
	}
	return message, nil
}
func AiAnalyzeStreamService(user model.User, records []model.InfractionRecord) (<-chan string, error) {
	return llm.AIModelStream(user, records)
}

func AiQuestionService(records []model.InfractionRecord, message string) (<-chan string, error) {
	return llm.AIQuestionStream(records, message)
}

func DriverClassifyService(records []model.InfractionRecord) (string, error) {
	message, err := llm.DriverClassify(records)
	if err != nil {
		return "", err
	}
	return message, nil
}

func AiReportService(records []model.InfractionRecord) (llm.ReportResult, error) {
	message, err := llm.AIReport(records)
	if err != nil {
		return llm.ReportResult{}, err
	}
	return message, nil
}

func AiReportUserService(records []model.InfractionRecord) (llm.AnalysisResult, error) {
	message, err := llm.AIReportUser(records)
	if err != nil {
		return llm.AnalysisResult{}, err
	}
	return message, nil
}
