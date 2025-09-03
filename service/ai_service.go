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
