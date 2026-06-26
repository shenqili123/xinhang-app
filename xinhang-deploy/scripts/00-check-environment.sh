#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤0：环境检查
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 环境检查"
echo "========================================"
echo ""

# 检查是否以 root 运行
if [ "$EUID" -ne 0 ]; then
    echo "[提示] 请使用 sudo 运行此脚本"
    echo "  用法: sudo bash $0"
    exit 1
fi

echo "--- 系统信息 ---"
echo "操作系统: $(cat /etc/os-release | grep PRETTY_NAME | cut -d'"' -f2)"
echo "内核版本: $(uname -r)"
echo "CPU 架构: $(uname -m)"
echo "主机名:   $(hostname)"
echo ""

echo "--- CPU 信息 ---"
echo "CPU 核心数: $(nproc)"
echo "CPU 型号:   $(grep 'model name' /proc/cpuinfo | head -1 | cut -d':' -f2 | xargs)"
echo ""

echo "--- 内存信息 ---"
free -h | head -2
echo ""

echo "--- 磁盘空间 ---"
df -h / | tail -1 | awk '{printf "总容量: %s, 已用: %s, 可用: %s, 使用率: %s\n", $2, $3, $4, $5}'
echo ""

# 检查磁盘空间是否充足（至少需要 5GB）
AVAIL_KB=$(df / | tail -1 | awk '{print $4}')
if [ "$AVAIL_KB" -lt 5000000 ]; then
    echo "[警告] 可用磁盘空间不足 5GB，部署可能会有问题！"
else
    echo "[✓] 磁盘空间充足"
fi
echo ""

echo "--- 网络检查 ---"
echo "网络接口:"
ip -brief addr show | grep -v "^lo"
echo ""

# 测试网络连通性
if ping -c 1 -W 3 mirrors.aliyun.com > /dev/null 2>&1; then
    echo "[✓] 可以访问外部网络（mirrors.aliyun.com）"
    echo "    建议使用在线安装模式（更快）"
    NETWORK_OK=1
else
    echo "[!] 无法访问外部网络"
    echo "    将使用离线安装模式"
    NETWORK_OK=0
fi
echo ""

echo "--- Ubuntu 版本检查 ---"
CODENAME=$(lsb_release -cs 2>/dev/null || echo "unknown")
VERSION=$(lsb_release -rs 2>/dev/null || echo "unknown")
echo "版本: Ubuntu $VERSION ($CODENAME)"
if [ "$CODENAME" = "noble" ]; then
    echo "[✓] Ubuntu 24.04，与离线安装包兼容"
elif [ "$CODENAME" = "jammy" ]; then
    echo "[警告] Ubuntu 22.04，离线安装包可能不兼容！"
    echo "       如果有网络，可以使用在线安装；否则需要 22.04 版本的离线包"
else
    echo "[警告] 非预期的 Ubuntu 版本，请联系技术支持"
fi
echo ""

echo "--- 已安装的相关软件 ---"
for pkg in postgresql redis-server nginx; do
    if command -v $pkg > /dev/null 2>&1 || dpkg -l | grep -q "^ii.*$pkg"; then
        echo "[已安装] $pkg"
    else
        echo "[未安装] $pkg （需要安装）"
    fi
done
echo ""

echo "========================================"
echo " 环境检查完成"
echo "========================================"
echo ""
echo "如果一切正常，请执行下一步："
echo "  sudo bash scripts/01-install-software.sh"
echo ""
