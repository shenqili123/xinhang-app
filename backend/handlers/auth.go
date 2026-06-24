package handlers

import (
	"net/http"
	"time"

	"xinhang-backend/database"
	"xinhang-backend/email"
	"xinhang-backend/middleware"
	"xinhang-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

func generateToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func userResponse(user *models.User) gin.H {
	return gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"phone":         user.Phone,
		"role":          user.Role,
		"emailVerified": user.EmailVerified,
		"createdAt":     user.CreatedAt,
	}
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请填写完整的注册信息"})
		return
	}

	if req.Code != "" {
		ok, msg := verifyCode(req.Email, req.Code)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": msg})
			return
		}
	} else if email.IsEnabled() {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入邮箱验证码"})
		return
	}

	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "该邮箱已被注册"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "服务器内部错误"})
		return
	}

	user := models.User{
		Name:          req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		PasswordHash:  string(hash),
		Role:          "user",
		EmailVerified: req.Code != "",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败，请稍后重试"})
		return
	}

	tokenString, err := generateToken(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "注册成功，请登录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"token":   tokenString,
		"user":    userResponse(&user),
	})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入邮箱和密码"})
		return
	}

	ip := c.ClientIP()
	locked, lockMsg := middleware.CheckLoginLock(ip, req.Email)
	if locked {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": lockMsg})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		middleware.RecordLoginFail(ip, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "邮箱或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		middleware.RecordLoginFail(ip, req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "邮箱或密码错误"})
		return
	}

	middleware.ClearLoginFail(ip, req.Email)

	tokenString, err := generateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   tokenString,
		"user":    userResponse(&user),
	})
}

func GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userResponse(&user)})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数格式错误"})
		return
	}

	updates := map[string]interface{}{}
	if name, ok := req["name"]; ok && name != "" && len(name) <= 100 {
		updates["name"] = name
	}
	if phone, ok := req["phone"]; ok && phone != "" && len(phone) <= 20 {
		updates["phone"] = phone
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "没有需要更新的信息",
			"user":    userResponse(&user),
		})
		return
	}

	if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新失败"})
		return
	}

	database.DB.First(&user, userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    userResponse(&user),
	})
}

func MyApplications(c *gin.Context) {
	userID, _ := c.Get("userID")

	var apps []models.Application
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取报名列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    apps,
		"total":   len(apps),
	})
}
