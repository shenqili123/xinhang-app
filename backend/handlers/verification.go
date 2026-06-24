package handlers

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"log"

	"xinhang-backend/cache"
	"xinhang-backend/email"

	"github.com/gin-gonic/gin"
)

type sendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func generateCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return fmt.Sprintf("%06d", n.Int64())
}

func SendVerificationCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请输入有效的邮箱地址"})
		return
	}

	if !email.IsEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "邮件服务暂未配置"})
		return
	}

	if cache.RDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"message": "验证码服务暂不可用"})
		return
	}

	cdKey := fmt.Sprintf("verify_cd:%s", req.Email)
	if cache.RDB.Exists(cache.Ctx, cdKey).Val() > 0 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "发送过于频繁，请60秒后重试"})
		return
	}

	countKey := fmt.Sprintf("verify_count:%s", req.Email)
	cnt, _ := cache.RDB.Get(cache.Ctx, countKey).Int64()
	if cnt >= 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "发送次数过多，请1小时后重试"})
		return
	}

	code := generateCode()

	cache.RDB.Set(cache.Ctx, fmt.Sprintf("verify:%s", req.Email), code, 5*time.Minute)
	cache.RDB.Del(cache.Ctx, fmt.Sprintf("verify_attempts:%s", req.Email))

	cache.RDB.Set(cache.Ctx, cdKey, "1", 60*time.Second)

	pipe := cache.RDB.Pipeline()
	pipe.Incr(cache.Ctx, countKey)
	pipe.Expire(cache.Ctx, countKey, time.Hour)
	pipe.Exec(cache.Ctx)

	if err := email.SendVerificationCode(req.Email, code); err != nil {
		log.Printf("SMTP send error to %s: %v", req.Email, err)
		cache.RDB.Del(cache.Ctx, fmt.Sprintf("verify:%s", req.Email))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "邮件发送失败，请检查邮箱地址是否正确"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送，请查收邮件"})
}

var testCodes map[string]string

func SetVerifyCodeForTest(emailAddr, code string) {
	if testCodes == nil {
		testCodes = make(map[string]string)
	}
	testCodes[emailAddr] = code
}

func verifyCode(emailAddr, code string) (bool, string) {
	if testCodes != nil {
		if stored, ok := testCodes[emailAddr]; ok {
			delete(testCodes, emailAddr)
			if stored == code {
				return true, ""
			}
			return false, "验证码错误"
		}
	}

	if cache.RDB == nil {
		return false, "验证码服务暂不可用"
	}

	attKey := fmt.Sprintf("verify_attempts:%s", emailAddr)
	att, _ := cache.RDB.Get(cache.Ctx, attKey).Int64()
	if att >= 5 {
		return false, "验证码错误次数过多，请重新获取"
	}

	codeKey := fmt.Sprintf("verify:%s", emailAddr)
	stored, err := cache.RDB.Get(cache.Ctx, codeKey).Result()
	if err != nil {
		return false, "验证码已过期，请重新获取"
	}

	if stored != code {
		cache.RDB.Incr(cache.Ctx, attKey)
		cache.RDB.Expire(cache.Ctx, attKey, 5*time.Minute)
		return false, "验证码错误"
	}

	cache.RDB.Del(cache.Ctx, codeKey)
	cache.RDB.Del(cache.Ctx, attKey)
	return true, ""
}
