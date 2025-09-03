package handler

import (
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AiAnalyzeHandler(c *gin.Context) {
	var ai request.AIAnalyzeRequest
	if err := c.ShouldBindJSON(&ai); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	user, err := repository.UserFindByIDCardNumber(ai.IDCardNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	records, err := service.GetLatestThreeInfractionRecordsService(ai)

	message, err := service.AiAnalyzeService(*user, records)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code: http.StatusOK,
		Data: gin.H{
			"content": message,
		},
		Message: "ai分析完成",
	})

}
