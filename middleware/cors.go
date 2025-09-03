package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允许所有来源的请求
		AllowOrigins: []string{"*"},
		// 允许所有HTTP方法
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		// 允许所有请求头
		AllowHeaders: []string{"*"},
		// 暴露所有响应头
		ExposeHeaders: []string{"*"},
		// 允许携带认证信息(cookie等)
		AllowCredentials: true,
		// 预检请求的缓存时间
		MaxAge: 12 * time.Hour,
	})
}
