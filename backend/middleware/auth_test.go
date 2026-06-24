package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
	gin.SetMode(gin.TestMode)
}

const testSecret = "test-jwt-secret-key"

func makeToken(claims jwt.MapClaims) string {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, _ := t.SignedString([]byte(testSecret))
	return s
}

func jsonBody(w *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &m)
	return m
}

func TestJWTAuth_NoHeader(t *testing.T) {
	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_BadFormat(t *testing.T) {
	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "BadFormat")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
	body := jsonBody(w)
	if body["message"] != "认证格式错误" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}

func TestJWTAuth_InvalidToken(t *testing.T) {
	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	token := makeToken(jwt.MapClaims{
		"userId": 1.0,
		"role":   "user",
		"exp":    time.Now().Add(-time.Hour).Unix(),
	})

	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestJWTAuth_MissingClaims(t *testing.T) {
	token := makeToken(jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401 for missing userId, got %d", w.Code)
	}
}

func TestJWTAuth_MissingRole(t *testing.T) {
	token := makeToken(jwt.MapClaims{
		"userId": 1.0,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401 for missing role, got %d", w.Code)
	}
}

func TestJWTAuth_ValidToken(t *testing.T) {
	token := makeToken(jwt.MapClaims{
		"userId": 42.0,
		"role":   "admin",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	var gotID uint
	var gotRole string

	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {
		gotID = c.GetUint("userID")
		if v, ok := c.Get("userRole"); ok {
			gotRole = v.(string)
		}
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if gotID != 42 {
		t.Errorf("expected userID=42, got %d", gotID)
	}
	if gotRole != "admin" {
		t.Errorf("expected role=admin, got %s", gotRole)
	}
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": 1.0, "role": "user",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := tk.SignedString([]byte("different-secret"))

	r := gin.New()
	r.GET("/test", JWTAuth(testSecret), func(c *gin.Context) {})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAdminOnly_NoRole(t *testing.T) {
	r := gin.New()
	r.GET("/test", AdminOnly(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminOnly_UserRole(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("userRole", "user")
		c.Next()
	}, AdminOnly(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminOnly_AdminRole(t *testing.T) {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("userRole", "admin")
		c.Next()
	}, AdminOnly(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
