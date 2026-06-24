package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"xinhang-backend/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func setupRedisForSecTest(t *testing.T) func() {
	addr := os.Getenv("TEST_REDIS_ADDR")
	pwd := os.Getenv("TEST_REDIS_PASSWORD")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}
	origRDB := cache.RDB
	cache.RDB = redis.NewClient(&redis.Options{Addr: addr, Password: pwd, DB: 15})
	cache.RDB.FlushDB(context.Background())
	return func() {
		cache.RDB.FlushDB(context.Background())
		cache.RDB = origRDB
	}
}

func TestIPProtection_NilRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	r := gin.New()
	r.Use(IPProtection())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestIPProtection_Normal(t *testing.T) {
	cleanup := setupRedisForSecTest(t)
	defer cleanup()

	r := gin.New()
	r.Use(IPProtection())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("request %d: expected 200, got %d", i, w.Code)
		}
	}
}

func TestCheckLoginLock_NilRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	locked, _ := CheckLoginLock("127.0.0.1", "test@test.com")
	if locked {
		t.Error("expected not locked when Redis nil")
	}
}

func TestRecordLoginFail_NilRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	RecordLoginFail("127.0.0.1", "test@test.com")
}

func TestClearLoginFail_NilRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	ClearLoginFail("127.0.0.1", "test@test.com")
}

func TestLoginLock_Integration(t *testing.T) {
	cleanup := setupRedisForSecTest(t)
	defer cleanup()

	ip := "10.0.0.1"
	email := "lock@test.com"

	for i := 0; i < 5; i++ {
		RecordLoginFail(ip, email)
	}

	locked, msg := CheckLoginLock(ip, email)
	if !locked {
		t.Error("expected locked after 5 failures")
	}
	if msg == "" {
		t.Error("expected lock message")
	}

	ClearLoginFail(ip, email)

	locked2, _ := CheckLoginLock(ip, email)
	if locked2 {
		t.Error("expected unlocked after clear")
	}
}
