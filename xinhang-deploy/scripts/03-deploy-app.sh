#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤3：部署应用文件
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 部署应用文件"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
APP_DIR="$DEPLOY_DIR/app"
TARGET_DIR="/data/xinhang-app/backend"

echo "源文件目录: $APP_DIR"
echo "目标目录:   $TARGET_DIR"
echo ""

# 检查源文件
if [ ! -f "$APP_DIR/xinhang-backend" ]; then
    echo "[错误] 找不到后端程序: $APP_DIR/xinhang-backend"
    exit 1
fi

if [ ! -d "$APP_DIR/dist" ]; then
    echo "[错误] 找不到前端文件: $APP_DIR/dist/"
    exit 1
fi

echo "[1/5] 创建目标目录..."
mkdir -p "$TARGET_DIR"
mkdir -p "$TARGET_DIR/uploads/photos"
mkdir -p "$TARGET_DIR/uploads/migration"
mkdir -p /data/backups
echo "      目录创建完成"

echo "[2/5] 复制后端程序..."
cp "$APP_DIR/xinhang-backend" "$TARGET_DIR/"
chmod +x "$TARGET_DIR/xinhang-backend"
echo "      $(ls -lh $TARGET_DIR/xinhang-backend | awk '{print $5}') 复制完成"

echo "[3/5] 复制前端文件..."
cp -r "$APP_DIR/dist" "$TARGET_DIR/"
echo "      $(du -sh $TARGET_DIR/dist | awk '{print $1}') 复制完成"

echo "[4/5] 复制配置模板..."
if [ -f "$APP_DIR/.env.template" ]; then
    cp "$APP_DIR/.env.template" "$TARGET_DIR/.env.template"
    echo "      .env.template 复制完成"
fi

echo "[5/5] 设置文件权限..."
chmod 755 "$TARGET_DIR"
chmod 755 "$TARGET_DIR/uploads"
chmod 755 "$TARGET_DIR/uploads/photos"
chmod 755 "$TARGET_DIR/uploads/migration"
echo "      权限设置完成"

echo ""
echo "--- 部署结果 ---"
echo "目录结构:"
echo "  $TARGET_DIR/"
echo "  ├── xinhang-backend    ($(ls -lh $TARGET_DIR/xinhang-backend | awk '{print $5}'))"
echo "  ├── dist/              (前端网页)"
echo "  ├── uploads/"
echo "  │   ├── photos/        (用户上传照片)"
echo "  │   └── migration/     (新闻配图，稍后导入)"
echo "  └── .env               (稍后配置)"
echo ""
echo "========================================"
echo " 应用文件部署完成！"
echo "========================================"
echo ""
echo "下一步："
echo "  sudo bash scripts/04-configure-env.sh"
echo ""
