#!/bin/bash
set -euo pipefail

BASE="http://localhost:8080"
DURATION=10
CONCURRENCY=100
TIMESTAMP=$(date +%s)

echo "============================================"
echo "   压力测试 - $(date '+%Y-%m-%d %H:%M:%S')"
echo "   并发: $CONCURRENCY | 持续: ${DURATION}s"
echo "============================================"

# 清除限流
redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning KEYS 'ratelimit:*' | xargs -r redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning DEL > /dev/null 2>&1

# 检查 wrk 是否可用
if ! command -v wrk &>/dev/null; then
  echo "安装 wrk..."
  apt-get update -qq && apt-get install -y -qq wrk > /dev/null 2>&1 || {
    echo "wrk 安装失败，使用 curl 并发替代"
    USE_CURL=1
  }
fi
USE_CURL=${USE_CURL:-0}

run_wrk() {
  local label="$1" method="$2" url="$3" script="${4:-}"
  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "📊 $label"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  if [ "$USE_CURL" = "1" ]; then
    local start end count ok fail_count
    start=$(date +%s%N)
    count=0; ok=0; fail_count=0
    local end_time=$(($(date +%s) + DURATION))
    while [ "$(date +%s)" -lt "$end_time" ]; do
      for j in $(seq 1 $CONCURRENCY); do
        if [ "$method" = "GET" ]; then
          CODE=$(curl -s -o /dev/null -w "%{http_code}" "$url" 2>/dev/null) &
        fi
        count=$((count+1))
      done
      wait
    done
    end=$(date +%s%N)
    local elapsed=$(( (end - start) / 1000000 ))
    echo "  请求总数: ~$count | 耗时: ${elapsed}ms"
    return
  fi

  if [ -n "$script" ]; then
    wrk -t4 -c${CONCURRENCY} -d${DURATION}s --latency -s "$script" "$url" 2>&1 | tail -20
  else
    wrk -t4 -c${CONCURRENCY} -d${DURATION}s --latency "$url" 2>&1 | tail -20
  fi
}

#---------------------------------------------------
echo ""
echo "📌 1. 静态页面（首页）"
#---------------------------------------------------
run_wrk "GET /" "GET" "$BASE/"

#---------------------------------------------------
echo ""
echo "📌 2. 健康检查"
#---------------------------------------------------
run_wrk "GET /health" "GET" "$BASE/health"

#---------------------------------------------------
echo ""
echo "📌 3. API - 注册（wrk Lua）"
#---------------------------------------------------
cat > /tmp/wrk_register.lua << 'LUAEOF'
local counter = 0
local thread_id = 0

function setup(thread)
  thread:set("id", thread_id)
  thread_id = thread_id + 1
end

function init(args)
  math.randomseed(os.time() + id * 1000)
end

request = function()
  counter = counter + 1
  local body = string.format(
    '{"name":"stress_%d_%d","email":"s%d_%d_%d@load.test","phone":"138%08d","password":"Test@123"}',
    id, counter, os.time(), id, counter, math.random(99999999))
  return wrk.format("POST", "/api/register", {
    ["Content-Type"] = "application/json"
  }, body)
end
LUAEOF
run_wrk "POST /api/register" "POST" "$BASE" "/tmp/wrk_register.lua"

#---------------------------------------------------
echo ""
echo "📌 4. API - 登录"
#---------------------------------------------------
# 先创建测试用户
curl -s -X POST "$BASE/api/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"stress_login\",\"email\":\"stress_login_${TIMESTAMP}@test.com\",\"phone\":\"13811111111\",\"password\":\"Test@123\"}" > /dev/null

# 清除限流
redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning KEYS 'ratelimit:*' | xargs -r redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning DEL > /dev/null 2>&1

cat > /tmp/wrk_login.lua << LUAEOF
request = function()
  local body = '{"email":"stress_login_${TIMESTAMP}@test.com","password":"Test@123"}'
  return wrk.format("POST", "/api/login", {
    ["Content-Type"] = "application/json"
  }, body)
end
LUAEOF
run_wrk "POST /api/login" "POST" "$BASE" "/tmp/wrk_login.lua"

#---------------------------------------------------
echo ""
echo "📌 5. API - 报名（wrk Lua）"
#---------------------------------------------------
# 清除限流
redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning KEYS 'ratelimit:*' | xargs -r redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning DEL > /dev/null 2>&1

cat > /tmp/wrk_apply.lua << 'LUAEOF'
local counter = 0
local thread_id = 0

function setup(thread)
  thread:set("id", thread_id)
  thread_id = thread_id + 1
end

function init(args)
  math.randomseed(os.time() + id * 1000)
end

request = function()
  counter = counter + 1
  local body = string.format(
    '{"studentName":"压测学生_%d_%d","birthDate":"2015-01-01","gender":"男","grade":%d,"parentName":"压测家长","phone":"139%08d","email":"apply_%d_%d_%d@load.test","currentSchool":"压测中学"}',
    id, counter, math.random(1,12), math.random(99999999), os.time(), id, counter)
  return wrk.format("POST", "/api/apply", {
    ["Content-Type"] = "application/json"
  }, body)
end
LUAEOF
run_wrk "POST /api/apply" "POST" "$BASE" "/tmp/wrk_apply.lua"

#---------------------------------------------------
echo ""
echo "📌 6. 混合负载（静态 + API）"
#---------------------------------------------------
# 清除限流
redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning KEYS 'ratelimit:*' | xargs -r redis-cli -h xinhang-redis-redis-redis.ns-0h7fttt7.svc -a '7mO7h4W2x4' --no-auth-warning DEL > /dev/null 2>&1

cat > /tmp/wrk_mixed.lua << 'LUAEOF'
local paths = {"/", "/about", "/academics", "/health", "/admission"}
local counter = 0

request = function()
  counter = counter + 1
  local path = paths[(counter % #paths) + 1]
  return wrk.format("GET", path)
end
LUAEOF
run_wrk "混合 GET 请求" "GET" "$BASE" "/tmp/wrk_mixed.lua"

#---------------------------------------------------
echo ""
echo "============================================"
echo "   压力测试完成"
echo "============================================"

# 最终统计
echo ""
echo "📌 数据库统计："
PGPASSWORD=7dl72vft psql -h xinhang-db-postgresql.ns-0h7fttt7.svc -U postgres -d xinhang -t -c "SELECT '用户数: ' || count(*) FROM users;" 2>/dev/null
PGPASSWORD=7dl72vft psql -h xinhang-db-postgresql.ns-0h7fttt7.svc -U postgres -d xinhang -t -c "SELECT '报名数: ' || count(*) FROM applications;" 2>/dev/null
