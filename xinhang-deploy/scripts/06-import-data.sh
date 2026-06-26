#!/bin/bash
# ============================================
# 新航学校网站部署 - 步骤6：导入数据
# ============================================
set -e

echo "========================================"
echo " 新航学校网站 - 导入新闻数据和图片"
echo "========================================"
echo ""

if [ "$EUID" -ne 0 ]; then
    echo "[错误] 请使用 sudo 运行此脚本"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
DATA_DIR="$DEPLOY_DIR/data"
TARGET_DIR="/data/xinhang-app/backend"

# 检查后端是否运行（需要先启动一次以创建数据表）
if ! systemctl is-active --quiet xinhang; then
    echo "[提示] 后端服务未运行，正在启动（需要先运行一次以创建数据表）..."
    systemctl start xinhang
    sleep 3
    if ! systemctl is-active --quiet xinhang; then
        echo "[错误] 后端服务启动失败，请先排查问题"
        echo "  查看日志: journalctl -u xinhang -n 20"
        exit 1
    fi
fi

echo ""

# 导入新闻 SQL
echo "[1/3] 导入新闻文章数据..."
SQL_FILE="$DATA_DIR/seed_news.sql"

if [ ! -f "$SQL_FILE" ]; then
    echo "[警告] 找不到 $SQL_FILE，跳过新闻导入"
else
    echo "  请输入数据库密码:"
    export PGPASSWORD
    read -s -p "  数据库密码: " PGPASSWORD
    echo ""

    # 检查是否已导入过
    NEWS_COUNT=$(psql -h localhost -U xinhang_user -d xinhang -t -c "SELECT COUNT(*) FROM news;" 2>/dev/null | xargs)
    
    if [ -n "$NEWS_COUNT" ] && [ "$NEWS_COUNT" -gt 0 ]; then
        echo "  数据库中已有 $NEWS_COUNT 条新闻"
        read -p "  是否重新导入？这会覆盖现有新闻数据 (y/n): " REIMPORT
        if [ "$REIMPORT" != "y" ]; then
            echo "  跳过新闻数据导入"
        else
            psql -h localhost -U xinhang_user -d xinhang -c "DELETE FROM news;" > /dev/null 2>&1
            psql -h localhost -U xinhang_user -d xinhang -f "$SQL_FILE" > /dev/null 2>&1
            echo "  [✓] 新闻数据已重新导入"
        fi
    else
        psql -h localhost -U xinhang_user -d xinhang -f "$SQL_FILE" > /dev/null 2>&1
        echo "  [✓] 新闻数据导入完成"
    fi

    # 重置序列号（避免后续新增新闻时 ID 冲突）
    psql -h localhost -U xinhang_user -d xinhang -c \
        "SELECT setval('news_id_seq', (SELECT COALESCE(MAX(id), 1) FROM news));" > /dev/null 2>&1 || true
    echo "  [✓] 数据库序列号已重置"
    
    # 显示导入结果
    FINAL_COUNT=$(psql -h localhost -U xinhang_user -d xinhang -t -c "SELECT COUNT(*) FROM news;" 2>/dev/null | xargs)
    echo "  当前新闻总数: $FINAL_COUNT 条"
fi
echo ""

# 导入新闻图片
echo "[2/3] 导入新闻配图..."
IMAGES_TAR="$DATA_DIR/news-images.tar"

if [ ! -f "$IMAGES_TAR" ]; then
    echo "[警告] 找不到图片包: $IMAGES_TAR"
    echo "       新闻文章中的配图将无法显示"
    echo "       如果图片包在 U 盘上的其他位置，请手动执行:"
    echo "       tar -xf <图片包路径> -C $TARGET_DIR/uploads/migration/"
else
    IMAGES_DIR="$TARGET_DIR/uploads/migration/images"

    if [ -d "$IMAGES_DIR" ] && [ "$(ls $IMAGES_DIR 2>/dev/null | wc -l)" -gt 100 ]; then
        echo "  目标目录已有 $(ls $IMAGES_DIR | wc -l) 个文件"
        read -p "  是否重新解压覆盖？(y/n): " REEXTRACT
        if [ "$REEXTRACT" != "y" ]; then
            echo "  跳过图片解压"
        else
            echo "  解压中...（约 1.3GB，可能需要 1~2 分钟）"
            tar -xf "$IMAGES_TAR" -C "$TARGET_DIR/uploads/migration/"
            echo "  [✓] 图片解压完成"
        fi
    else
        mkdir -p "$TARGET_DIR/uploads/migration"
        echo "  解压中...（约 1.3GB，可能需要 1~2 分钟）"
        tar -xf "$IMAGES_TAR" -C "$TARGET_DIR/uploads/migration/"
        echo "  [✓] 图片解压完成"
    fi

    # 验证
    if [ -d "$IMAGES_DIR" ]; then
        IMG_COUNT=$(ls "$IMAGES_DIR" | wc -l)
        echo "  图片数量: $IMG_COUNT 张"
    fi
fi
echo ""

# 验证图片可访问
echo "[3/3] 验证数据是否正常..."

# 测试图片 HTTP 访问
SAMPLE_IMG=$(ls "$TARGET_DIR/uploads/migration/images/" 2>/dev/null | head -1)
if [ -n "$SAMPLE_IMG" ]; then
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/uploads/migration/images/$SAMPLE_IMG" 2>/dev/null || echo "000")
    if [ "$HTTP_CODE" = "200" ]; then
        echo "[✓] 新闻图片通过 HTTP 可正常访问"
    else
        echo "[!] 图片 HTTP 访问返回 $HTTP_CODE"
        echo "    可能需要重启后端: sudo systemctl restart xinhang"
    fi
fi

# 测试新闻 API
NEWS_API=$(curl -s "http://localhost:8080/api/news?page=1&page_size=1" 2>/dev/null)
if echo "$NEWS_API" | grep -q '"total"'; then
    echo "[✓] 新闻 API 正常返回数据"
else
    echo "[!] 新闻 API 可能有问题，请检查"
fi

unset PGPASSWORD

echo ""
echo "========================================"
echo " 数据导入完成！"
echo "========================================"
echo ""
echo "下一步（验收检查）："
echo "  sudo bash scripts/07-verify-all.sh"
echo ""
