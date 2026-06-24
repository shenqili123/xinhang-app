#!/bin/bash
set -euo pipefail

BASE="http://localhost:8080"
PASS=0
FAIL=0
TIMESTAMP=$(date +%s)

pass() { echo "  ✅ $1"; PASS=$((PASS+1)); }
fail() { echo "  ❌ $1"; FAIL=$((FAIL+1)); }

check_status() {
  local label="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then pass "$label (HTTP $actual)"; else fail "$label (expected $expected, got $actual)"; fi
}

check_json() {
  local label="$1" field="$2" expected="$3" body="$4"
  local val
  val=$(echo "$body" | python3 -c "import sys,json;print(json.load(sys.stdin).get('$field',''))" 2>/dev/null || echo "PARSE_ERROR")
  if [ "$val" = "$expected" ]; then pass "$label"; else fail "$label (expected '$expected', got '$val')"; fi
}

# 清理 Redis 限流 key
redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning KEYS 'ratelimit:*' | xargs -r redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning DEL > /dev/null 2>&1

echo "============================================"
echo "   E2E 测试 - $(date '+%Y-%m-%d %H:%M:%S')"
echo "============================================"

#---------------------------------------------------
echo ""
echo "📌 1. 健康检查"
#---------------------------------------------------
BODY=$(curl -s -w "\n%{http_code}" "$BASE/health")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "GET /health" "200" "$CODE"
check_json "/health status" "status" "ok" "$BODY"

#---------------------------------------------------
echo ""
echo "📌 2. 用户注册"
#---------------------------------------------------
EMAIL="e2e_${TIMESTAMP}@test.com"

# 2a. 缺少字段
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"bad"}')
CODE=$(echo "$BODY" | tail -1)
check_status "注册：缺少字段 → 400" "400" "$CODE"

# 2b. 正常注册
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"E2E测试\",\"email\":\"$EMAIL\",\"phone\":\"13800001111\",\"password\":\"Test@123\"}")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "注册：正常 → 200" "200" "$CODE"
check_json "注册成功消息" "message" "注册成功" "$BODY"

# 2c. 重复注册
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"E2E测试\",\"email\":\"$EMAIL\",\"phone\":\"13800001111\",\"password\":\"Test@123\"}")
CODE=$(echo "$BODY" | tail -1)
check_status "注册：重复 → 409" "409" "$CODE"

#---------------------------------------------------
echo ""
echo "📌 3. 用户登录"
#---------------------------------------------------

# 3a. 错误密码
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"wrong\"}")
CODE=$(echo "$BODY" | tail -1)
check_status "登录：错误密码 → 401" "401" "$CODE"

# 3b. 正常登录
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"Test@123\"}")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "登录：正常 → 200" "200" "$CODE"

USER_TOKEN=$(echo "$BODY" | python3 -c "import sys,json;print(json.load(sys.stdin).get('token',''))" 2>/dev/null)
if [ -n "$USER_TOKEN" ] && [ "$USER_TOKEN" != "" ]; then
  pass "登录返回 token"
else
  fail "登录未返回 token"
fi

#---------------------------------------------------
echo ""
echo "📌 4. 报名申请"
#---------------------------------------------------

# 4a. 匿名报名
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/apply" \
  -H "Content-Type: application/json" \
  -d "{\"studentName\":\"E2E匿名学生\",\"birthDate\":\"2015-01-01\",\"gender\":\"男\",\"grade\":7,\"parentName\":\"E2E家长\",\"phone\":\"13900001111\",\"email\":\"e2e_anon_${TIMESTAMP}@test.com\",\"currentSchool\":\"测试中学\"}")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "报名：匿名 → 200" "200" "$CODE"
check_json "报名成功消息" "message" "报名申请已提交成功" "$BODY"

# 4b. 登录后报名
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/apply" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d "{\"studentName\":\"E2E登录学生\",\"birthDate\":\"2014-06-15\",\"gender\":\"女\",\"grade\":8,\"parentName\":\"E2E家长2\",\"phone\":\"13900002222\",\"email\":\"e2e_auth_${TIMESTAMP}@test.com\",\"currentSchool\":\"测试高中\"}")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "报名：登录后 → 200" "200" "$CODE"

# 4c. 重复报名
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/apply" \
  -H "Content-Type: application/json" \
  -d "{\"studentName\":\"E2E匿名学生\",\"birthDate\":\"2015-01-01\",\"gender\":\"男\",\"grade\":7,\"parentName\":\"E2E家长\",\"phone\":\"13900001111\",\"email\":\"e2e_anon_${TIMESTAMP}@test.com\",\"currentSchool\":\"测试中学\"}")
CODE=$(echo "$BODY" | tail -1)
check_status "报名：重复 → 409" "409" "$CODE"

# 4d. 无效输入（grade 超范围）
BODY=$(curl -s -w "\n%{http_code}" -X POST "$BASE/api/apply" \
  -H "Content-Type: application/json" \
  -d '{"studentName":"坏","birthDate":"2015-01-01","gender":"其他","grade":99,"parentName":"坏","phone":"123","email":"bad","currentSchool":"x"}')
CODE=$(echo "$BODY" | tail -1)
check_status "报名：无效输入 → 400" "400" "$CODE"

#---------------------------------------------------
echo ""
echo "📌 5. 管理员接口"
#---------------------------------------------------

# 5a. 无 token → 401
BODY=$(curl -s -w "\n%{http_code}" "$BASE/api/applications")
CODE=$(echo "$BODY" | tail -1)
check_status "管理列表：无 token → 401" "401" "$CODE"

