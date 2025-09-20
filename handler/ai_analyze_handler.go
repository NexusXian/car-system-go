package handler

import (
	"car-system-go/repository"
	"car-system-go/request"
	"car-system-go/response"
	"car-system-go/service"
	"car-system-go/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

// DriverClassifyHandler HTTP处理器
func DriverClassifyHandler(c *gin.Context) {
	// 查询所有违规记录
	records, err := repository.InfractionRecordFindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("获取违规记录失败: %v", err),
		})
		return
	}

	// 检查是否有记录
	if records == nil || len(*records) == 0 {
		c.JSON(http.StatusOK, utils.Response{
			Code:    http.StatusOK,
			Message: "success",
			Data: response.DriverClassificationResult{
				ClassificationTime: time.Now().Format("2006-01-02 15:04:05"),
				TotalDrivers:       0,
				Categories: []struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				}{
					{Type: "安全型", Count: 0},
					{Type: "激进型", Count: 0},
					{Type: "疲劳型", Count: 0},
					{Type: "未知型", Count: 0},
				},
			},
		})
		return
	}

	// 调用服务层进行司机分类
	result, err := service.DriverClassifyService(*records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("AI分类计算失败: %v", err),
		})
		return
	}

	// 解析 AI 输出的 JSON
	var data response.DriverClassificationResult
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		// 如果解析失败，尝试解析为通用interface{}
		var fallbackData interface{}
		if parseErr := json.Unmarshal([]byte(result), &fallbackData); parseErr != nil {
			c.JSON(http.StatusInternalServerError, utils.Response{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("解析AI输出失败: %v, 原始响应: %s", err, result),
			})
			return
		}
		data = response.DriverClassificationResult{
			ClassificationTime: time.Now().Format("2006-01-02 15:04:05"),
			TotalDrivers:       0,
			Categories: []struct {
				Type  string `json:"type"`
				Count int    `json:"count"`
			}{
				{Type: "解析错误", Count: 0},
			},
		}
	}

	// 成功返回
	c.JSON(http.StatusOK, utils.Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func AiReportHandler(c *gin.Context) {
	records, err := repository.InfractionRecordFindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	res, err := service.AiReportService(*records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func AiReportUser(c *gin.Context) {
	var req request.AIAnalyzeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	records, err := repository.InfractionRecordFindByIDCardNumber(req.IDCardNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	res, err := service.AiReportUserService(*records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}
