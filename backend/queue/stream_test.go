package queue

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"xinhang-backend/cache"
	"xinhang-backend/models"

	"github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) {
	addr := os.Getenv("TEST_REDIS_ADDR")
	pwd := os.Getenv("TEST_REDIS_PASSWORD")
	if addr == "" {
		t.Skip("TEST_REDIS_ADDR not set")
	}
	cache.RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       15,
	})
	cache.RDB.FlushDB(context.Background())
}

func cleanupTestRedis() {
	if cache.RDB != nil {
		cache.RDB.FlushDB(context.Background())
	}
}

func TestIsEnabled_Default(t *testing.T) {
	enabled = false
	if IsEnabled() {
		t.Error("expected IsEnabled=false by default")
	}
}

func TestConnectQueue_NilRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	ConnectQueue()
	if IsEnabled() {
		t.Error("expected queue disabled when Redis is nil")
	}
}

func TestPublishApplication_Disabled(t *testing.T) {
	enabled = false
	err := PublishApplication(&models.Application{Email: "test@test.com"})
	if err != nil {
		t.Errorf("expected nil error when disabled, got %v", err)
	}
}

func TestStartConsumer_Disabled(t *testing.T) {
	enabled = false
	called := false
	StartConsumer(func(app *models.Application) error {
		called = true
		return nil
	})
	if called {
		t.Error("handler should not be called when disabled")
	}
}

func TestClose(t *testing.T) {
	Close()
}

func TestConnectQueue_Success(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		cache.RDB = origRDB
		enabled = origEnabled
		cleanupTestRedis()
	}()

	setupTestRedis(t)
	ConnectQueue()
	if !IsEnabled() {
		t.Error("expected queue to be enabled")
	}
}

func TestConnectQueue_Idempotent(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		cache.RDB = origRDB
		enabled = origEnabled
		cleanupTestRedis()
	}()

	setupTestRedis(t)
	ConnectQueue()
	ConnectQueue()
	if !IsEnabled() {
		t.Error("expected queue to be enabled after second connect")
	}
}

func TestPublishApplication_Integration(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		cache.RDB = origRDB
		enabled = origEnabled
		cleanupTestRedis()
	}()

	setupTestRedis(t)
	ConnectQueue()

	app := &models.Application{
		StudentName: "测试", BirthDate: "2015-01-01", Gender: "男",
		Grade: 7, ParentName: "家长", Phone: "139",
		Email: "pub@test.com",
	}
	err := PublishApplication(app)
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}

	length, _ := cache.RDB.XLen(context.Background(), streamName).Result()
	if length != 1 {
		t.Errorf("expected stream length=1, got %d", length)
	}
}

func TestConnectQueue_ErrorPath(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		cache.RDB = origRDB
		enabled = origEnabled
	}()

	cache.RDB = redis.NewClient(&redis.Options{Addr: "bad-host:6379"})
	ConnectQueue()
	if IsEnabled() {
		t.Error("expected queue disabled on bad Redis")
	}
}

func TestPublishApplication_Integration_Marshal(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		cache.RDB = origRDB
		enabled = origEnabled
		cleanupTestRedis()
	}()

	setupTestRedis(t)
	ConnectQueue()

	app := &models.Application{
		StudentName: "编码测试", Email: "marshal@test.com",
		BirthDate: "2015-01-01", Gender: "男", Grade: 7,
		ParentName: "家长", Phone: "139",
	}
	err := PublishApplication(app)
	if err != nil {
		t.Fatalf("publish failed: %v", err)
	}
}

func TestStartConsumer_BadData(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		Close()
		time.Sleep(200 * time.Millisecond)
		cleanupTestRedis()
		cache.RDB = origRDB
		enabled = origEnabled
	}()

	setupTestRedis(t)
	ConnectQueue()

	cache.RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"email": "bad@test.com",
			"data":  "not valid json{{{",
		},
	})

	cache.RDB.XAdd(context.Background(), &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"email": "no-data@test.com",
			"other": "no data field",
		},
	})

	handled := make(chan bool, 1)
	StartConsumer(func(a *models.Application) error {
		handled <- true
		return nil
	})

	select {
	case <-handled:
		t.Error("handler should not be called for bad data")
	case <-time.After(8 * time.Second):
		t.Log("correctly skipped bad messages")
	}
}

func TestStartConsumer_HandlerError(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		Close()
		time.Sleep(200 * time.Millisecond)
		cleanupTestRedis()
		cache.RDB = origRDB
		enabled = origEnabled
	}()

	setupTestRedis(t)
	ConnectQueue()

	app := &models.Application{
		StudentName: "错误处理", Email: "err@test.com",
		BirthDate: "2015-01-01", Gender: "男", Grade: 7,
		ParentName: "家长", Phone: "139",
	}
	PublishApplication(app)

	received := make(chan bool, 1)
	StartConsumer(func(a *models.Application) error {
		received <- true
		return fmt.Errorf("simulated handler error")
	})

	select {
	case <-received:
		t.Log("handler was called and returned error (message not ACKed)")
	case <-time.After(10 * time.Second):
		t.Error("consumer did not process message within 10s")
	}
}

func TestStartConsumer_Integration(t *testing.T) {
	origRDB := cache.RDB
	origEnabled := enabled
	defer func() {
		Close()
		time.Sleep(200 * time.Millisecond)
		cleanupTestRedis()
		cache.RDB = origRDB
		enabled = origEnabled
	}()

	setupTestRedis(t)
	ConnectQueue()

	app := &models.Application{
		StudentName: "消费者测试", BirthDate: "2015-01-01", Gender: "女",
		Grade: 8, ParentName: "家长B", Phone: "138",
		Email: "consumer@test.com",
	}
	PublishApplication(app)

	received := make(chan *models.Application, 1)
	StartConsumer(func(a *models.Application) error {
		received <- a
		return nil
	})

	select {
	case got := <-received:
		if got.StudentName != "消费者测试" {
			t.Errorf("expected 消费者测试, got %s", got.StudentName)
		}
		if got.Email != "consumer@test.com" {
			t.Errorf("expected consumer@test.com, got %s", got.Email)
		}
	case <-time.After(10 * time.Second):
		t.Error("consumer did not receive message within 10s")
	}
}
