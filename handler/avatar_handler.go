package handler

import (
	"car-system-go/database"
	"car-system-go/model"
	"car-system-go/request"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MaxFileSizeMB = 10
	UploadDir     = "images/avatars" // 上传目录
)

// AdminUploadAvatarsHandler 处理管理员头像上传
func AdminUploadAvatarsHandler(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req request.AdminAvatarUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("请求参数错误: %v", err.Error()),
		})
		return
	}
	adminID := req.AdminID

	// 2. 查找用户
	var admin model.Admin
	if result := database.DB.Where("admin_id = ?", adminID).First(&admin); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("查询用户失败: %v", result.Error.Error()),
		})
		return
	}

	// 3. 获取上传的头像文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("获取表单数据失败: %v", err.Error()),
		})
		return
	}

	files := form.File["avatar"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供头像文件"})
		return
	}
	if len(files) > 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "一次只能上传一张头像图片"})
		return
	}
	file := files[0]

	// 4. 验证文件大小
	maxSize := int64(MaxFileSizeMB * 1024 * 1024)
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("文件过大，最大允许 %dMB", MaxFileSizeMB),
		})
		return
	}

	// 5. 验证文件类型（同时检查扩展名和MIME类型）
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	allowedMimes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	// 检查扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件类型无效，只允许 JPG, JPEG, PNG, GIF 格式",
		})
		return
	}

	// 检查MIME类型
	fileHeader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("打开文件失败: %v", err.Error()),
		})
		return
	}
	defer fileHeader.Close()

	buffer := make([]byte, 512)
	_, err = fileHeader.Read(buffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("读取文件信息失败: %v", err.Error()),
		})
		return
	}

	mimeType := http.DetectContentType(buffer)
	if !allowedMimes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "文件内容不符合图片格式要求",
		})
		return
	}

	// 6. 确保上传目录存在
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("创建上传目录失败: %v", err.Error()),
		})
		return
	}

	// 7. 生成新文件名并保存
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(UploadDir, newFilename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("保存文件失败: %v", err.Error()),
		})
		return
	}
	publicPath := "/" + filePath

	// 8. 删除旧头像文件（如果存在且不同）
	if admin.AvatarUrl != "" && admin.AvatarUrl != publicPath {
		oldFilePath := strings.TrimPrefix(admin.AvatarUrl, "/")
		// 检查文件是否存在且是常规文件
		if info, err := os.Stat(oldFilePath); err == nil && !info.IsDir() {
			if err := os.Remove(oldFilePath); err != nil {
				fmt.Printf("警告: 删除旧头像失败 %s: %v\n", oldFilePath, err)
			} else {
				fmt.Printf("已删除旧头像: %s\n", oldFilePath)
			}
		} else if !os.IsNotExist(err) {
			fmt.Printf("检查旧头像文件时出错 %s: %v\n", oldFilePath, err)
		}
	}

	// 9. 更新数据库中的头像URL（注意字段名与结构体一致）
	if result := database.DB.Model(&admin).Update("AvatarUrl", publicPath); result.Error != nil {
		// 更新失败时清理刚上传的文件
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("警告: 回滚头像文件失败 %s: %v\n", filePath, err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("更新头像信息失败: %v", result.Error.Error()),
		})
		return
	}

	// 10. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"code":    http.StatusOK,
		"data": gin.H{
			"avatarUrl": publicPath,
		},
	})
}
