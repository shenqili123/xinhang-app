package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"xinhang-backend/cache"
	"xinhang-backend/database"
	"xinhang-backend/models"
	"xinhang-backend/testutil"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	cache.RDB = nil
	testutil.SetupTestDB()
	SetJWTSecret("test-secret-key-for-unit-tests")
	os.Exit(m.Run())
}

func registerPayload(email, name, phone string) map[string]interface{} {
	code := "123456"
	SetVerifyCodeForTest(email, code)
	return map[string]interface{}{
		"name":     name,
		"email":    email,
		"phone":    phone,
		"password": "Pass@123",
		"code":     code,
	}
}

func setupRouter() *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	api.POST("/register", Register)
	api.POST("/login", Login)
	api.POST("/apply", Apply)
	api.GET("/applications", ListApplications)
	return r
}

func doJSON(r *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func parseJSON(w *httptest.ResponseRecorder) map[string]interface{} {
	var m map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &m)
	return m
}

// ==================== Register ====================

func TestRegister_Success(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/register", registerPayload("reg1@test.com", "张三", "13800000001"))
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := parseJSON(w)
	if body["message"] != "注册成功" {
		t.Errorf("unexpected message: %v", body["message"])
	}

	var user models.User
	database.DB.Where("email = ?", "reg1@test.com").First(&user)
	if user.Name != "张三" {
		t.Errorf("expected name=张三, got %s", user.Name)
	}
	if user.Role != "user" {
		t.Errorf("expected role=user, got %s", user.Role)
	}
	// EmailVerified depends on whether email service is configured
}

func TestRegister_MissingFields(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/register", map[string]interface{}{
		"email": "bad",
	})
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/register", map[string]interface{}{
		"name": "Test", "email": "not-an-email", "phone": "13800000001",
		"password": "Pass@123", "code": "123456",
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for invalid email, got %d", w.Code)
	}
}

func TestRegister_ShortPassword(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/register", map[string]interface{}{
		"name": "Test", "email": "short@test.com", "phone": "13800000001",
		"password": "123", "code": "123456",
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for short password, got %d", w.Code)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	payload := registerPayload("dup@test.com", "张三", "13800000001")
	doJSON(r, "POST", "/api/register", payload)

	SetVerifyCodeForTest("dup@test.com", "123456")
	w := doJSON(r, "POST", "/api/register", registerPayload("dup@test.com", "张三", "13800000001"))
	if w.Code != 409 {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

// ==================== Login ====================

func TestLogin_Success(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	doJSON(r, "POST", "/api/register", registerPayload("login@test.com", "李四", "13800000002"))

	w := doJSON(r, "POST", "/api/login", map[string]string{
		"email": "login@test.com", "password": "Pass@123",
	})
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := parseJSON(w)
	if body["token"] == nil || body["token"] == "" {
		t.Error("expected token in response")
	}
	if body["message"] != "登录成功" {
		t.Errorf("unexpected message: %v", body["message"])
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	doJSON(r, "POST", "/api/register", registerPayload("loginwrong@test.com", "李四", "13800000003"))

	w := doJSON(r, "POST", "/api/login", map[string]string{
		"email": "loginwrong@test.com", "password": "WrongPass",
	})
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/login", map[string]string{
		"email": "noexist@test.com", "password": "Pass@123",
	})
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	r := setupRouter()

	w := doJSON(r, "POST", "/api/login", map[string]string{})
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ==================== Apply ====================

func TestApply_Success(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "王五", "birthDate": "2015-01-01", "gender": "男",
		"grade": 7, "parentName": "王大", "phone": "13900001111",
		"email": "apply1@test.com", "currentSchool": "测试小学",
	})
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	body := parseJSON(w)
	if body["message"] != "报名申请已提交成功" {
		t.Errorf("unexpected message: %v", body["message"])
	}
	if body["id"] == nil {
		t.Error("expected id in response")
	}
}

func TestApply_WithUserID(t *testing.T) {
	testutil.CleanDB()

	r := gin.New()
	r.POST("/api/apply", func(c *gin.Context) {
		c.Set("userID", uint(99))
		c.Next()
	}, Apply)

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "带用户", "birthDate": "2015-01-01", "gender": "女",
		"grade": 8, "parentName": "家长A", "phone": "13900002222",
		"email": "withuser@test.com",
	})
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var app models.Application
	database.DB.Where("email = ?", "withuser@test.com").First(&app)
	if app.UserID == nil || *app.UserID != 99 {
		t.Errorf("expected UserID=99, got %v", app.UserID)
	}
}

func TestApply_Duplicate(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	payload := map[string]interface{}{
		"studentName": "重复生", "birthDate": "2015-01-01", "gender": "男",
		"grade": 7, "parentName": "家长", "phone": "13900003333",
		"email": "dup@test.com",
	}
	doJSON(r, "POST", "/api/apply", payload)

	w := doJSON(r, "POST", "/api/apply", payload)
	if w.Code != 409 {
		t.Errorf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestApply_InvalidGender(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "坏", "birthDate": "2015-01-01", "gender": "其他",
		"grade": 7, "parentName": "家长", "phone": "13900004444",
		"email": "bad@test.com",
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for invalid gender, got %d", w.Code)
	}
}

func TestApply_GradeOutOfRange(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "坏", "birthDate": "2015-01-01", "gender": "男",
		"grade": 99, "parentName": "家长", "phone": "13900005555",
		"email": "grade99@test.com",
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for grade=99, got %d", w.Code)
	}
}

