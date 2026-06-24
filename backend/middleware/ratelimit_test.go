package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"xinhang-backend/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TestRateLimit_RedisNil(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	r := gin.New()
	r.POST("/test", RateLimit(1, time.Minute), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("request %d: expected 200 (Redis nil = no limit), got %d", i, w.Code)
		}
	}
}

func TestRateLimit_RedisError(t *testing.T) {
	origRDB := cache.RDB
	defer func() { cache.RDB = origRDB }()

	cache.RDB = redis.NewClient(&redis.Options{Addr: "bad-host:6379"})

	r := gin.New()
	r.POST("/test", RateLimit(1, time.Minute), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200 on Redis error (fail-open), got %d", w.Code)
	}
}

func TestRateLimit_WithRedis(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	pwd := os.Getenv("TEST_REDIS_PASSWORD")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}

	origRDB := cache.RDB
	defer func() { cache.RDB = origRDB }()

	cache.RDB = redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: 15})
	ctx := context.Background()
	cache.RDB.FlushDB(ctx)
	defer cache.RDB.FlushDB(ctx)

	limit := 3
	r := gin.New()
	r.POST("/api/ratelimit-test", RateLimit(limit, time.Minute), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	cache.RDB.FlushDB(ctx)

	for i := 1; i <= limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/ratelimit-test", nil)
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("request %d: expected 200, got %d, body: %s", i, w.Code, w.Body.String())
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/ratelimit-test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 429 {
		keys, _ := cache.RDB.Keys(ctx, "ratelimit:*").Result()
		for _, k := range keys {
			val, _ := cache.RDB.Get(ctx, k).Result()
			t.Logf("  key=%s val=%s", k, val)
		}
		t.Errorf("over-limit request: expected 429, got %d", w.Code)
	}
}
