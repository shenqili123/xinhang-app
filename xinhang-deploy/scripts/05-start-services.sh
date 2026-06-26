#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤5：启动服务
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 启动服务"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
TARGET_DIR="/data/xinhang-app/backend"

# 检查前置条件
if [ ! -f "$TARGET_DIR/xinhang-backend" ]; then
    echo "[错误] 后端程序不存在，请先执行步骤3"
    exit 1
fi

if [ ! -f "$TARGET_DIR/.env" ]; then
    echo "[错误] 环境变量未配置，请先执行步骤4"
    exit 1
fi

# 安装 systemd 服务
echo "[1/4] 安装 systemd 服务..."

cat > /etc/systemd/system/xinhang.service << 'EOF'
[Unit]
Description=Xinhang School Website Backend
After=network.target postgresql.service redis-server.service

[Service]
Type=simple
User=root
WorkingDirectory=/data/xinhang-app/backend
ExecStart=/data/xinhang-app/backend/xinhang-backend
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
echo "      服务文件已安装"

# 启动后端服务
echo "[2/4] 启动后端服务..."
systemctl start xinhang
systemctl enable xinhang
sleep 2

if systemctl is-active --quiet xinhang; then
    echo "      [✓] 后端服务启动成功"
else
    echo "      [✗] 后端服务启动失败！"
    echo "      查看日志: journalctl -u xinhang -n 20"
    journalctl -u xinhang -n 10 --no-pager
    exit 1
fi

# 配置 Nginx
echo "[3/4] 配置 Nginx 反向代理..."

cat > /etc/nginx/sites-available/xinhang << 'EOF'
server {
    listen 80;
    server_name _;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
EOF

# 启用配置
ln -sf /etc/nginx/sites-available/xinhang /etc/nginx/sites-enabled/xinhang
rm -f /etc/nginx/sites-enabled/default

# 测试并启动 Nginx
if nginx -t 2>/dev/null; then
    systemctl start nginx
    systemctl reload nginx
    echo "      [✓] Nginx 配置成功"
else
    echo "      [✗] Nginx 配置有错误"
    nginx -t
    exit 1
fi

# 配置防火墙
echo "[4/4] 配置防火墙..."
if command -v ufw > /dev/null 2>&1; then
    ufw allow 22/tcp > /dev/null 2>&1 || true
    ufw allow 80/tcp > /dev/null 2>&1 || true
    ufw allow 443/tcp > /dev/null 2>&1 || true
    # 如果 ufw 未激活，不强制启用（避免把自己锁在外面）
    if ufw status | grep -q "inactive"; then
        echo "      防火墙未激活，已添加规则但未启用"
        echo "      （如需启用: sudo ufw enable）"
    else
        echo "      [✓] 防火墙规则已添加（22, 80, 443）"
    fi
else
    echo "      ufw 未安装，跳过防火墙配置"
fi
echo ""

# 测试访问
echo "--- 服务启动验证 ---"
sleep 1
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/ 2>/dev/null || echo "000")

if [ "$HTTP_CODE" = "200" ]; then
    echo "[✓] 网站通过 Nginx (端口80) 访问正常！"
elif [ "$HTTP_CODE" = "000" ]; then
    echo "[!] 无法连接，正在检查直连..."
    HTTP_CODE_DIRECT=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/ 2>/dev/null || echo "000")
    if [ "$HTTP_CODE_DIRECT" = "200" ]; then
        echo "    后端(8080)正常，Nginx 可能有问题"
    else
        echo "    后端也无法连接，请检查日志: journalctl -u xinhang -n 20"
    fi
else
    echo "[!] HTTP 返回码: $HTTP_CODE（预期 200）"
fi

echo ""

# 显示访问地址
echo "--- 访问地址 ---"
SERVER_IP=$(hostname -I | awk '{print $1}')
echo "服务器 IP: $SERVER_IP"
echo ""
echo "在内网其他电脑的浏览器打开："
echo "  http://$SERVER_IP"
echo ""

echo "========================================"
echo " 服务启动完成！"
echo "========================================"
echo ""
echo "下一步（导入新闻数据，可选）："
echo "  sudo bash scripts/06-import-data.sh"
echo ""
echo "常用管理命令："
echo "  查看状态: sudo systemctl status xinhang"
echo "  查看日志: sudo journalctl -u xinhang -f"
echo "  重启服务: sudo systemctl restart xinhang"
echo ""
