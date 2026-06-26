package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"xinhang-backend/database"
	"xinhang-backend/models"

	"github.com/gin-gonic/gin"
)

var idNumberRegex = regexp.MustCompile(`^\d{17}[\dXx]$`)

func validateIDNumber(id string) bool {
	if len(id) != 18 {
		return false
	}
	if !idNumberRegex.MatchString(id) {
		return false
	}
	weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	checkCodes := "10X98765432"
	sum := 0
	for i := 0; i < 17; i++ {
		sum += int(id[i]-'0') * weights[i]
	}
	expected := checkCodes[sum%11]
	last := id[17]
	if last == 'x' {
		last = 'X'
	}
	return last == expected
}

func assignExamRoomAndSeat(grade int) (string, string) {
	var totalInGrade int64
	database.DB.Model(&models.Application{}).Where("grade = ?", grade).Count(&totalInGrade)
	seatsPerRoom := 30
	room := int(totalInGrade)/seatsPerRoom + 1
	seat := int(totalInGrade)%seatsPerRoom + 1
	if seat == 0 {
		seat = rand.Intn(30) + 1
	}
	return fmt.Sprintf("%02d", room), fmt.Sprintf("%02d", seat)
}

func Apply(c *gin.Context) {
	var req models.ApplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("表单验证失败: %s", err.Error())})
		return
	}

	if !validateIDNumber(req.IDNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "身份证号格式不正确，请输入18位有效身份证号码"})
		return
	}

	now := time.Now()
	permitNo := fmt.Sprintf("XH%s%04d", now.Format("20060102"), now.UnixMilli()%10000)
	examRoom, seatNumber := assignExamRoomAndSeat(req.Grade)

	app := models.Application{
		PermitNo:      permitNo,
		StudentName:   req.StudentName,
		BirthDate:     req.BirthDate,
		Gender:        req.Gender,
		Division:      req.Division,
		Grade:         req.Grade,
		IDNumber:      req.IDNumber,
		Photo:         req.Photo,
		ExamRoom:      examRoom,
		SeatNumber:    seatNumber,
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
		"message":    "报名申请已提交成功",
		"id":         app.ID,
		"permitNo":   app.PermitNo,
		"examRoom":   app.ExamRoom,
		"seatNumber": app.SeatNumber,
		"photo":      app.Photo,
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
