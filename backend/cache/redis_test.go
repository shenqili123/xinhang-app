package cache

import (
	"os"
	"testing"
	"time"

	"xinhang-backend/config"

	"github.com/redis/go-redis/v9"
)

func getTestRedisConfig() *config.Config {
	return &config.Config{
		RedisAddr:     os.Getenv("TEST_REDIS_ADDR"),
		RedisPassword: os.Getenv("TEST_REDIS_PASSWORD"),
		RedisDB:       15,
	}
}

func skipIfNoRedis(t *testing.T) *config.Config {
	cfg := getTestRedisConfig()
	if cfg.RedisAddr == "" {
		t.Skip("TEST_REDIS_ADDR not set, skipping Redis integration test")
	}
	return cfg
}

func TestCheckRateLimit_NilRedis(t *testing.T) {
	origRDB := RDB
	RDB = nil
	defer func() { RDB = origRDB }()

	allowed, err := CheckRateLimit("test-key", 5, time.Minute)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if !allowed {
		t.Error("expected allowed=true when Redis is nil")
	}
}

func TestConnectRedis_Success(t *testing.T) {
	cfg := skipIfNoRedis(t)
	origRDB := RDB
	defer func() { RDB = origRDB }()

	ConnectRedis(cfg)
	if RDB == nil {
		t.Fatal("expected RDB to be set after ConnectRedis")
	}
	RDB.FlushDB(Ctx)
}

func TestConnectRedis_BadAddr(t *testing.T) {
	origRDB := RDB
	defer func() { RDB = origRDB }()

	cfg := &config.Config{
		RedisAddr:     "bad-host-does-not-exist:9999",
		RedisPassword: "",
		RedisDB:       15,
	}
	ConnectRedis(cfg)
	if RDB != nil {
		t.Error("expected RDB=nil after failed connect")
	}
}

func TestCheckRateLimit_Integration(t *testing.T) {
	cfg := skipIfNoRedis(t)
	origRDB := RDB
	defer func() { RDB = origRDB }()

	ConnectRedis(cfg)
	if RDB == nil {
		t.Fatal("Redis not connected")
	}
	RDB.FlushDB(Ctx)

	limit := 3
	key := "test:ratelimit:integration"

	for i := 1; i <= limit; i++ {
		allowed, err := CheckRateLimit(key, limit, 10*time.Second)
		if err != nil {
			t.Fatalf("request %d: unexpected error: %v", i, err)
		}
		if !allowed {
			t.Fatalf("request %d: should be allowed (within limit %d)", i, limit)
		}
	}

	allowed, err := CheckRateLimit(key, limit, 10*time.Second)
	if err != nil {
		t.Fatalf("over-limit: unexpected error: %v", err)
	}
	if allowed {
		t.Error("request should be rejected (over limit)")
	}

	RDB.FlushDB(Ctx)
}

func TestCheckRateLimit_RedisError(t *testing.T) {
	origRDB := RDB
	defer func() { RDB = origRDB }()

	RDB = redis.NewClient(&redis.Options{
		Addr: "bad-host-does-not-exist:6379",
	})

	allowed, err := CheckRateLimit("err-key", 5, time.Minute)
	if err == nil {
		t.Log("no error (DNS might have resolved unexpectedly)")
	}
	if !allowed {
		t.Error("expected allowed=true on Redis error (fail-open)")
	}
}
