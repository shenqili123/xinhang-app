# 山东新航实验国际学校 — 全栈网站项目

## 项目概述

将学校原有的静态 HTML 官网重构为 **Vue 3 + Go + PostgreSQL** 全栈应用，部署在 **Sealos** 云平台上。支持中英双语、用户注册登录、学生报名（含照片上传和准考证生成）、新闻资讯系统等完整功能。

---

## 当前进度总览

| 阶段 | 状态 | 说明 |
|------|------|------|
| 静态 HTML 官网 | ✅ 已弃用 | 原始 6 个页面，曾部署在上码平台 |
| Vue3 前端重构 | ✅ 已完成 | 完全按同事提供的新设计稿重写，6 展示页 + 7 功能页 |
| Go 后端 API | ✅ 已完成 | Gin + GORM + Redis + 照片上传 + 准考证生成 |
| PostgreSQL 数据库 | ✅ 运行中 | 4 张表：users, applications, news, categories |
| Redis 缓存/限流 | ✅ 运行中 | 验证码缓存 + API 频率限制 |
| 邮件验证系统 | ✅ 已完成 | QQ 邮箱 SMTP，注册时发送验证码 |
| 新闻资讯系统 | ✅ 已完成 | 从旧站迁移 300+ 篇文章，支持分类/分页/详情 |
| 准考证系统 | ✅ 已完成 | 报名后自动生成传统风格准考证 + PDF 下载 |
| 照片上传 | ✅ 已完成 | 报名时上传考生照片，存储于服务器 |
| 身份证校验 | ✅ 已完成 | 18位身份证号格式+校验位验证 |
| 生产部署 | ✅ 运行中 | Go serve 前端 dist（方案 A），8080 端口 |

---

## 技术架构图

```
用户浏览器 (PC/Mobile)
    │
    ▼
┌──────────────────────────────────────────────────────────────────┐
│                    Go 后端 (Gin) :8080                            │
│  ┌─────────────────┐    ┌────────────────┐    ┌──────────────┐  │
│  │ 静态文件服务     │    │ REST API 路由   │    │ 文件上传服务  │  │
│  │ /assets, /images│    │ /api/*         │    │ /uploads/*   │  │
│  │ dist/index.html │    │ JSON 请求/响应  │    │ photos/      │  │
│  └─────────────────┘    └────────────────┘    └──────────────┘  │
│           │                      │                     │         │
│           ▼                      ▼                     ▼         │
│  ┌──────────────┐    ┌──────────────┐    ┌───────────────────┐  │
│  │ Vue3 SPA     │    │ 业务逻辑层    │    │ 本地文件系统       │  │
│  │ 前端 dist    │    │ handlers/    │    │ ./uploads/photos/ │  │
│  └──────────────┘    └──────────────┘    └───────────────────┘  │
│                              │                                   │
│              ┌───────────────┼───────────────┐                   │
│              ▼               ▼               ▼                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ PostgreSQL   │  │    Redis     │  │ QQ SMTP      │          │
│  │ 持久化数据    │  │ 缓存/限流    │  │ 邮件发送      │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└──────────────────────────────────────────────────────────────────┘
      Sealos 内网                Sealos 内网           外部服务
```

---

## 前端架构

### 技术栈
- **Vue 3.5** + **Vite 8** + **Vue Router 4**
- **中英双语**：自定义 composable `useLanguage` + `t(en, zh)` 函数
- **滚动动画**：自定义 composable `useReveal` (IntersectionObserver)
- **样式**：纯 CSS，无 UI 框架依赖
- **PDF 生成**：html2canvas + jsPDF（准考证下载）
- **二维码**：后端 Go 生成 QR Code 图片

### 目录结构

