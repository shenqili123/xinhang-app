#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤7：验收检查
# ============================================

echo "========================================"
echo " 新航学校网站 - 部署验收检查"
echo "========================================"
echo ""
echo "检查时间: $(date)"
echo ""

PASS=0
FAIL=0
WARN=0

check_pass() {
    echo "[✓] $1"
    PASS=$((PASS + 1))
}

check_fail() {
    echo "[✗] $1"
    FAIL=$((FAIL + 1))
}

check_warn() {
    echo "[!] $1"
    WARN=$((WARN + 1))
}

echo "--- 1. 系统服务状态 ---"

if systemctl is-active --quiet postgresql; then
    check_pass "PostgreSQL 运行中"
else
    check_fail "PostgreSQL 未运行"
fi

if systemctl is-active --quiet redis-server; then
    check_pass "Redis 运行中"
else
    check_fail "Redis 未运行"
fi

if systemctl is-active --quiet xinhang; then
    check_pass "网站后端服务运行中"
else
    check_fail "网站后端服务未运行"
fi

if systemctl is-active --quiet nginx; then
    check_pass "Nginx 运行中"
else
    check_fail "Nginx 未运行"
fi

echo ""
echo "--- 2. 端口监听 ---"

if ss -tlnp | grep -q ":5432"; then
    check_pass "PostgreSQL 监听 5432 端口"
else
    check_fail "PostgreSQL 未监听"
fi

if ss -tlnp | grep -q ":6379"; then
    check_pass "Redis 监听 6379 端口"
else
    check_fail "Redis 未监听"
fi

if ss -tlnp | grep -q ":8080"; then
    check_pass "后端监听 8080 端口"
else
    check_fail "后端未监听 8080"
fi

if ss -tlnp | grep -q ":80"; then
    check_pass "Nginx 监听 80 端口"
else
    check_fail "Nginx 未监听 80"
fi

echo ""
echo "--- 3. HTTP 访问测试 ---"

# 测试首页
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/ 2>/dev/null)
if [ "$HTTP_CODE" = "200" ]; then
    check_pass "首页访问正常 (HTTP $HTTP_CODE)"
else
    check_fail "首页访问失败 (HTTP $HTTP_CODE)"
fi

