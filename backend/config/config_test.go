package config

import (
	"os"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	os.Clearenv()
	cfg := Load()

	if cfg.DBHost != "localhost" {
		t.Errorf("expected DBHost=localhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "5432" {
		t.Errorf("expected DBPort=5432, got %s", cfg.DBPort)
	}
	if cfg.DBUser != "postgres" {
		t.Errorf("expected DBUser=postgres, got %s", cfg.DBUser)
	}
	if cfg.DBName != "xinhang" {
		t.Errorf("expected DBName=xinhang, got %s", cfg.DBName)
	}
	if cfg.DBMaxConns != 100 {
		t.Errorf("expected DBMaxConns=100, got %d", cfg.DBMaxConns)
	}
	if cfg.DBIdleConns != 20 {
		t.Errorf("expected DBIdleConns=20, got %d", cfg.DBIdleConns)
	}
	if cfg.Port != "8080" {
		t.Errorf("expected Port=8080, got %s", cfg.Port)
	}
	if cfg.JWTSecret != "xinhang-secret-change-me" {
		t.Errorf("expected default JWTSecret, got %s", cfg.JWTSecret)
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("expected RedisAddr=localhost:6379, got %s", cfg.RedisAddr)
	}
	if cfg.RedisDB != 0 {
		t.Errorf("expected RedisDB=0, got %d", cfg.RedisDB)
	}
}

func TestLoadFromEnv(t *testing.T) {
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "9999")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpwd")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_MAX_CONNS", "50")
	os.Setenv("DB_IDLE_CONNS", "10")
	os.Setenv("JWT_SECRET", "my-secret")
	os.Setenv("PORT", "3000")
	os.Setenv("REDIS_ADDR", "redis:6380")
	os.Setenv("REDIS_PASSWORD", "redisPwd")
	os.Setenv("REDIS_DB", "2")
	defer os.Clearenv()

	cfg := Load()

	if cfg.DBHost != "testhost" {
		t.Errorf("expected testhost, got %s", cfg.DBHost)
	}
	if cfg.DBPort != "9999" {
		t.Errorf("expected 9999, got %s", cfg.DBPort)
	}
	if cfg.DBMaxConns != 50 {
		t.Errorf("expected 50, got %d", cfg.DBMaxConns)
	}
	if cfg.DBIdleConns != 10 {
		t.Errorf("expected 10, got %d", cfg.DBIdleConns)
	}
	if cfg.JWTSecret != "my-secret" {
		t.Errorf("expected my-secret, got %s", cfg.JWTSecret)
	}
	if cfg.Port != "3000" {
		t.Errorf("expected 3000, got %s", cfg.Port)
	}
	if cfg.RedisAddr != "redis:6380" {
		t.Errorf("expected redis:6380, got %s", cfg.RedisAddr)
	}
	if cfg.RedisPassword != "redisPwd" {
		t.Errorf("expected redisPwd, got %s", cfg.RedisPassword)
	}
	if cfg.RedisDB != 2 {
		t.Errorf("expected 2, got %d", cfg.RedisDB)
	}
}

func TestGetEnvIntInvalid(t *testing.T) {
	os.Setenv("BAD_INT", "notanumber")
	defer os.Unsetenv("BAD_INT")

	val := getEnvInt("BAD_INT", 42)
	if val != 42 {
		t.Errorf("expected fallback 42 for invalid int, got %d", val)
	}
}

func TestDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "myhost",
		DBPort:     "5432",
		DBUser:     "myuser",
		DBPassword: "mypass",
		DBName:     "mydb",
	}
	dsn := cfg.DSN()
	expected := "host=myhost port=5432 user=myuser password=mypass dbname=mydb sslmode=disable TimeZone=Asia/Shanghai"
	if dsn != expected {
		t.Errorf("DSN mismatch:\n  got:      %s\n  expected: %s", dsn, expected)
	}
}