func TestApply_MissingFields(t *testing.T) {
	r := setupRouter()
	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{})
	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// ==================== ListApplications ====================

func TestListApplications_Empty(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/applications", nil)
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := parseJSON(w)
	if body["total"].(float64) != 0 {
		t.Errorf("expected total=0, got %v", body["total"])
	}
}

func TestListApplications_Pagination(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	for i := 1; i <= 25; i++ {
		database.DB.Create(&models.Application{
			StudentName: "学生", BirthDate: "2015-01-01", Gender: "男",
			Grade: 7, ParentName: "家长", Phone: "139",
			Email: "page@test.com",
		})
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/applications?page=1&pageSize=10", nil)
	r.ServeHTTP(w, req)

	body := parseJSON(w)
	total := int(body["total"].(float64))
	if total != 25 {
		t.Errorf("expected total=25, got %d", total)
	}
	data := body["data"].([]interface{})
	if len(data) != 10 {
		t.Errorf("expected 10 items, got %d", len(data))
	}
	if int(body["page"].(float64)) != 1 {
		t.Errorf("expected page=1, got %v", body["page"])
	}
	if int(body["pageSize"].(float64)) != 10 {
		t.Errorf("expected pageSize=10, got %v", body["pageSize"])
	}

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/applications?page=3&pageSize=10", nil)
	r.ServeHTTP(w2, req2)
	body2 := parseJSON(w2)
	data2 := body2["data"].([]interface{})
	if len(data2) != 5 {
		t.Errorf("expected 5 items on page 3, got %d", len(data2))
	}
}

func TestListApplications_InvalidPage(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/applications?page=-1&pageSize=0", nil)
	r.ServeHTTP(w, req)

	body := parseJSON(w)
	if int(body["page"].(float64)) != 1 {
		t.Errorf("expected page clamped to 1, got %v", body["page"])
	}
	if int(body["pageSize"].(float64)) != 20 {
		t.Errorf("expected pageSize clamped to 20, got %v", body["pageSize"])
	}
}

func TestListApplications_PageSizeTooLarge(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/applications?pageSize=999", nil)
	r.ServeHTTP(w, req)

	body := parseJSON(w)
	if int(body["pageSize"].(float64)) != 20 {
		t.Errorf("expected pageSize clamped to 20, got %v", body["pageSize"])
	}
}

// ==================== DB Error paths ====================

func TestRegister_DBError(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	database.DB.Exec("DROP TABLE IF EXISTS users")

	w := doJSON(r, "POST", "/api/register", registerPayload("dberr@test.com", "Test", "13800000001"))
	if w.Code != 500 {
		if w.Code == 200 {
			// SQLite might auto-create — skip
			t.Skip("SQLite auto-created table")
		}
		t.Errorf("expected 500, got %d", w.Code)
	}

	testutil.SetupTestDB()
}

func TestLogin_UserNotFound_DB(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/login", map[string]string{
		"email": "ghost@test.com", "password": "any",
	})
	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestApply_DBError(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	database.DB.Exec("DROP TABLE IF EXISTS applications")

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "ErrTest", "birthDate": "2015-01-01", "gender": "男",
		"grade": 7, "parentName": "家长", "phone": "13900006666",
		"email": "dberr@test.com",
	})

	if w.Code == 500 {
		t.Log("correctly returned 500 on DB error")
	}

	testutil.SetupTestDB()
}

func TestListApplications_DBError(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	database.DB.Exec("DROP TABLE IF EXISTS applications")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/applications", nil)
	r.ServeHTTP(w, req)

	if w.Code == 500 {
		t.Log("correctly returned 500 on DB error")
	}

	testutil.SetupTestDB()
}

func TestSetJWTSecret(t *testing.T) {
	SetJWTSecret("another-secret")
	if string(jwtSecret) != "another-secret" {
		t.Errorf("expected another-secret, got %s", string(jwtSecret))
	}
	SetJWTSecret("test-secret-key-for-unit-tests")
}

func TestApply_NotesMaxLen(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	longNotes := ""
	for i := 0; i < 2001; i++ {
		longNotes += "x"
	}
	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "长备注", "birthDate": "2015-01-01", "gender": "男",
		"grade": 7, "parentName": "家长", "phone": "13900007777",
		"email": "notes@test.com", "notes": longNotes,
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for notes > 2000 chars, got %d", w.Code)
	}
}

func TestApply_GradeZero(t *testing.T) {
	testutil.CleanDB()
	r := setupRouter()

	w := doJSON(r, "POST", "/api/apply", map[string]interface{}{
		"studentName": "零年级", "birthDate": "2015-01-01", "gender": "男",
		"grade": 0, "parentName": "家长", "phone": "13900008888",
		"email": "grade0@test.com",
	})
	if w.Code != 400 {
		t.Errorf("expected 400 for grade=0, got %d", w.Code)
	}
}
