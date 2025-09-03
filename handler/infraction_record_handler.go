package handler

import (
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InfractionRecordFindByIDCardNumberHandler(c *gin.Context) {
	var req request.InfractionRecordFindByIDCardNumber
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	records, err := service.InfractionRecordFindByIDCardNumberService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Data:    records,
		Message: "查询记录成功",
	})
}

func InfractionRecordFindByRealNameNumberHandler(c *gin.Context) {
	var req request.InfractionRecordFindByRealName
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	records, err := service.InfractionRecordFindByRealNameService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Data:    records,
		Message: "查询记录成功",
	})
}

func InfractionRecordFindAllHandler(c *gin.Context) {
	records, err := service.InfractionRecordFindAllService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Data:    records,
		Message: "查询成功",
	})
}
