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
		EmailVerified: email.IsEnabled(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "注册失败，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"role":   user.Role,
		"exp":    time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   tokenString,
	})
}
