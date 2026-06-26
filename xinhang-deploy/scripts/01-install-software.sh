#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤1：安装基础软件
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 安装基础软件"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

# 设置时区
echo "[1/4] 设置系统时区为 Asia/Shanghai..."
timedatectl set-timezone Asia/Shanghai
echo "      当前时间: $(date)"
echo ""

# 检测网络
echo "[2/4] 检测网络连接..."
OFFLINE_MODE=0
if ping -c 1 -W 3 mirrors.aliyun.com > /dev/null 2>&1; then
    echo "      网络可用，将使用在线安装模式"
else
    echo "      网络不可用，将使用离线安装模式"
    OFFLINE_MODE=1
fi
echo ""

# 安装软件
echo "[3/4] 安装 PostgreSQL、Redis、Nginx..."
echo ""

if [ "$OFFLINE_MODE" -eq 0 ]; then
    # 在线模式：配置镜像源并 apt install
    echo "  → 配置阿里云镜像源..."
    cp /etc/apt/sources.list /etc/apt/sources.list.backup 2>/dev/null || true

    CODENAME=$(lsb_release -cs)
    cat > /etc/apt/sources.list << EOF
deb http://mirrors.aliyun.com/ubuntu/ ${CODENAME} main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ ${CODENAME}-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ ${CODENAME}-backports main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ ${CODENAME}-security main restricted universe multiverse
EOF

    echo "  → 更新软件包索引..."
    apt-get update -qq

    echo "  → 安装 PostgreSQL..."
    apt-get install -y -qq postgresql postgresql-contrib postgresql-client > /dev/null

    echo "  → 安装 Redis..."
    apt-get install -y -qq redis-server redis-tools > /dev/null

    echo "  → 安装 Nginx..."
    apt-get install -y -qq nginx > /dev/null

    echo "  → 安装辅助工具..."
    apt-get install -y -qq curl wget unzip tar htop net-tools > /dev/null

else
    # 离线模式：使用预下载的 deb 包
    PKGDIR="$DEPLOY_DIR/offline-packages"
    if [ ! -d "$PKGDIR" ] || [ -z "$(ls $PKGDIR/*.deb 2>/dev/null)" ]; then
        echo "[错误] 找不到离线安装包目录: $PKGDIR"
        echo "       请确保 offline-packages/ 文件夹中有 .deb 文件"
        exit 1
    fi

    echo "  → 使用离线安装包 ($PKGDIR)..."
    echo "  → 共 $(ls $PKGDIR/*.deb | wc -l) 个包"

    # dpkg 安装所有包（忽略依赖错误，后面 fix）
    dpkg -i $PKGDIR/*.deb 2>/dev/null || true

    # 尝试修复依赖关系（如果有网络的话）
    apt-get -f install -y 2>/dev/null || true

    echo "  → 离线安装完成"
fi
echo ""

# 启动服务
echo "[4/4] 启动并设置开机自启..."

echo "  → PostgreSQL..."
systemctl start postgresql 2>/dev/null || true
systemctl enable postgresql 2>/dev/null || true

echo "  → Redis..."
systemctl start redis-server 2>/dev/null || true
systemctl enable redis-server 2>/dev/null || true

echo "  → Nginx（暂时先不配置）..."
systemctl stop nginx 2>/dev/null || true
systemctl enable nginx 2>/dev/null || true

echo ""

# 验证
echo "--- 安装结果验证 ---"
ALL_OK=1

if systemctl is-active --quiet postgresql; then
    echo "[✓] PostgreSQL 运行中"
else
    echo "[✗] PostgreSQL 未运行"
    ALL_OK=0
fi

if systemctl is-active --quiet redis-server; then
    echo "[✓] Redis 运行中"
else
    echo "[✗] Redis 未运行"
    ALL_OK=0
fi

if command -v nginx > /dev/null 2>&1; then
    echo "[✓] Nginx 已安装"
else
    echo "[✗] Nginx 安装失败"
    ALL_OK=0
fi

echo ""
if [ "$ALL_OK" -eq 1 ]; then
    echo "========================================"
    echo " 软件安装完成！"
    echo "========================================"
    echo ""
    echo "下一步："
    echo "  sudo bash scripts/02-setup-database.sh"
else
    echo "[警告] 部分软件安装失败，请检查上面的错误信息"
    echo "       可以尝试重新运行此脚本"
fi
echo ""
