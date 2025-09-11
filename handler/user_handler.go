package handler

import (
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegisterHandler(c *gin.Context) {
	var req request.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err := service.UserRegisterService(req); err != nil {
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

func UserFindAllHandler(c *gin.Context) {
	users, err := service.UserFindAllService()
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
		Message: "查询成功!",
		Data:    users,
	})
}
func UserInfractionCreateHandler(c *gin.Context) {
	var req request.UserInfractionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	if err := service.UserInfractionCreateService(req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Data:    nil,
		Message: "违规行为已记录",
	})
}

func UserFindHandler(c *gin.Context) {
	var req request.UserFindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	user, err := service.UserFindService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Message: "查询成功",
		Data:    user,
	})
}

func UserFindAllInfoHandler(c *gin.Context) {
	resp, err := service.UserFindAllInfoService()
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
		Message: "查询成功",
		Data:    resp,
	})
}
