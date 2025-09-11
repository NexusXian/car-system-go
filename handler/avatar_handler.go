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
	MaxFileSizeMB = 2
	UploadDir     = "images/avatars" // 上传目录，根据你的项目调整
)

func UploadAvatars(c *gin.Context) {
	// 1. 获取 email
	var req request.AdminAvatarUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("请求参数错误: %v", err.Error())})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("查询用户失败: %v", result.Error.Error())})
		return
	}

	// 3. 获取头像文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("获取 multipart 表单失败: %v", err.Error())})
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

	// 4. 验证文件大小和类型
	if file.Size > MaxFileSizeMB*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("文件过大，最大允许 %dMB", MaxFileSizeMB)})
		return
	}
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件类型无效。只允许 JPG, JPEG, PNG, GIF 格式。"})
		return
	}

	// 5. 生成新文件名并保存
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(UploadDir, newFilename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存文件失败: %v", err.Error())})
		return
	}
	publicPath := "/" + filePath

	// 6. 删除旧头像文件（如果存在）
	if admin.AvatarUrl != "" && admin.AvatarUrl != publicPath {
		oldAvatarFilePath := strings.TrimPrefix(admin.AvatarUrl, "/")
		if _, err := os.Stat(oldAvatarFilePath); err == nil {
			if err := os.Remove(oldAvatarFilePath); err != nil {
				fmt.Printf("删除旧头像失败: %v\n", err)
			} else {
				fmt.Printf("已删除旧头像: %s\n", oldAvatarFilePath)
			}
		}
	}

	if result := database.DB.Model(&admin).Select("AvatarUrl").Updates(map[string]interface{}{
		"avatar": publicPath,
	}); result.Error != nil {
		_ = os.Remove(strings.TrimPrefix(publicPath, "/")) // 清理新文件
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("更新用户头像路径失败: %v", result.Error.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "头像上传成功",
		"code":    http.StatusOK,
		"data": gin.H{
			"avatar": publicPath,
		},
	})
}