```
xinhang-app/
├── index.html
├── vite.config.js                  # 构建输出到 backend/dist
├── package.json
├── public/
│   └── assets/                     # 静态图片资源 (jpg/png/svg)
├── src/
│   ├── main.js                     # Vue 挂载入口
│   ├── App.vue                     # 根组件 (Header + RouterView + Footer)
│   ├── assets/
│   │   └── styles.css              # 全局样式 (~3100行)
│   ├── router/
│   │   └── index.js                # 路由配置 + 守卫 (auth/guest)
│   ├── composables/
│   │   ├── useLanguage.js          # 中英切换 (localStorage 持久化)
│   │   ├── useAuth.js              # 认证状态管理 (token + user)
│   │   └── useReveal.js            # 滚动入场动画
│   ├── components/
│   │   ├── SiteHeader.vue          # 导航栏 (子菜单 + 语言切换 + 登录状态)
│   │   └── SiteFooter.vue          # 页脚
│   └── views/
│       ├── HomeView.vue            # 首页
│       ├── AboutView.vue           # 关于新航
│       ├── AcademicsView.vue       # 学术课程
│       ├── CampusView.vue          # 校园展示
│       ├── StudentLifeView.vue     # 学生生活
│       ├── ApplyView.vue           # 报名申请 + 准考证生成
│       ├── RegisterView.vue        # 用户注册
│       ├── LoginView.vue           # 用户登录
│       ├── ProfileView.vue         # 个人中心
│       ├── NewsView.vue            # 新闻列表
│       ├── NewsDetailView.vue      # 新闻详情
│       ├── VerifyView.vue          # 准考证验证 (教职工)
│       └── AdmissionView.vue       # 招生信息 (redirect → /apply)
└── backend/                        # Go 后端 (同仓库)
```

### 路由表

| 路径 | 页面 | 权限 | 说明 |
|------|------|------|------|
| `/` | HomeView | 公开 | 首页 |
| `/about` | AboutView | 公开 | 关于新航 |
| `/academics` | AcademicsView | 公开 | 学术课程 |
| `/campus` | CampusView | 公开 | 校园展示 |
| `/student-life` | StudentLifeView | 公开 | 学生生活 |
| `/apply` | ApplyView | 需登录 | 报名 + 准考证 |
| `/register` | RegisterView | 仅游客 | 用户注册 |
| `/login` | LoginView | 仅游客 | 用户登录 |
| `/profile` | ProfileView | 需登录 | 个人中心 |
| `/news` | NewsView | 公开 | 新闻列表 |
| `/news/:id` | NewsDetailView | 公开 | 新闻详情 |
| `/verify` | VerifyView | 公开 | 准考证验证 |

### 前端认证机制

```
useAuth.js composable:
├── user (ref)         ← localStorage('xinhang-user') JSON
├── token (ref)        ← localStorage('xinhang-token') JWT string
├── isLoggedIn (computed)
├── authHeader()       → { Authorization: 'Bearer xxx' }
├── setAuth(token, user) → 写入 localStorage + 更新 ref
└── logout()           → 清空 localStorage + 重置 ref
```

路由守卫逻辑：
- `meta.requiresAuth: true` → 未登录跳转 `/login?redirect=原路径`
- `meta.guestOnly: true` → 已登录跳转 `/profile`

---

## 后端架构

### 技术栈
- **Go 1.25** + **Gin** (HTTP 框架)
- **GORM** (ORM，自动迁移数据库)
- **PostgreSQL 16** (主数据库)
- **Redis** (验证码缓存 + API 限流)
- **JWT** (golang-jwt/v5，72小时有效期)
- **bcrypt** (密码加密)
- **go-qrcode** (准考证二维码生成)
- **QQ SMTP** (邮件发送)

### 目录结构