# 测试健康检查
HEALTH=$(curl -s http://localhost/health 2>/dev/null)
if echo "$HEALTH" | grep -q '"ok"'; then
    check_pass "健康检查接口正常"
else
    check_fail "健康检查接口异常"
fi

# 测试新闻 API
NEWS_RESP=$(curl -s "http://localhost/api/news?page=1&page_size=1" 2>/dev/null)
if echo "$NEWS_RESP" | grep -q '"total"'; then
    TOTAL=$(echo "$NEWS_RESP" | grep -o '"total":[0-9]*' | cut -d: -f2)
    check_pass "新闻 API 正常（共 $TOTAL 篇文章）"
else
    check_warn "新闻 API 无数据（可能未导入新闻）"
fi

# 测试新闻图片
SAMPLE_IMG=$(ls /data/xinhang-app/backend/uploads/migration/images/ 2>/dev/null | head -1)
if [ -n "$SAMPLE_IMG" ]; then
    IMG_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/uploads/migration/images/$SAMPLE_IMG" 2>/dev/null)
    if [ "$IMG_CODE" = "200" ]; then
        check_pass "新闻配图可正常加载"
    else
        check_warn "新闻配图加载失败 (HTTP $IMG_CODE)"
    fi
else
    check_warn "未找到新闻配图文件"
fi

# 测试静态资源
STATIC_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/assets/ 2>/dev/null)
# assets 目录不直接访问，测试首页包含的 JS
JS_FILE=$(ls /data/xinhang-app/backend/dist/assets/*.js 2>/dev/null | head -1 | xargs basename 2>/dev/null)
if [ -n "$JS_FILE" ]; then
    JS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/assets/$JS_FILE" 2>/dev/null)
    if [ "$JS_CODE" = "200" ]; then
        check_pass "前端 JS 资源加载正常"
    else
        check_fail "前端 JS 资源加载失败"
    fi
fi

echo ""
echo "--- 4. 数据检查 ---"

# 检查数据库中的表
TABLE_COUNT=$(sudo -u postgres psql -d xinhang -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | xargs)
if [ -n "$TABLE_COUNT" ] && [ "$TABLE_COUNT" -gt 0 ]; then
    check_pass "数据库有 $TABLE_COUNT 个表"
else
    check_warn "数据库表可能未创建"
fi

# 检查新闻数量
NEWS_COUNT=$(sudo -u postgres psql -d xinhang -t -c "SELECT COUNT(*) FROM news;" 2>/dev/null | xargs)
if [ -n "$NEWS_COUNT" ] && [ "$NEWS_COUNT" -gt 0 ]; then
    check_pass "新闻数据已导入（$NEWS_COUNT 条）"
else
    check_warn "新闻数据未导入"
fi

# 检查图片文件
IMG_COUNT=$(ls /data/xinhang-app/backend/uploads/migration/images/ 2>/dev/null | wc -l)
if [ "$IMG_COUNT" -gt 100 ]; then
    check_pass "新闻配图已就位（$IMG_COUNT 张）"
elif [ "$IMG_COUNT" -gt 0 ]; then
    check_warn "新闻配图数量偏少（$IMG_COUNT 张，预期 2000+）"
else
    check_warn "新闻配图未导入"
fi

echo ""
echo "--- 5. 开机自启 ---"

if systemctl is-enabled --quiet postgresql 2>/dev/null; then
    check_pass "PostgreSQL 已设为开机自启"
else
    check_warn "PostgreSQL 未设为开机自启"
fi

if systemctl is-enabled --quiet redis-server 2>/dev/null; then
    check_pass "Redis 已设为开机自启"
else
    check_warn "Redis 未设为开机自启"
fi

if systemctl is-enabled --quiet xinhang 2>/dev/null; then
    check_pass "网站后端已设为开机自启"
else
    check_warn "网站后端未设为开机自启"
fi

if systemctl is-enabled --quiet nginx 2>/dev/null; then
    check_pass "Nginx 已设为开机自启"
else
    check_warn "Nginx 未设为开机自启"
fi

echo ""
echo "========================================"
echo " 检查结果汇总"
echo "========================================"
echo ""
echo " 通过: $PASS 项"
echo " 失败: $FAIL 项"
echo " 警告: $WARN 项"
echo ""

if [ "$FAIL" -eq 0 ]; then
    echo "🎉 部署验收通过！"
    echo ""
    SERVER_IP=$(hostname -I | awk '{print $1}')
    echo "=== 网站访问地址 ==="
    echo ""
    echo "  http://$SERVER_IP"
    echo ""
    echo "在学校内网的任何电脑浏览器中输入以上地址即可访问。"
    echo ""
    echo "=== 后续操作 ==="
    echo ""
    echo "1. 创建管理员账号："
    echo "   - 先在网站上注册一个普通账号"
    echo "   - 然后执行: sudo -u postgres psql -d xinhang"
    echo "   - 输入: UPDATE users SET role='admin' WHERE email='你的邮箱';"
    echo ""
    echo "2. 如果需要修改配置："
    echo "   nano /data/xinhang-app/backend/.env"
    echo "   sudo systemctl restart xinhang"
    echo ""
    echo "3. 查看日志："
    echo "   sudo journalctl -u xinhang -f"
    echo ""
else
    echo "⚠️  有 $FAIL 项检查未通过，请根据上面的错误信息排查"
    echo ""
    echo "常用排错命令："
    echo "  sudo journalctl -u xinhang -n 50    # 查看后端日志"
    echo "  sudo systemctl status postgresql     # 查看数据库状态"
    echo "  sudo nginx -t                        # 检查 Nginx 配置"
    echo "  cat /data/xinhang-app/backend/.env   # 查看配置文件"
fi
echo ""
