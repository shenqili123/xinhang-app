package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请上传照片文件"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "照片文件不能超过5MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"message": "仅支持 JPG/PNG 格式的照片"})
		return
	}

	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建上传目录失败"})
		return
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), randomString(8), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "保存文件失败"})
		return
	}

	url := "/uploads/photos/" + filename
	c.JSON(http.StatusOK, gin.H{
		"message": "上传成功",
		"url":     url,
	})
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1)
	}
	return string(b)
}