```
backend/
├── main.go              # 入口：路由注册 + 中间件 + 静态文件服务
├── xinhang-backend      # 编译后的可执行文件
├── config/
│   └── config.go        # 环境变量加载 (PORT, DB, Redis, SMTP, JWT)
├── database/
│   └── database.go      # PostgreSQL 连接池 (maxOpen=100, maxIdle=20)
├── cache/
│   └── redis.go         # Redis 客户端初始化
├── email/
│   └── email.go         # SMTP 邮件发送 (验证码)
├── queue/
│   └── kafka.go         # Redis Streams 异步消息队列
├── middleware/
│   ├── auth.go          # JWT 鉴权 + AdminOnly 中间件
│   ├── ratelimit.go     # Redis 滑动窗口限流
│   └── ip.go            # IP 保护中间件
├── handlers/
│   ├── auth.go          # 注册/登录/Profile CRUD
│   ├── application.go   # 报名提交 + 身份证校验 + 考场分配
│   ├── upload.go        # 照片上传处理
│   ├── qrcode.go        # 准考证二维码生成
│   ├── verification.go  # 准考证验证 (教职工)
│   └── news.go          # 新闻 CRUD + 分类统计
├── models/
│   ├── user.go          # User 数据模型
│   ├── application.go   # Application 数据模型
│   └── news.go          # News 数据模型
├── dist/                # 前端构建输出 (Vite build)
│   ├── index.html
│   └── assets/          # JS/CSS 带 hash 文件名
└── uploads/
    └── photos/          # 考生照片存储目录
```

### API 接口列表

| 方法 | 路径 | 限流 | 权限 | 说明 |
|------|------|------|------|------|
| POST | `/api/send-code` | 5/min | 公开 | 发送邮箱验证码 |
| POST | `/api/register` | 10/min | 公开 | 用户注册 |
| POST | `/api/login` | 20/min | 公开 | 用户登录 → JWT |
| POST | `/api/upload-photo` | 10/min | 需登录 | 上传考生照片 |
| POST | `/api/apply` | 5/min | 需登录 | 提交报名申请 |
| GET | `/api/profile` | — | 需登录 | 获取个人信息 |
| PUT | `/api/profile` | — | 需登录 | 更新个人信息 |
| GET | `/api/my-applications` | — | 需登录 | 我的报名记录 |
| GET | `/api/applications` | — | 管理员 | 所有报名列表 |
| GET | `/api/permit-qr` | — | 公开 | 生成准考证二维码图片 |
| POST | `/api/verify-qr` | — | 公开 | 验证准考证 (需PIN) |
| GET | `/api/query-permit` | — | 公开 | 按手机号/准考证号查询 |
| GET | `/api/news` | — | 公开 | 新闻列表 (分页+分类) |
| GET | `/api/news/:id` | — | 公开 | 新闻详情 |
| GET | `/api/news-categories` | — | 公开 | 新闻分类统计 |
| POST | `/api/news` | — | 管理员 | 创建新闻 |
| PUT | `/api/news/:id` | — | 管理员 | 编辑新闻 |
| DELETE | `/api/news/:id` | — | 管理员 | 删除新闻 |

---

## 数据库设计

### 连接方式
- **ORM**: GORM（自动迁移，启动时自动创建/更新表结构）
- **连接字符串**: `host=xxx port=5432 user=xxx password=xxx dbname=xinhang sslmode=disable TimeZone=Asia/Shanghai`
- **连接池**: maxOpen=100, maxIdle=20

### 表结构

#### users 表
```sql
CREATE TABLE users (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL,
    email          VARCHAR(200) UNIQUE NOT NULL,
    phone          VARCHAR(20),
    password_hash  VARCHAR(200) NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    role           VARCHAR(20) DEFAULT 'user',  -- 'user' | 'admin'
    created_at     TIMESTAMP DEFAULT NOW()
);
```

