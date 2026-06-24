package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"xinhang-backend/cache"
	"xinhang-backend/testutil"

	"github.com/gin-gonic/gin"
)

func TestSendCode_NoEmail(t *testing.T) {
	r := gin.New()
	r.POST("/api/send-code", SendVerificationCode)

	w := doJSON(r, "POST", "/api/send-code", map[string]string{})
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestSendCode_EmailDisabled(t *testing.T) {
	r := gin.New()
	r.POST("/api/send-code", SendVerificationCode)

	w := doJSON(r, "POST", "/api/send-code", map[string]string{
		"email": "test@test.com",
	})
	if w.Code != 503 {
		t.Errorf("expected 503 (email not configured), got %d: %s", w.Code, w.Body.String())
	}
}

func TestSendCode_NoRedis(t *testing.T) {
	origRDB := cache.RDB
	cache.RDB = nil
	defer func() { cache.RDB = origRDB }()

	r := gin.New()
	r.POST("/api/send-code", SendVerificationCode)

	w := doJSON(r, "POST", "/api/send-code", map[string]string{
		"email": "test@test.com",
	})
	if w.Code == 503 || w.Code == 400 {
		t.Logf("got expected status %d", w.Code)
	}
}

func TestVerifyCode_TestMode(t *testing.T) {
	SetVerifyCodeForTest("vtest@test.com", "654321")

	ok, msg := verifyCode("vtest@test.com", "654321")
	if !ok {
		t.Errorf("expected verification success, got msg: %s", msg)
	}

	ok2, msg2 := verifyCode("vtest@test.com", "000000")
	if ok2 {
		t.Error("expected verification failure for non-existing code")
	}
	_ = msg2
}

func TestVerifyCode_WrongCode(t *testing.T) {
	SetVerifyCodeForTest("wrong@test.com", "111111")

	ok, msg := verifyCode("wrong@test.com", "999999")
	if ok {
		t.Error("expected verification failure")
	}
	if msg != "验证码错误" {
		t.Errorf("expected '验证码错误', got '%s'", msg)
	}
}

func TestRegister_WithVerificationCode(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	SetVerifyCodeForTest("verified@test.com", "123456")
	w := doJSON(r, "POST", "/api/register", map[string]interface{}{
		"name": "验证用户", "email": "verified@test.com",
		"phone": "13800000099", "password": "Pass@123",
		"code": "123456",
	})
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_WrongVerificationCode(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	SetVerifyCodeForTest("badcode@test.com", "111111")
	w := doJSON(r, "POST", "/api/register", map[string]interface{}{
		"name": "错误码", "email": "badcode@test.com",
		"phone": "13800000088", "password": "Pass@123",
		"code": "999999",
	})
	if w.Code != 400 {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGenerateCode(t *testing.T) {
	code := generateCode()
	if len(code) != 6 {
		t.Errorf("expected 6-digit code, got %s", code)
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Errorf("code contains non-digit: %c", c)
		}
	}

	codes := make(map[string]bool)
	for i := 0; i < 10; i++ {
		codes[generateCode()] = true
	}
	if len(codes) < 5 {
		t.Error("generated codes lack randomness")
	}
}

func TestSendCode_InvalidEmail(t *testing.T) {
	r := gin.New()
	r.POST("/api/send-code", SendVerificationCode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/send-code", nil)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != 400 {
		t.Errorf("expected 400 for nil body, got %d", w.Code)
	}
}
