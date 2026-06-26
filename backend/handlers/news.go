package handlers

import (
	"net/http"
	"strconv"

	"xinhang-backend/database"
	"xinhang-backend/models"

	"github.com/gin-gonic/gin"
)

func ListNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "12"))
	category := c.Query("category")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}

	query := database.DB.Model(&models.News{}).Where("published = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var total int64
	query.Count(&total)

	var news []models.News
	if err := query.Select("id, title, summary, category, cover_image, author_name, source, published, published_at, created_at").
		Order("published_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取新闻列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     news,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func GetNews(c *gin.Context) {
	id := c.Param("id")

	var news models.News
	if err := database.DB.Where("id = ? AND published = ?", id, true).First(&news).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "新闻不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": news})
}

func GetCategories(c *gin.Context) {
	type CategoryCount struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	var results []CategoryCount
	database.DB.Model(&models.News{}).
		Select("category, count(*) as count").
		Where("published = ?", true).
		Group("category").
		Order("count DESC").
		Find(&results)

	c.JSON(http.StatusOK, gin.H{"data": results})
}

func CreateNews(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数格式错误"})
		return
	}

	if news.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "标题不能为空"})
		return
	}

	if err := database.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "data": news})
}

func UpdateNews(c *gin.Context) {
	id := c.Param("id")

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "新闻不存在"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数格式错误"})
		return
	}

	if err := database.DB.Model(&news).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}

	database.DB.First(&news, id)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "data": news})
}

func DeleteNews(c *gin.Context) {
	id := c.Param("id")

	var news models.News
	if err := database.DB.First(&news, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "新闻不存在"})
		return
	}

	if err := database.DB.Delete(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