#### applications 表
```sql
CREATE TABLE applications (
    id             SERIAL PRIMARY KEY,
    user_id        INTEGER REFERENCES users(id),
    permit_no      VARCHAR(30) UNIQUE NOT NULL,    -- 准考证号 XH+日期+序号
    student_name   VARCHAR(100) NOT NULL,
    birth_date     VARCHAR(20),
    gender         VARCHAR(20) NOT NULL,           -- 'Male' | 'Female'
    grade          INTEGER NOT NULL,               -- 1=一年级, 7=七年级, 10=高一, 0=插班
    id_number      VARCHAR(30),                    -- 身份证号 (18位)
    photo          VARCHAR(500),                   -- 考生照片路径
    exam_room      VARCHAR(10),                    -- 考场号 (自动分配)
    seat_number    VARCHAR(10),                    -- 座位号 (自动分配)
    boarding_need  VARCHAR(30),
    parent_name    VARCHAR(100) NOT NULL,
    phone          VARCHAR(20) NOT NULL,
    relationship   VARCHAR(30),
    email          VARCHAR(200),
    current_school VARCHAR(200),
    track          VARCHAR(100),
    visit_date     VARCHAR(20),
    notes          TEXT,
    status         VARCHAR(20) DEFAULT 'pending',  -- pending|approved|rejected|exam
    created_at     TIMESTAMP DEFAULT NOW()
);
```

#### news 表
```sql
CREATE TABLE news (
    id           SERIAL PRIMARY KEY,
    title        VARCHAR(500) NOT NULL,
    summary      TEXT,
    content      TEXT,
    category     VARCHAR(50),         -- campus|news|teachers|activities|...
    author_name  VARCHAR(100),
    cover_image  VARCHAR(500),
    source_url   VARCHAR(500),
    published_at TIMESTAMP,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);
```

---

## 核心功能详解

### 1. 注册/登录系统

**注册流程：**
1. 用户填写姓名、手机、邮箱、密码
2. 点击"发送验证码" → 后端生成6位随机码存入 Redis (5分钟有效)
3. 通过 QQ SMTP 发送验证码邮件
4. 用户输入验证码 → 后端校验 Redis 中的码
5. 校验通过 → bcrypt 加密密码 → 写入 users 表
6. 返回 JWT token + user 信息

**登录流程：**
1. 用户输入邮箱+密码
2. 后端查询 users 表，bcrypt.Compare 密码
3. 匹配成功 → 生成 JWT (72h 有效期) → 返回 token + user

**技术要点：**
- JWT 载荷: `{ userId, role, exp }`
- 密码: bcrypt cost=10
- 验证码: 6位数字，Redis key = `verify:email`，TTL 5分钟
- 限流: 注册 10次/分钟，登录 20次/分钟

### 2. 报名系统（含准考证生成）

**报名流程：**
1. 用户填写学生信息（姓名、性别、年级、身份证号、学校等）
2. 上传考生照片（JPG/PNG，≤5MB）→ 存储在服务器 `./uploads/photos/`
3. 前端校验身份证号（18位格式 + 校验位算法）
4. 后端再次校验身份证号（双重验证）
5. 自动生成准考证号：`XH` + 日期(20060102) + 4位序号
6. 自动分配考场号和座位号（每30人一间考场，按年级分组）
7. 写入 applications 表
8. 返回准考证号 + 考场信息
9. 前端生成传统风格准考证页面 + 二维码
10. 支持下载为 PDF 文件

**身份证校验算法：**
```
位权: [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
校验码表: "10X98765432"
计算: 前17位加权求和 % 11 → 查校验码表 → 与第18位比对
```

**照片上传：**
- 接口: `POST /api/upload-photo` (multipart/form-data)
- 限制: JPG/PNG，最大 5MB
- 存储: `./uploads/photos/{timestamp}_{random}.{ext}`
- 返回: `{ url: "/uploads/photos/xxx.jpg" }`

**准考证样式：**
仿传统高考准考证设计，包含：
- 学校名称标题 + 年级对应考试子标题
- 考生照片（右侧）
- 信息字段（准考证号、姓名、性别、证件号码、报考年级、考点、考场、座位）
- 二维码（用于教职工扫码验证）
- 考试时间表（笔试一/二、面试一/二）
- 注意事项

### 3. 新闻资讯系统

**数据来源：** 从旧网站迁移 300+ 篇文章（HTML 内容清洗后导入）

**分类体系：**
- 校园动态 (campus)
- 新闻报道 (news)
- 名师风采 (teachers)
- 学生活动 (activities)
- 学子风采 (highlights)
- 媒体聚焦 (media)
- 教育科研 (research)
- 德育天地 (moral_education)
- 招生专栏 (admission)
- 通知公告 (notice)
- 等 15+ 个分类

