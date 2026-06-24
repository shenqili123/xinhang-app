package cache

import (
	"context"
	"fmt"
	"log"
	"time"

	"xinhang-backend/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client
var Ctx = context.Background()

var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCR", key)
if current == 1 then
    redis.call("EXPIRE", key, window)
end
return current
`)

func ConnectRedis(cfg *config.Config) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
		PoolSize: 50,
	})

	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		log.Printf("WARNING: Redis not available (%v), rate limiting disabled", err)
		RDB = nil
		return
	}
	log.Printf("Redis connected at %s", cfg.RedisAddr)
}

func CheckRateLimit(key string, maxRequests int, window time.Duration) (bool, error) {
	if RDB == nil {
		return true, nil
	}

	rk := fmt.Sprintf("ratelimit:%s", key)
	windowSec := int(window.Seconds())

	count, err := rateLimitScript.Run(Ctx, RDB, []string{rk}, maxRequests, windowSec).Int64()
	if err != nil {
		return true, err
	}
	return count <= int64(maxRequests), nil
}
