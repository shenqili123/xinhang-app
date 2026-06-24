package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"xinhang-backend/database"
	"xinhang-backend/models"

	"github.com/gin-gonic/gin"
)

func Apply(c *gin.Context) {
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("表单验证失败: %s", err.Error())})
		return
	}

	now := time.Now()
	permitNo := fmt.Sprintf("XH%s%04d", now.Format("20060102"), now.UnixMilli()%10000)

	app := models.Application{
		PermitNo:      permitNo,
		StudentName:   req.StudentName,
		BirthDate:     req.BirthDate,
		Gender:        req.Gender,
		Grade:         req.Grade,
		IDNumber:      req.IDNumber,
		BoardingNeed:  req.BoardingNeed,
		ParentName:    req.ParentName,
		Phone:         req.Phone,
		Relationship:  req.Relationship,
		Email:         req.Email,
		CurrentSchool: req.CurrentSchool,
		Track:         req.Track,
		VisitDate:     req.VisitDate,
		Notes:         req.Notes,
		Status:        "pending",
	}

	userID, _ := c.Get("userID")
	uid := userID.(uint)
	app.UserID = &uid

	var count int64
	database.DB.Model(&models.Application{}).
		Where("phone = ? AND student_name = ? AND grade = ?", app.Phone, app.StudentName, app.Grade).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "该学生已提交过报名申请，请勿重复提交"})
		return
	}

	if err := database.DB.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "提交失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "报名申请已提交成功",
		"id":       app.ID,
		"permitNo": app.PermitNo,
	})
}

func QueryPermit(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入报名号或手机号"})
		return
	}

	var app models.Application
	err := database.DB.Where("permit_no = ? OR phone = ?", keyword, keyword).
		Order("created_at DESC").First(&app).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"found": false, "message": "未找到匹配的报名记录，请检查输入"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"found":     true,
		"permitNo":  app.PermitNo,
		"student":   app.StudentName,
		"grade":     app.Grade,
		"school":    app.CurrentSchool,
		"status":    app.Status,
		"createdAt": app.CreatedAt,
	})
}

func ListApplications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	database.DB.Model(&models.Application{}).Count(&total)

	var apps []models.Application
	if err := database.DB.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取报名列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"data":     apps,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