**前端功能：**
- 分类筛选标签（显示各分类文章数量）
- 分页加载（每页 12 篇）
- 新闻卡片（封面图 + 标题 + 摘要 + 作者 + 日期）
- 详情页（全文 HTML 渲染，DOMPurify 安全过滤）

### 4. 准考证验证系统

**用途：** 考试当天，教职工扫描准考证二维码验证真伪

**流程：**
1. 教职工打开 `/verify` 页面
2. 输入验证 PIN 码 + 扫码得到的文本
3. 后端校验 PIN 正确性 + 查询 applications 表
4. 返回验证结果（有效/无效 + 学生信息）

---

## 前后端连接方式

### 开发环境
```javascript
// vite.config.js
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080'  // 开发时代理到后端
    }
  },
  build: {
    outDir: '../backend/dist'  // 构建输出到后端目录
  }
})
```

### 生产环境
Go 后端同时 serve 前端静态文件（方案 A，一体化部署）：
```go
// main.go
r.Static("/assets", "./dist/assets")     // 前端 JS/CSS
r.Static("/uploads", "./uploads")         // 用户上传文件
r.NoRoute(func(c *gin.Context) {
    c.File("./dist/index.html")           // SPA 路由回退
})
```

### 请求流程示例（报名）
```
1. 浏览器 → POST /api/upload-photo (FormData: photo file)
   ← { url: "/uploads/photos/1234_abc.jpg" }

2. 浏览器 → POST /api/apply (JSON + Authorization header)
   请求体: { studentName, gender, grade, idNumber, photo, ... }
   ← { permitNo: "XH202606258752", examRoom: "01", seatNumber: "15" }

3. 浏览器 → GET /api/permit-qr?no=XH202606258752&student=张三
   ← PNG 图片 (QR Code)
```

---

## 环境变量配置

```env
# 服务器
PORT=8080

# 数据库
DB_HOST=xinhang-pg-postgresql.ns-xxx.svc
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=xxxxxxxx
DB_NAME=xinhang

# Redis
REDIS_ADDR=xinhang-redis-redis-redis.ns-xxx.svc:6379
REDIS_PASSWORD=xxxxxxxx

# JWT
JWT_SECRET=your-secret-key

# 邮件
SMTP_HOST=smtp.qq.com
SMTP_PORT=465
SMTP_USER=1486308808@qq.com
SMTP_PASS=xxxxxxxx (QQ邮箱授权码)
```

---

## 部署方式

### Sealos 云平台
- **Go 后端**: DevBox 容器（Go 1.25 + Node 22）
- **PostgreSQL**: Sealos 数据库服务 (16.4, 0.5C/512M)
- **Redis**: Sealos Redis 服务

### 构建与启动
```bash
# 前端构建
cd xinhang-app && npm run build
# → 输出到 backend/dist/

# 后端编译
cd backend && go build -o xinhang-backend .

# 启动服务
cd backend && ./xinhang-backend
# Server starting on :8080
```

### 数据库迁移
GORM AutoMigrate 自动处理，启动时自动：
- 创建缺失的表
- 添加新增字段（如 photo, exam_room, seat_number）
- 不会删除已有字段或数据

---

## 管理员功能

- **账号**: `test@example.com` / `123456`
- **查看所有报名**: `GET /api/applications` (需 admin 角色)
- **新闻管理**: POST/PUT/DELETE `/api/news` (需 admin 角色)

---

## 下一步计划

1. ⬜ 配置 Sealos 外网访问（域名/端口）
2. ⬜ 正式域名绑定
3. ⬜ 考试时间确定后更新准考证时间表
4. ⬜ 成绩发布功能
5. ⬜ 管理后台界面（目前仅 API）
6. ⬜ 清理测试数据

---

*文档创建时间：2026-06-21*
*最后更新：2026-06-25*
*当前阶段：核心功能全部完成，准备外网部署*
