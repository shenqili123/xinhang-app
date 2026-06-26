#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤2：配置数据库
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 配置数据库"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

# 检查 PostgreSQL 是否运行
if ! systemctl is-active --quiet postgresql; then
    echo "[错误] PostgreSQL 未运行，请先执行步骤1"
    exit 1
fi

echo "此脚本将创建网站所需的数据库和用户。"
echo ""

# 获取密码
echo "请设置数据库密码（至少8位，包含字母和数字）："
echo "（这个密码后面配置环境变量时需要用到，请记住！）"
echo ""
read -s -p "输入密码: " DB_PASSWORD
echo ""
read -s -p "再输入一次确认: " DB_PASSWORD_CONFIRM
echo ""
echo ""

if [ "$DB_PASSWORD" != "$DB_PASSWORD_CONFIRM" ]; then
    echo "[错误] 两次输入的密码不一致，请重新运行此脚本"
    exit 1
fi

if [ ${#DB_PASSWORD} -lt 8 ]; then
    echo "[警告] 密码少于8位，建议使用更长的密码"
    read -p "是否继续？(y/n): " CONTINUE
    if [ "$CONTINUE" != "y" ]; then
        exit 1
    fi
fi

echo "[1/4] 创建数据库 'xinhang'..."
sudo -u postgres psql -c "CREATE DATABASE xinhang;" 2>/dev/null || echo "      (数据库可能已存在，跳过)"

echo "[2/4] 创建用户 'xinhang_user'..."
sudo -u postgres psql -c "CREATE USER xinhang_user WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null || \
    sudo -u postgres psql -c "ALTER USER xinhang_user WITH PASSWORD '$DB_PASSWORD';" 2>/dev/null
echo "      用户创建/更新成功"

echo "[3/4] 授予权限..."
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE xinhang TO xinhang_user;"
sudo -u postgres psql -d xinhang -c "GRANT ALL ON SCHEMA public TO xinhang_user;"
echo "      权限授予成功"

echo "[4/4] 配置认证方式..."
# 找到 pg_hba.conf 文件
PG_VERSION=$(ls /etc/postgresql/ | head -1)
PG_HBA="/etc/postgresql/$PG_VERSION/main/pg_hba.conf"

if [ -f "$PG_HBA" ]; then
    # 备份
    cp "$PG_HBA" "${PG_HBA}.backup"

    # 修改 local 认证方式为 md5
    sed -i 's/^local\s\+all\s\+all\s\+peer/local   all             all                                     md5/' "$PG_HBA"

    # 确保有 localhost 的 md5 认证行
    if ! grep -q "host.*all.*all.*127.0.0.1/32.*md5" "$PG_HBA"; then
        echo "host    all             all             127.0.0.1/32            md5" >> "$PG_HBA"
    fi

    # 重启 PostgreSQL 使配置生效
    systemctl restart postgresql
    echo "      认证配置已更新"
else
    echo "[警告] 找不到 pg_hba.conf，请手动配置"
fi

echo ""

# 验证连接
echo "--- 验证数据库连接 ---"
export PGPASSWORD="$DB_PASSWORD"
if psql -h localhost -U xinhang_user -d xinhang -c "SELECT 1;" > /dev/null 2>&1; then
    echo "[✓] 数据库连接成功！"
else
    echo "[✗] 数据库连接失败"
    echo "    请检查密码和认证配置"
    echo "    手动测试: psql -h localhost -U xinhang_user -d xinhang"
    unset PGPASSWORD
    exit 1
fi
unset PGPASSWORD

echo ""
echo "========================================"
echo " 数据库配置完成！"
echo "========================================"
echo ""
echo "请记住你的数据库密码: （已设置，不显示）"
echo ""
echo "下一步："
echo "  sudo bash scripts/03-deploy-app.sh"
echo ""
