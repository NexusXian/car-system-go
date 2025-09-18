package handler

import (
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/service"
	"car-system-go/utils"
	"io"
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

func AiAnalyzeStreamHandler(c *gin.Context) {
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
		})
		return
	}

	records, err := service.GetLatestThreeInfractionRecordsService(ai)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 获取流式通道
	streamChan, err := service.AiAnalyzeStreamService(*user, records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	// 设置响应头为流式传输
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 流式推送数据
	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				// 通道关闭，结束流
				return false
			}
			// 向前端写入当前片段（可根据需要包装为 JSON 格式）
			if _, err := w.Write([]byte(chunk)); err != nil {
				return false
			}
			// 刷新缓冲区，确保前端及时收到
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			// 客户端断开连接，终止流
			return false
		}
	})
}

//ai问答这块

func AiAnswerHandler(c *gin.Context) {
	var ai request.AIAnswerRequest
	if err := c.ShouldBindJSON(&ai); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	records, err := repository.InfractionRecordFindAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	// 获取流式通道
	streamChan, err := service.AiQuestionService(*records, ai.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case chunk, ok := <-streamChan:
			if !ok {
				return false
			}
			if _, err := w.Write([]byte(chunk)); err != nil {
				return false
			}
			c.Writer.Flush()
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
