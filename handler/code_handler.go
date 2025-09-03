package handler

import (
	"car-system-go/request"
	"car-system-go/response"
	"car-system-go/service"
	"car-system-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func VerificationCodeGetHandler(c *gin.Context) {
	var req request.VerificationCodeGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if req.AdminID == "" {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: "工号不能为空！",
		})
		return
	}

	if err := service.VerificationCodeGetService(req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: "验证码生成失败",
		})
		return
	}

	resp := response.VerificationCodeGetResponse{
		AdminID:          req.AdminID,
		VerificationCode: utils.CaptchaStore[req.AdminID],
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Message: "生成验证码成功！",
		Data:    resp,
	})
}
