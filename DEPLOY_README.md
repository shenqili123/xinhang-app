# 山东新航实验外国语学校 — 网站部署指南

## 项目简介

本项目为学校官方网站全栈应用，包含：
- **前端**：Vue 3 单页应用（中英双语）
- **后端**：Go (Gin) REST API 服务
- **数据库**：PostgreSQL
- **缓存**：Redis

部署后提供以下功能：学校官网展示、用户注册/登录、学生报名（含照片上传和准考证生成/下载）、新闻资讯系统、准考证验证等。

---

## 环境要求

| 组件 | 最低版本 | 说明 |
|------|----------|------|
| Go | 1.21+ | 后端编译运行 |
| Node.js | 18+ | 前端构建（运行时不需要） |
| PostgreSQL | 14+ | 主数据库 |
| Redis | 6+ | 验证码缓存 + API 限流 |
| 操作系统 | Linux / macOS / Windows | 均可 |

---

## 目录结构

```
xinhang-app/
├── DEPLOY_README.md         ← 本文件
├── package.json             # 前端依赖声明
├── vite.config.js           # 前端构建配置
├── index.html               # 前端入口 HTML
├── public/                  # 前端静态资源（图片等）
│   └── assets/
├── src/                     # 前端 Vue 源代码
│   ├── views/
│   ├── components/
│   ├── composables/
│   ├── router/
│   └── assets/
└── backend/                 # Go 后端源代码
    ├── main.go              # 后端入口
    ├── go.mod               # Go 依赖声明
    ├── go.sum               # Go 依赖锁定
    ├── .env.example         # 环境变量模板 ← 需复制为 .env 并填写
    ├── config/              # 配置加载
    ├── database/            # 数据库连接
    ├── cache/               # Redis 连接
    ├── email/               # 邮件发送
    ├── handlers/            # API 处理器
    ├── middleware/          # 中间件（鉴权/限流）
    ├── models/              # 数据模型
    ├── dist/                # ← 前端构建输出（部署时使用）
    └── uploads/             # ← 用户上传文件（运行时生成）
```

---

## 部署步骤

### 第一步：安装依赖环境

#### 安装 Go
```bash
# Ubuntu/Debian
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version  # 确认安装成功
```

#### 安装 Node.js（仅用于构建前端）
```bash
# 使用 nvm 或直接安装
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node --version  # 确认安装成功
npm --version
```

#### 安装 PostgreSQL
```bash
sudo apt install -y postgresql postgresql-contrib
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### 安装 Redis
```bash
sudo apt install -y redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

---

### 第二步：配置数据库

```bash
# 进入 PostgreSQL
sudo -u postgres psql

# 创建数据库和用户
CREATE DATABASE xinhang;
CREATE USER xinhang_user WITH PASSWORD '你的数据库密码';
GRANT ALL PRIVILEGES ON DATABASE xinhang TO xinhang_user;
\q
```

> 注意：表结构会在后端首次启动时自动创建（GORM AutoMigrate），无需手动建表。

---

### 第三步：配置环境变量

```bash
cd backend/
cp .env.example .env
```

编辑 `.env` 文件，填入你的实际配置：

```env
# 服务端口
PORT=8080

# 数据库（修改为你的实际信息）
DB_HOST=localhost
DB_PORT=5432
DB_USER=xinhang_user
DB_PASSWORD=你的数据库密码
DB_NAME=xinhang
DB_MAX_CONNS=50
DB_IDLE_CONNS=10

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 密钥（请更换为随机字符串）
JWT_SECRET=请替换为一个随机的32位以上字符串

# 邮件发送（用于注册验证码）
# 如果使用 QQ 邮箱，SMTP_PASSWORD 填授权码（非邮箱密码）
# 其他邮箱请相应修改 SMTP_HOST 和 SMTP_PORT
SMTP_HOST=smtp.qq.com
SMTP_PORT=465
SMTP_USER=your-email@qq.com
SMTP_PASSWORD=你的邮箱授权码
SMTP_FROM=your-email@qq.com

# 准考证验证密码（教职工扫码验证时使用）
VERIFY_PIN=xinhang2026
```

