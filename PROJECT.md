# 山东新航实验国际学校 — 全栈网站项目

## 项目概述

将学校原有的静态 HTML 官网迁移为 **Vue3 + Go + PostgreSQL** 全栈应用，部署在 **Sealos** 云平台上，支持用户注册和学生报名功能。

---

## 当前进度总览

| 阶段 | 状态 | 说明 |
|------|------|------|
| 静态 HTML 官网 | ✅ 已完成 | 原始 6 个页面，已部署在上码（现已弃用） |
| Vue3 前端改造 | ✅ 已完成 | 6 个页面迁移 + 注册/登录/报名 3 个新页面 |
| Sealos 服务器 | ✅ 已创建 | Go 开发环境已建好，已装 Go 1.22 + Node 22 |
| Sealos PostgreSQL | ✅ 已连接 | xinhang 数据库已建，users + applications 表已自动创建 |
| Sealos Redis | ✅ 已连接 | 用于缓存和 API 限流 |
| 异步消息队列 | ✅ 已完成 | Redis Streams 替代 Kafka（Kafka 在 Sealos 有 Operator bug） |
| Go 后端 | ✅ 已完成 | Gin + GORM + Redis + Redis Streams |
| 前后端联调 | ✅ 已完成 | 注册/登录/报名/管理员列表 4 个 API 全部测试通过 |
| 生产部署 | ✅ 运行中 | Go serve 前端 dist（方案 A），8080 端口运行中 |

---

## 第一阶段目标

> 把 HTML 网页搬上去，看上去像回事，并能支持**注册**和**报名**功能。

具体来说：
1. 任何人可以访问网站（公开）
2. 用户可以注册账号（邮箱 + 手机 + 密码）
3. 用户可以提交报名申请（学生信息、家长信息、年级等）
4. 后台可以查看所有报名记录（管理员）
5. 预计报名量级：**~1000 条**

---

## 技术架构

```
用户浏览器
    │
    ▼
┌─────────────┐     /api/*       ┌─────────────┐      SQL       ┌──────────────┐
│  Vue3 前端   │ ──────────────► │  Go 后端     │ ────────────► │ PostgreSQL   │
│  (Nginx)    │ ◄────────────── │  (Gin)      │ ◄──────────── │  数据库       │
└─────────────┘     JSON        └─────────────┘               └──────────────┘
     Sealos 容器                    Sealos 容器                   Sealos 数据库
```

---

## 前端项目详情（已完成）

### 技术栈
- **Vue 3** + **Vite 8** + **Vue Router 4**
- 中英双语切换（composable: `useLanguage`）
- 样式：原始 CSS 全量迁移，无框架依赖

### 项目结构

```
xinhang-frontend/
├── index.html                        # Vite 入口
├── vite.config.js                    # Vite 配置，/api 代理到 localhost:8080
├── package.json
├── public/
│   └── images/                       # 21 张图片资源（jpg/png/svg）
├── src/
│   ├── main.js                       # Vue 挂载入口
│   ├── App.vue                       # 根组件（Header + RouterView + Footer）
│   ├── assets/
│   │   └── styles.css                # 全局样式（原 CSS + 表单样式）
│   ├── router/
│   │   └── index.js                  # 9 条路由配置
│   ├── composables/
│   │   └── useLanguage.js            # 中英双语切换逻辑
│   ├── components/
│   │   ├── SiteHeader.vue            # 公共头部导航（含注册按钮）
│   │   └── SiteFooter.vue            # 公共底部
│   └── views/
│       ├── HomeView.vue              # 首页（原 index.html）
│       ├── AboutView.vue             # 关于新航
│       ├── AcademicsView.vue         # 学术课程
│       ├── AdmissionView.vue         # 招生信息
│       ├── CampusView.vue            # 校园展示
│       ├── StudentLifeView.vue       # 学生生活
│       ├── RegisterView.vue          # ✨ 注册页（新增）
│       ├── LoginView.vue             # ✨ 登录页（新增）
│       └── ApplyView.vue             # ✨ 报名申请页（新增）
```

### 路由表

| 路径 | 页面 | 说明 |
|------|------|------|
| `/` | HomeView | 首页 |
| `/about` | AboutView | 关于新航 |
| `/academics` | AcademicsView | 学术课程 |
| `/admission` | AdmissionView | 招生信息 |
| `/campus` | CampusView | 校园展示 |
| `/student-life` | StudentLifeView | 学生生活 |
| `/register` | RegisterView | 用户注册 |
| `/login` | LoginView | 用户登录 |
| `/apply` | ApplyView | 提交报名申请 |

### 前端已预留的 API 接口调用

| 接口 | 方法 | 来源页面 | 用途 |
|------|------|----------|------|
| `/api/register` | POST | RegisterView | 用户注册 |
| `/api/login` | POST | LoginView | 用户登录，返回 JWT token |
| `/api/apply` | POST | ApplyView | 提交报名申请 |

#### 注册请求体
```json
{
  "name": "张三",
  "email": "zhangsan@example.com",
  "phone": "13800138000",
  "password": "123456"
}
```

#### 登录请求体
```json
{
  "email": "zhangsan@example.com",
  "password": "123456"
}
```
登录成功响应需包含 `{ "token": "jwt-token-here" }`，前端会存到 `localStorage('xinhang-token')`。

