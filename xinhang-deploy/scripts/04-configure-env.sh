#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤4：配置环境变量
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 配置环境变量"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

TARGET_DIR="/data/xinhang-app/backend"
ENV_FILE="$TARGET_DIR/.env"

if [ ! -d "$TARGET_DIR" ]; then
    echo "[错误] 应用目录不存在，请先执行步骤3"
    exit 1
fi

echo "此脚本将引导你配置网站的运行参数。"
echo ""

# 生成 JWT 密钥
JWT_SECRET=$(openssl rand -hex 32)
echo "[✓] 已自动生成 JWT 密钥"
echo ""

# 获取数据库密码
echo "--- 数据库配置 ---"
read -s -p "请输入步骤2中设置的数据库密码: " DB_PASSWORD
echo ""
echo ""

# 邮件配置
echo "--- 邮件配置（用于发送注册验证码）---"
echo "如果暂时没有准备好邮箱信息，可以输入 'skip' 跳过"
echo "（跳过后网站可以运行，但注册功能中的邮件验证码无法发送）"
echo ""
read -p "邮箱地址（如 abc@qq.com，或输入 skip 跳过）: " SMTP_USER

if [ "$SMTP_USER" = "skip" ] || [ -z "$SMTP_USER" ]; then
    SMTP_HOST="smtp.qq.com"
    SMTP_PORT="465"
    SMTP_USER="placeholder@qq.com"
    SMTP_PASSWORD="placeholder"
    SMTP_FROM="placeholder@qq.com"
    echo "  已跳过邮件配置（后续可在 $ENV_FILE 中修改）"
else
    # 根据邮箱后缀自动判断 SMTP 服务器
    DOMAIN=$(echo "$SMTP_USER" | cut -d'@' -f2)
    case "$DOMAIN" in
        qq.com)
            SMTP_HOST="smtp.qq.com"
            SMTP_PORT="465"
            ;;
        163.com)
            SMTP_HOST="smtp.163.com"
            SMTP_PORT="465"
            ;;
        126.com)
            SMTP_HOST="smtp.126.com"
            SMTP_PORT="465"
            ;;
        *)
            read -p "SMTP 服务器地址: " SMTP_HOST
            read -p "SMTP 端口（通常465或587）: " SMTP_PORT
            ;;
    esac
    echo ""
    echo "SMTP 服务器: $SMTP_HOST:$SMTP_PORT"
    read -s -p "请输入邮箱授权码（不是登录密码！）: " SMTP_PASSWORD
    echo ""
    SMTP_FROM="$SMTP_USER"
fi
echo ""

# 验证 PIN
echo "--- 准考证验证密码 ---"
echo "这是教职工扫描准考证二维码时需要输入的验证密码"
read -p "验证密码（默认 xinhang2026，直接回车使用默认值）: " VERIFY_PIN
if [ -z "$VERIFY_PIN" ]; then
    VERIFY_PIN="xinhang2026"
fi
echo ""

# 生成 .env 文件
echo "[生成配置文件] $ENV_FILE"
cat > "$ENV_FILE" << EOF
# ========================================
# 山东新航实验外国语学校 - 后端环境变量配置
# 生成时间: $(date '+%Y-%m-%d %H:%M:%S')
# ========================================

# 服务端口
PORT=8080

# PostgreSQL 数据库
DB_HOST=localhost
DB_PORT=5432
DB_USER=xinhang_user
DB_PASSWORD=$DB_PASSWORD
DB_NAME=xinhang
DB_MAX_CONNS=50
DB_IDLE_CONNS=10

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 认证密钥（自动生成）
JWT_SECRET=$JWT_SECRET

# 邮件发送 (SMTP)
SMTP_HOST=$SMTP_HOST
SMTP_PORT=$SMTP_PORT
SMTP_USER=$SMTP_USER
SMTP_PASSWORD=$SMTP_PASSWORD
SMTP_FROM=$SMTP_FROM

# 准考证验证 PIN
VERIFY_PIN=$VERIFY_PIN
EOF

chmod 600 "$ENV_FILE"
echo "[✓] 配置文件已生成（权限已设为仅 root 可读）"
echo ""

echo "========================================"
echo " 环境变量配置完成！"
echo "========================================"
echo ""
echo "配置文件位置: $ENV_FILE"
echo "如需修改配置: nano $ENV_FILE"
echo ""
echo "下一步："
echo "  sudo bash scripts/05-start-services.sh"
echo ""
