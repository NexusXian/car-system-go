package handler

import (
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BirthdayHandler(c *gin.Context) {
	var req request.UserBirthDayRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	username, day, isComing, err := service.BirthdayService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	message := fmt.Sprintf("距离%v的生日还有%v天，请提前为%v准备生日祝福", username, day, username)

	if isComing {
		c.JSON(http.StatusOK, utils.Response{
			Code:    http.StatusOK,
			Message: message,
		})
	}
}
