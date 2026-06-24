package middleware

import (
	"fmt"
	"net/http"
	"time"

	"xinhang-backend/cache"

	"github.com/gin-gonic/gin"
)

func IPProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cache.RDB == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		blockKey := fmt.Sprintf("ip_block:%s", ip)

		if cache.RDB.Exists(cache.Ctx, blockKey).Val() > 0 {
			c.JSON(http.StatusForbidden, gin.H{"message": "您的IP已被临时封禁，请30分钟后重试"})
			c.Abort()
			return
		}

		countKey := fmt.Sprintf("ip_req:%s", ip)
		count, _ := cache.RDB.Incr(cache.Ctx, countKey).Result()
		if count == 1 {
			cache.RDB.Expire(cache.Ctx, countKey, time.Minute)
		}

		if count > 200 {
			cache.RDB.Set(cache.Ctx, blockKey, "1", 30*time.Minute)
			c.JSON(http.StatusForbidden, gin.H{"message": "请求异常，IP已被临时封禁"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckLoginLock(ip, emailAddr string) (bool, string) {
	if cache.RDB == nil {
		return false, ""
	}

	lockKey := fmt.Sprintf("login_lock:%s:%s", ip, emailAddr)
	if cache.RDB.Exists(cache.Ctx, lockKey).Val() > 0 {
		ttl := cache.RDB.TTL(cache.Ctx, lockKey).Val()
		mins := int(ttl.Minutes()) + 1
		return true, fmt.Sprintf("登录失败次数过多，请%d分钟后重试", mins)
	}
	return false, ""
}

func RecordLoginFail(ip, emailAddr string) {
	if cache.RDB == nil {
		return
	}

	failKey := fmt.Sprintf("login_fail:%s:%s", ip, emailAddr)
	count, _ := cache.RDB.Incr(cache.Ctx, failKey).Result()
	if count == 1 {
		cache.RDB.Expire(cache.Ctx, failKey, 15*time.Minute)
	}

	if count >= 5 {
		lockKey := fmt.Sprintf("login_lock:%s:%s", ip, emailAddr)
		cache.RDB.Set(cache.Ctx, lockKey, "1", 15*time.Minute)
	}
}

func ClearLoginFail(ip, emailAddr string) {
	if cache.RDB == nil {
		return
	}
	cache.RDB.Del(cache.Ctx,
		fmt.Sprintf("login_fail:%s:%s", ip, emailAddr),
		fmt.Sprintf("login_lock:%s:%s", ip, emailAddr),
	)
}