**关于邮箱配置：**
- QQ 邮箱：登录 QQ 邮箱 → 设置 → 账户 → 开启 SMTP → 获取授权码
- 163 邮箱：SMTP_HOST=smtp.163.com，SMTP_PORT=465
- 企业邮箱：根据你的邮件服务商配置

---

### 第四步：构建前端

```bash
# 在项目根目录（xinhang-app/）
npm install         # 安装前端依赖
npm run build       # 构建，输出到 backend/dist/
```

构建成功后，`backend/dist/` 目录会包含：
- `index.html`
- `assets/` (JS/CSS 文件)

---

### 第五步：编译并运行后端

```bash
cd backend/

# 下载 Go 依赖
go mod download

# 编译
go build -o xinhang-backend .

# 运行
./xinhang-backend
```

看到以下输出说明启动成功：
```
Database connected (maxOpen=50, maxIdle=10)
Database migration completed
Redis connected at localhost:6379
Server starting on :8080
```

访问 `http://你的服务器IP:8080` 即可看到网站。

---

### 第六步：创建管理员账号

启动后，先通过网站注册一个普通账号，然后手动将其升级为管理员：

```bash
sudo -u postgres psql -d xinhang

UPDATE users SET role = 'admin' WHERE email = '你注册的邮箱';
\q
```

管理员可以：
- 通过 API 查看所有报名记录 (`GET /api/applications`)
- 管理新闻（创建/编辑/删除）

---

## 生产环境部署建议

### 使用 systemd 管理服务

创建文件 `/etc/systemd/system/xinhang.service`：
```ini
[Unit]
Description=Xinhang School Website
After=network.target postgresql.service redis-server.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/xinhang-app/backend
ExecStart=/path/to/xinhang-app/backend/xinhang-backend
Restart=always
RestartSec=5
Environment=GIN_MODE=release

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable xinhang
sudo systemctl start xinhang
```

### 使用 Nginx 反向代理（推荐）

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 强制 HTTPS（如果有证书）
    # return 301 https://$host$request_uri;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 文件上传大小限制
        client_max_body_size 10M;
    }
}
```

### HTTPS 配置（Let's Encrypt）
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

---

## 数据导入

### 导入新闻数据（推荐）

我们已导出 300+ 篇从旧网站迁移的新闻文章，位于 `backend/database/seed_news.sql`。

首次启动后端（确保表已自动创建）后执行：

```bash
psql -h localhost -U xinhang_user -d xinhang -f backend/database/seed_news.sql
```

输入数据库密码后即可完成导入。导入后访问网站的"新闻"页面即可看到所有文章。

### 手动添加新闻（通过 API）

```bash
# 先登录管理员获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"your_password"}' | python3 -c "import sys,json;print(json.load(sys.stdin)['token'])")

# 创建新闻
curl -X POST http://localhost:8080/api/news \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "新闻标题",
    "summary": "摘要",
    "content": "<p>正文 HTML</p>",
    "category": "campus",
    "authorName": "作者"
  }'
```

---

## 常见问题

### Q: 启动时报 "Failed to connect to database"
A: 检查 PostgreSQL 是否启动，.env 中数据库信息是否正确。确保数据库已创建。

### Q: 启动时报 "Redis connected" 失败
A: 检查 Redis 是否启动。如果不需要 Redis，需要修改代码移除 Redis 依赖（不推荐）。

### Q: 邮件发送失败
A: 检查 SMTP 配置。QQ 邮箱需要开启 SMTP 服务并使用授权码（不是登录密码）。

### Q: 注册时收不到验证码
A: 1) 检查邮箱 SMTP 配置 2) 检查垃圾箱 3) QQ 邮箱授权码是否正确

### Q: 前端页面空白
A: 确保 `backend/dist/` 目录存在且包含 `index.html`。如果没有，重新执行 `npm run build`。

### Q: 照片上传失败
A: 确保 `backend/uploads/` 目录存在且有写入权限：`mkdir -p uploads/photos && chmod 755 uploads`

### Q: 如何修改端口
A: 修改 `.env` 中的 `PORT=8080` 为其他端口。

---

## API 文档

详细 API 说明请参考 `PROJECT.md` 文件。

---

## 技术支持

如有部署问题，请联系开发团队。

*版本：2.0*
*更新日期：2026-06-25*