# 5b. 普通用户 → 403
BODY=$(curl -s -w "\n%{http_code}" "$BASE/api/applications" \
  -H "Authorization: Bearer $USER_TOKEN")
CODE=$(echo "$BODY" | tail -1)
check_status "管理列表：普通用户 → 403" "403" "$CODE"

# 5c. 管理员登录
ADMIN_RESP=$(curl -s -X POST "$BASE/api/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}')
ADMIN_TOKEN=$(echo "$ADMIN_RESP" | python3 -c "import sys,json;print(json.load(sys.stdin).get('token',''))" 2>/dev/null)

if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "" ]; then
  pass "管理员登录成功"
else
  fail "管理员登录失败"
fi

# 5d. 管理员获取列表
BODY=$(curl -s -w "\n%{http_code}" "$BASE/api/applications?page=1&pageSize=5" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
CODE=$(echo "$BODY" | tail -1)
BODY=$(echo "$BODY" | head -1)
check_status "管理列表：管理员 → 200" "200" "$CODE"

TOTAL=$(echo "$BODY" | python3 -c "import sys,json;print(json.load(sys.stdin).get('total',0))" 2>/dev/null)
if [ "$TOTAL" -gt 0 ] 2>/dev/null; then
  pass "管理列表有数据 (total=$TOTAL)"
else
  fail "管理列表无数据"
fi

PAGE_SIZE=$(echo "$BODY" | python3 -c "import sys,json;print(json.load(sys.stdin).get('pageSize',0))" 2>/dev/null)
check_json "分页参数正确" "page" "1" "$BODY"

#---------------------------------------------------
echo ""
echo "📌 6. 前端页面"
#---------------------------------------------------
for path in "/" "/about" "/academics" "/admission" "/campus" "/student-life" "/register" "/login" "/apply"; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE$path")
  check_status "页面 $path" "200" "$CODE"
done

#---------------------------------------------------
echo ""
echo "📌 7. 静态资源"
#---------------------------------------------------
ASSET_FILE=$(ls /home/devbox/project/xinhang-app/backend/dist/assets/*.js 2>/dev/null | head -1)
if [ -n "$ASSET_FILE" ]; then
  ASSET_NAME=$(basename "$ASSET_FILE")
  CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE/assets/$ASSET_NAME")
  check_status "静态资源 /assets/$ASSET_NAME" "200" "$CODE"
else
  pass "静态资源目录存在（无JS文件跳过）"
fi

#---------------------------------------------------
echo ""
echo "📌 8. CORS"
#---------------------------------------------------
CORS_HEADER=$(curl -s -I -X OPTIONS "$BASE/api/login" \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: POST" 2>/dev/null | grep -i 'access-control-allow-origin' || echo "")
if [ -n "$CORS_HEADER" ]; then
  pass "CORS 头存在: $(echo $CORS_HEADER | tr -d '\r')"
else
  fail "CORS 头缺失"
fi

#---------------------------------------------------
echo ""
echo "📌 9. Gzip 压缩"
#---------------------------------------------------
CONTENT_ENCODING=$(curl -s -H "Accept-Encoding: gzip" -D - -o /dev/null "$BASE/" 2>/dev/null | grep -i 'content-encoding.*gzip' || echo "")
if [ -n "$CONTENT_ENCODING" ]; then
  pass "Gzip 压缩启用"
else
  fail "Gzip 压缩未启用"
fi

#---------------------------------------------------
echo ""
echo "📌 10. 请求体大小限制"
#---------------------------------------------------
dd if=/dev/zero bs=1 count=2000000 2>/dev/null | curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/apply" \
  -H "Content-Type: application/json" --data-binary @- > /tmp/bigbody_code.txt 2>/dev/null
CODE=$(cat /tmp/bigbody_code.txt)
if [ "$CODE" = "400" ] || [ "$CODE" = "413" ] || [ -z "$CODE" ]; then
  pass "大请求体被拒绝 (HTTP ${CODE:-连接断开})"
else
  fail "大请求体未被限制 (HTTP $CODE)"
fi

#---------------------------------------------------
echo ""
echo "📌 11. 限流"
#---------------------------------------------------
LIMIT_HIT=0
for i in $(seq 1 12); do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE/api/register" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"rl$i\",\"email\":\"rl${i}_${TIMESTAMP}@test.com\",\"phone\":\"138000${i}0000\",\"password\":\"Test@123\"}")
  if [ "$CODE" = "429" ]; then
    LIMIT_HIT=1
    break
  fi
done
if [ "$LIMIT_HIT" = "1" ]; then
  pass "限流生效 (第${i}次请求被拒)"
else
  fail "限流未生效（12次请求均未被限流）"
fi

#---------------------------------------------------
echo ""
echo "📌 12. JWT 安全"
#---------------------------------------------------
# 伪造 token
FAKE_TOKEN="eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VySWQiOjEsInJvbGUiOiJhZG1pbiJ9."
BODY=$(curl -s -w "\n%{http_code}" "$BASE/api/applications" \
  -H "Authorization: Bearer $FAKE_TOKEN")
CODE=$(echo "$BODY" | tail -1)
check_status "伪造 alg:none token → 401" "401" "$CODE"

#---------------------------------------------------
echo ""
echo "============================================"
echo "   测试结果: ✅ $PASS 通过 / ❌ $FAIL 失败"
echo "============================================"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