#### 报名请求体
```json
{
  "studentName": "李四",
  "birthDate": "2012-05-15",
  "gender": "male",
  "grade": 7,
  "parentName": "李大强",
  "phone": "13900139000",
  "email": "parent@example.com",
  "currentSchool": "济南某某小学",
  "notes": "希望了解国际方向"
}
```
报名请求会附带 `Authorization: Bearer <token>` 头（如果用户已登录）。

---

## Go 后端（待搭建）

### 需要实现的内容

1. **框架**：推荐 Gin
2. **数据库驱动**：`pgx` 或 `gorm`
3. **API 路由**：
   - `POST /api/register` — 注册
   - `POST /api/login` — 登录（返回 JWT）
   - `POST /api/apply` — 提交报名
   - `GET /api/applications` — 获取报名列表（管理员，需鉴权）
4. **数据库表**：
   - `users` — id, name, email, phone, password_hash, role, created_at
   - `applications` — id, user_id(可选), student_name, birth_date, gender, grade, parent_name, phone, email, current_school, notes, created_at
5. **鉴权**：JWT，存在 Authorization header
6. **密码**：bcrypt 加密存储

### 数据库连接信息

创建数据库后，在 Sealos 数据库详情页获取：
- Host: `xinhang-db-postgresql.ns-xxx.svc.cluster.local`（内网地址）
- Port: `5432`
- Username: `postgres`
- Password: （查看详情页）
- Database: `postgres`（可创建新的 `xinhang` 数据库）

连接字符串格式：
```
postgresql://postgres:PASSWORD@HOST:5432/xinhang?sslmode=disable
```

---

## Sealos 部署配置

### 数据库（已创建）
- 类型：PostgreSQL 16.4
- CPU: 0.5 Core / 内存: 512MB~1G
- 磁盘: 3Gi
- max_connections: 50
- time_zone: Asia/Shanghai

### Go 后端服务器（已创建）
- 模板：Go 1.22.5
- 在此服务器上：
  1. 克隆 Git 仓库
  2. 搭建 Go 后端
  3. 构建并运行
  4. 打包前端 dist 用 Nginx 或 Go 直接 serve

### 生产部署方案（后续）
- **方案 A（简单）**：Go 后端同时 serve 前端静态文件，一个容器搞定
- **方案 B（标准）**：Nginx 容器 serve 前端 + 反代 /api 到 Go 容器

推荐第一阶段用**方案 A**，最简单。

---

## 已完成待办

1. ✅ SSH 连接 Sealos 服务器（通过 Cursor Remote）
2. ✅ 从 GitHub 拉取代码到服务器
3. ✅ 在服务器上安装 Node.js 22 + Go 1.22
4. ✅ 搭建 Go 后端项目（Gin + GORM + Redis + Redis Streams）
5. ✅ 创建数据库 xinhang + 自动建表（users, applications）
6. ✅ 实现 4 个 API 接口（注册/登录/报名/管理列表）
7. ✅ 前端打包 `npm run build` → dist/
8. ✅ Go 后端 serve 前端 dist + API（方案 A 一体部署）
9. ✅ 端到端测试（注册 → 登录 → 报名 → 后台查看）
10. ✅ Redis 连接（缓存 + 限流）
11. ✅ 管理员账号创建（test@example.com / 123456）

## 下一步待办

1. ✅ ~~Kafka~~ → 已用 Redis Streams 替代（Sealos Kafka Operator 有 advertised listener bug，broker 无法启动）
2. ⬜ 配置 Sealos 外网访问（域名/端口）
3. ⬜ 升级 PostgreSQL 资源（4C/4GB/2实例）
4. ⬜ 正式域名绑定
5. ⬜ 清理测试数据（压测/E2E 产生的测试用户和报名）

## 管理员账号

- 邮箱：`test@example.com`
- 密码：`123456`
- 登录后用 token 访问 `GET /api/applications` 可查看所有报名

## 后端架构

```
backend/
├── main.go              # 入口：整合所有组件 + gzip 压缩
├── config/config.go     # 统一配置（DB/Redis）
├── database/database.go # 连接池（maxOpen=100, maxIdle=20）
├── cache/redis.go       # Redis 缓存 + 限流
├── queue/kafka.go       # Redis Streams 异步消息队列（已从 Kafka 迁移）
├── middleware/
│   ├── auth.go          # JWT 鉴权 + 管理员权限
│   └── ratelimit.go     # Redis 限流中间件
├── handlers/
│   ├── auth.go          # 注册(10/min) + 登录(20/min)
│   └── application.go   # 报名(5/min, Redis Stream异步) + 列表
└── models/              # User + Application 数据模型
```

## 环境变量 (.env)

详见 `backend/.env.example`

---

## 原始项目历史

- 原始静态 HTML 网站位于 `c:\Users\14863\Desktop\xinhang-webb-inspired-demo`
- 曾部署在上码平台（upma.site），现已弃用
- 图片资源全部在 `public/images/` 目录下
- 支持中英文双语切换

---

*文档创建时间：2026-06-21*
*最后更新：2026-06-22*
*当前阶段：后端全部完成（Redis + Redis Streams），等待外网开放和域名绑定*
