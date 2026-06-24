package middleware

import (
	"net/http"
	"time"

	"xinhang-backend/cache"

	"github.com/gin-gonic/gin"
)

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		allowed, err := cache.CheckRateLimit(key, maxRequests, window)
		if err != nil {
			c.Next()
			return
		}
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"message": "请求过于频繁，请稍后再试"})
			c.Abort()
			return
		}
		c.Next()
	}
}
