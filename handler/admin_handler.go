package handler

import (
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRegisterHandler(c *gin.Context) {
	var req request.AdminRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: "数据绑定失败",
			Data:    nil,
		})
		return
	}

	if err := service.AdminRegisterService(req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Message: "注册成功！",
		Data:    nil,
	})
}

func AdminLoginHandler(c *gin.Context) {
	var req request.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	resp, err := service.AdminLoginService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Data:    resp,
		Message: "登录成功！",
	})
}

func AdminFindPasswordHandler(c *gin.Context) {
	var req request.AdminFindPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := service.AdminFindPasswordService(req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Message: "密码更新成功！",
		Data:    nil,
	})
}
