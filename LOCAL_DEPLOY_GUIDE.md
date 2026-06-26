# 济南新航实验外国语学校网站 — 本地服务器部署完全指南

> 本指南面向对 Linux 运维不太熟悉的人员编写，每一步都有详细说明。
> 只要按顺序执行，即可在学校本地 Ubuntu 服务器上完成部署。

---

## 目录

1. [服务器硬件准备与 RAID5 配置](#1-服务器硬件准备与-raid5-配置)
2. [Ubuntu 系统初始配置](#2-ubuntu-系统初始配置)
3. [更换阿里云镜像源（加速下载）](#3-更换阿里云镜像源加速下载)
4. [安装必要软件](#4-安装必要软件)
5. [将代码传到服务器](#5-将代码传到服务器)
6. [配置数据库](#6-配置数据库)
7. [配置 Redis](#7-配置-redis)
8. [配置后端环境变量](#8-配置后端环境变量)
9. [构建前端](#9-构建前端)
10. [编译并运行后端](#10-编译并运行后端)
11. [导入初始数据](#11-导入初始数据)
12. [设为开机自启服务](#12-设为开机自启服务)
13. [配置 Nginx 反向代理](#13-配置-nginx-反向代理)
14. [内网暴露端口与防火墙设置](#14-内网暴露端口与防火墙设置)
15. [HTTPS 配置（可选）](#15-https-配置可选)
16. [日常维护与故障排查](#16-日常维护与故障排查)

---

## 1. 服务器硬件准备与 RAID5 配置

### 什么是 RAID5？
RAID5 使用 3 块或以上硬盘，将数据分散存储并生成校验信息。如果其中 1 块硬盘坏了，数据不会丢失，更换新盘后可自动恢复。这样能保证学校数据的安全性。

### 硬件要求
- **CPU**：4核以上
- **内存**：8GB 以上（推荐 16GB）
- **硬盘**：至少 3 块相同容量的硬盘（用于 RAID5）
- **网络**：接入学校内网，有固定 IP 地址

### 安装 Ubuntu 时配置 RAID5（推荐方式）

如果服务器有硬件 RAID 控制器（如 Dell PERC、HP SmartArray），建议在 BIOS/UEFI 中直接配置 RAID5，再安装 Ubuntu。这种方式最稳定。

### 软件 RAID5 配置（如果没有硬件 RAID）

安装完 Ubuntu 后，假设你有 3 块额外的数据盘 `/dev/sdb`、`/dev/sdc`、`/dev/sdd`：

```bash
# 第一步：安装 RAID 管理工具
sudo apt update
sudo apt install -y mdadm

# 第二步：查看可用硬盘（确认盘符）
lsblk

# 你会看到类似这样的输出：
# sda       系统盘（已安装 Ubuntu）
# sdb       数据盘1
# sdc       数据盘2
# sdd       数据盘3

# 第三步：创建 RAID5 阵列
# 注意：这会清除 sdb/sdc/sdd 上的所有数据！
sudo mdadm --create /dev/md0 --level=5 --raid-devices=3 /dev/sdb /dev/sdc /dev/sdd

# 系统会提示确认，输入 y 回车

# 第四步：等待 RAID 同步完成（可能需要几小时）
# 查看同步进度：
cat /proc/mdstat

# 第五步：格式化 RAID 分区
sudo mkfs.ext4 /dev/md0

# 第六步：创建挂载目录并挂载
sudo mkdir -p /data
sudo mount /dev/md0 /data

# 第七步：设置开机自动挂载
sudo mdadm --detail --scan | sudo tee -a /etc/mdadm/mdadm.conf
sudo update-initramfs -u

# 在 /etc/fstab 添加一行
echo '/dev/md0  /data  ext4  defaults  0  2' | sudo tee -a /etc/fstab

# 第八步：验证
df -h /data
# 应该能看到 /data 已挂载且容量正确
```

> **后续所有网站文件都放在 `/data/` 目录下**，这样数据就存储在 RAID5 阵列上了。

---

## 2. Ubuntu 系统初始配置

安装好 Ubuntu Server 后，首先进行基础配置：

```bash
# 用你安装时创建的用户登录服务器
# 如果是远程操作，用 SSH 连接：
# ssh 你的用户名@服务器IP地址

# 更新系统（第一次可能需要几分钟）
sudo apt update
sudo apt upgrade -y

# 安装常用工具
sudo apt install -y curl wget git unzip tar vim nano htop net-tools

# 设置时区为中国
sudo timedatectl set-timezone Asia/Shanghai

# 确认时间正确
date
```

---

## 3. 更换阿里云镜像源（加速下载）

Ubuntu 默认从国外服务器下载软件，速度很慢。更换为阿里云镜像源后，下载速度会快很多。

```bash
# 第一步：备份原始源文件（以防出错可恢复）
sudo cp /etc/apt/sources.list /etc/apt/sources.list.backup

# 第二步：查看你的 Ubuntu 版本代号
lsb_release -cs
# 会输出类似 "noble"（24.04）或 "jammy"（22.04）

# 第三步：替换为阿里云源
# 如果你是 Ubuntu 24.04 (noble)：
sudo tee /etc/apt/sources.list << 'EOF'
deb http://mirrors.aliyun.com/ubuntu/ noble main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ noble-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ noble-backports main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ noble-security main restricted universe multiverse
EOF

# 如果你是 Ubuntu 22.04 (jammy)，把上面所有 "noble" 替换为 "jammy"

# 第四步：更新软件源索引
sudo apt update

# 如果看到类似 "Hit:1 http://mirrors.aliyun.com/ubuntu noble InRelease" 说明成功了
```

> **提示**：如果 `sudo apt update` 报错，说明版本代号写错了。用 `lsb_release -cs` 确认后重新修改。

---

## 4. 安装必要软件

### 4.1 安装 Go 语言（后端运行环境）

```bash
# 下载 Go（使用国内镜像加速）
wget https://mirrors.aliyun.com/golang/go1.22.5.linux-amd64.tar.gz

# 如果上面的链接下载慢，也可以用：
# wget https://studygolang.com/dl/golang/go1.22.5.linux-amd64.tar.gz

# 解压到 /usr/local
sudo tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
echo 'export GOPROXY=https://goproxy.cn,direct' | sudo tee -a /etc/profile.d/go.sh
source /etc/profile.d/go.sh

# 验证安装
go version
# 应该输出：go version go1.22.5 linux/amd64

# 清理下载文件
rm go1.22.5.linux-amd64.tar.gz
```

> `GOPROXY=https://goproxy.cn` 设置让 Go 下载依赖时走国内镜像，速度更快。

### 4.2 安装 Node.js（仅用于构建前端，构建完成后不再需要）

```bash
# 使用 NodeSource 安装 Node.js 20
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# 验证
node --version   # 应该显示 v20.x.x
npm --version    # 应该显示 10.x.x

# 设置 npm 国内镜像（加速下载前端依赖）
npm config set registry https://registry.npmmirror.com
```

### 4.3 安装 PostgreSQL 数据库

```bash
# 安装
sudo apt install -y postgresql postgresql-contrib

# 启动并设置开机自启
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 检查状态（应该显示 active (running)）
sudo systemctl status postgresql
```

### 4.4 安装 Redis 缓存

```bash
# 安装
sudo apt install -y redis-server

# 启动并设置开机自启
sudo systemctl start redis-server
sudo systemctl enable redis-server

# 验证（应该返回 PONG）
redis-cli ping
```

### 4.5 安装 Nginx（反向代理，让网站用 80 端口访问）

```bash
sudo apt install -y nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## 5. 将代码传到服务器

### 方式一：使用 U 盘 / 移动硬盘

```bash
# 将部署压缩包 xinhang-school-deploy.zip 拷贝到 U 盘
# 将 U 盘插入服务器

# 查看 U 盘设备名
lsblk
# 通常是 /dev/sde1 或类似

# 挂载 U 盘
sudo mkdir -p /mnt/usb
sudo mount /dev/sde1 /mnt/usb

# 复制文件到数据目录
cp /mnt/usb/xinhang-school-deploy.zip /data/

# 卸载 U 盘
sudo umount /mnt/usb
```

### 方式二：通过网络传输（SCP / SFTP）

在你自己的电脑上执行（不是在服务器上）：

```bash
# Windows 用户可以用 WinSCP 图形界面传文件
# 或者在 PowerShell / 终端中执行：
scp xinhang-school-deploy.zip 用户名@服务器IP:/data/
```

### 解压代码

```bash
cd /data
unzip xinhang-school-deploy.zip

# 解压后你会得到一个 xinhang-app 文件夹
ls xinhang-app/
# 应该看到：backend/  src/  public/  package.json  vite.config.js 等文件
```

---

## 6. 配置数据库

```bash
# 进入 PostgreSQL 命令行
sudo -u postgres psql
```

在 PostgreSQL 提示符中输入以下命令（每行输入后按回车）：

```sql
-- 创建数据库
CREATE DATABASE xinhang;

-- 创建用户（请把 'YourSecurePassword123' 改为你自己的密码）
-- 密码要求：至少8位，包含字母和数字
CREATE USER xinhang_user WITH PASSWORD 'YourSecurePassword123';

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE xinhang TO xinhang_user;

-- 连接到 xinhang 数据库并授予 schema 权限
\c xinhang
GRANT ALL ON SCHEMA public TO xinhang_user;

-- 退出
\q
```

> **重要**：记住你设置的密码，后面配置环境变量时需要用到！

### 验证数据库连接

```bash
# 测试能否用新用户连接
psql -h localhost -U xinhang_user -d xinhang -c "SELECT 1;"
# 输入刚才设置的密码，如果显示一个表格说明成功
```

如果连接失败，可能需要修改 PostgreSQL 认证配置：

```bash
# 编辑认证配置文件
sudo nano /etc/postgresql/16/main/pg_hba.conf
# 注：16 是版本号，你的可能是 14 或 15，用 ls /etc/postgresql/ 查看

# 找到这一行（大约在文件末尾）：
# local   all   all   peer
# 改为：
# local   all   all   md5

# 同时确保有这一行：
# host    all   all   127.0.0.1/32   md5

# 保存退出（Ctrl+O 保存，Ctrl+X 退出）

# 重启 PostgreSQL
sudo systemctl restart postgresql

# 再次测试连接
psql -h localhost -U xinhang_user -d xinhang -c "SELECT 1;"
```

---

## 7. 配置 Redis

Redis 默认配置即可使用，无需额外修改。如果想设置密码（可选）：

```bash
# 编辑 Redis 配置
sudo nano /etc/redis/redis.conf

# 找到 # requirepass foobared，去掉注释并改为你的密码：
# requirepass YourRedisPassword

# 保存退出后重启
sudo systemctl restart redis-server
```

> 如果设置了 Redis 密码，后面的 `.env` 文件中也要填写对应密码。

---

## 8. 配置后端环境变量

```bash
cd /data/xinhang-app/backend

# 从模板创建配置文件
cp .env.example .env

# 编辑配置文件
nano .env
```

修改 `.env` 文件内容如下（带 ← 的地方需要你修改）：

```env
# 服务端口（后端监听的端口）
PORT=8080

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=xinhang_user
DB_PASSWORD=YourSecurePassword123        ← 改为第6步设置的数据库密码
DB_NAME=xinhang
DB_MAX_CONNS=50
DB_IDLE_CONNS=10

# Redis 配置
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=                           ← 如果第7步设置了Redis密码就填这里，否则留空
REDIS_DB=0

# JWT 密钥（用于用户登录鉴权，必须修改为随机字符串）
# 执行这个命令生成一个：openssl rand -hex 32
JWT_SECRET=请替换为随机字符串             ← 必须修改！

# 邮件配置（用于发送注册验证码）
SMTP_HOST=smtp.qq.com                    ← 根据你用的邮箱修改
SMTP_PORT=465
SMTP_USER=your-email@qq.com             ← 改为学校的发件邮箱
SMTP_PASSWORD=your_authorization_code   ← 改为邮箱授权码（不是登录密码！）
SMTP_FROM=your-email@qq.com             ← 同上面的邮箱地址

# 准考证验证密码（教职工扫码验证时输入）
VERIFY_PIN=xinhang2026                   ← 可以改为学校自定义的密码
```

### 生成 JWT 密钥

```bash
# 执行这个命令，会输出一串随机字符
openssl rand -hex 32
# 复制输出的字符串，粘贴到 .env 中的 JWT_SECRET= 后面
```

### 获取邮箱授权码（以 QQ 邮箱为例）

1. 登录 QQ 邮箱网页版
2. 点击"设置" → "账户"
3. 找到"POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务"
4. 开启"IMAP/SMTP服务"
5. 按提示验证后获取授权码
6. 将授权码填入 `SMTP_PASSWORD`

保存配置文件：按 `Ctrl+O` 保存，`Ctrl+X` 退出 nano。

---

## 9. 构建前端

```bash
# 回到项目根目录
cd /data/xinhang-app

# 安装前端依赖（首次需要下载，可能需要几分钟）
npm install

# 如果下载慢，确认已设置了国内镜像：
# npm config set registry https://registry.npmmirror.com
# 然后重新执行 npm install

# 构建前端（输出到 backend/dist/ 目录）
npm run build

# 构建成功后检查
ls backend/dist/
# 应该看到 index.html 和 assets/ 文件夹
```

> 构建成功后，前端代码已变成纯静态文件，Node.js 就不再需要了。

---

## 10. 编译并运行后端

```bash
cd /data/xinhang-app/backend

# 下载 Go 依赖（首次需要时间，已配置国内镜像会快很多）
go mod download

# 编译为可执行文件
go build -o xinhang-backend .

# 先手动运行测试（看是否有报错）
./xinhang-backend
```

如果看到以下输出，说明启动成功：
```
Database connected (maxOpen=50, maxIdle=10)
Database migration completed
Redis connected at localhost:6379
Serving frontend from ./dist
Server starting on :8080
```

按 `Ctrl+C` 停止程序（后面我们会设置为系统服务自动运行）。

### 测试访问

在服务器上执行：
```bash
curl http://localhost:8080/
# 如果返回 HTML 内容，说明网站正常运行
```

---

## 11. 导入初始数据

### 导入新闻数据

新闻系统包含两部分：**数据库记录**（文章标题、正文等）和**新闻配图**（文章中引用的图片文件）。两者都需要导入，否则新闻页面的图片会显示不出来。

#### 第一步：导入新闻文章数据（SQL）

```bash
# 确保后端已经启动过一次（自动创建了数据表）
# 然后导入新闻数据
psql -h localhost -U xinhang_user -d xinhang -f /data/xinhang-app/backend/database/seed_news.sql
# 输入数据库密码
```

导入成功后，数据库中 `news` 表会包含 300+ 篇新闻文章。

#### 第二步：导入新闻配图（重要！）

新闻文章的正文 HTML 中引用了 `/uploads/migration/images/xxx.jpg` 路径下的图片。这些图片文件单独打包在 `xinhang-news-images.tar.gz` 中（约 1.3GB，包含 2000+ 张图片）。

**如果不导入这些图片，新闻文章中的配图将全部显示为空白/裂图。**

```bash
# 在服务器上创建目标目录
mkdir -p /data/xinhang-app/backend/uploads/migration

# 将图片压缩包传到服务器（与代码包分开传输，因为有1.3GB）
# 方式1：U盘拷贝
cp /mnt/usb/xinhang-news-images.tar.gz /data/

# 方式2：SCP传输（在你的电脑上执行）
# scp xinhang-news-images.tar.gz 用户名@服务器IP:/data/

# 解压到正确位置
cd /data/xinhang-app/backend/uploads/migration
tar -xzf /data/xinhang-news-images.tar.gz

# 验证：应该看到 images/ 文件夹，里面有大量 .jpg/.png 文件
ls images/ | head -5
ls images/ | wc -l
# 应该显示 2021 个文件
```

解压后的目录结构应该是：
```
/data/xinhang-app/backend/uploads/
└── migration/
    └── images/
        ├── 000e3d22767a45b0a9b192f06f4cc132.jpg
        ├── 008e3d78e7db4b2ba994dba3add93442.jpg
        ├── ... （共 2021 个图片文件）
```

#### 验证新闻图片是否正常

```bash
# 确保后端正在运行
sudo systemctl status xinhang

# 测试图片是否可以访问（用任意一个图片文件名测试）
curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/uploads/migration/images/000e3d22767a45b0a9b192f06f4cc132.jpg
# 应该返回 200，如果返回 404 说明路径不对
```

如果返回 404，检查：
1. 目录路径是否正确：`ls /data/xinhang-app/backend/uploads/migration/images/`
2. 后端工作目录是否是 `/data/xinhang-app/backend/`（在 systemd 服务文件中的 `WorkingDirectory`）

> **提示**：如果学校以后不需要旧新闻文章，可以跳过图片导入。新发布的新闻文章不会用到这些旧图片。

---

### 部署包清单

给学校的完整交付物应包含：

| 文件 | 大小 | 说明 |
|------|------|------|
| `xinhang-school-deploy.zip` | ~33MB | 网站源代码 + 构建产物 + 配置模板 + 新闻SQL |
| `xinhang-news-images.tar.gz` | ~1.3GB | 新闻文章配图（2021张） |
| `LOCAL_DEPLOY_GUIDE.md` | - | 本部署指南（已包含在zip中） |

> 由于新闻图片有 1.3GB，建议用 U 盘或移动硬盘传输，不建议通过网络传。

### 创建管理员账号

1. 先在网站上正常注册一个账号（需要邮箱验证码）
2. 然后在数据库中将该账号提升为管理员：

```bash
sudo -u postgres psql -d xinhang

-- 将指定邮箱的用户设为管理员
UPDATE users SET role = 'admin' WHERE email = '你注册时用的邮箱';

-- 验证
SELECT email, role FROM users WHERE role = 'admin';

\q
```

---

## 12. 设为开机自启服务

让网站在服务器重启后自动运行：

```bash
# 创建 systemd 服务文件
sudo tee /etc/systemd/system/xinhang.service << 'EOF'
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

# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start xinhang

# 设置开机自启
sudo systemctl enable xinhang

# 检查状态
sudo systemctl status xinhang
# 应该显示 active (running)
```

### 常用管理命令

```bash
# 查看服务状态
sudo systemctl status xinhang

# 停止服务
sudo systemctl stop xinhang

# 重启服务（修改配置后需要重启）
sudo systemctl restart xinhang

# 查看实时日志
sudo journalctl -u xinhang -f

# 查看最近100行日志
sudo journalctl -u xinhang -n 100
```

---

## 13. 配置 Nginx 反向代理

Nginx 让用户可以通过 80 端口（不用加 :8080）访问网站：

```bash
# 创建 Nginx 配置文件
sudo tee /etc/nginx/sites-available/xinhang << 'EOF'
server {
    listen 80;
    server_name _;

    # 文件上传大小限制（学生照片）
    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持（如果后续需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
EOF

# 启用配置（创建软链接）
sudo ln -sf /etc/nginx/sites-available/xinhang /etc/nginx/sites-enabled/xinhang

# 删除默认配置（避免冲突）
sudo rm -f /etc/nginx/sites-enabled/default

# 测试配置是否正确
sudo nginx -t
# 应该显示：syntax is ok / test is successful

# 重新加载 Nginx
sudo systemctl reload nginx
```

现在用浏览器访问 `http://服务器IP` 就能看到网站了（不需要加 :8080）。

---

## 14. 内网暴露端口与防火墙设置

### 检查并开放防火墙端口

```bash
# 查看防火墙状态
sudo ufw status

# 如果防火墙是 active 状态，需要开放端口：
sudo ufw allow 22/tcp      # SSH（远程管理用）
sudo ufw allow 80/tcp      # HTTP（网站访问）
sudo ufw allow 443/tcp     # HTTPS（如果以后配了SSL）

# 启用防火墙（如果还没启用）
sudo ufw enable

# 再次确认
sudo ufw status
# 应该显示 22、80、443 端口都是 ALLOW
```

### 确认网站可从内网其他电脑访问

```bash
# 查看服务器的内网 IP
ip addr show
# 找到类似 192.168.x.x 或 10.x.x.x 的地址

# 或者：
hostname -I
```

在学校内网的其他电脑上，打开浏览器访问：
```
http://服务器IP地址
```

例如：`http://192.168.1.100`

### 如果内网其他电脑访问不了

1. **检查服务器防火墙**：确保上面的 ufw 规则已添加
2. **检查交换机/路由器**：确保服务器和访问电脑在同一网段，或路由已配置
3. **检查 Nginx 是否监听所有地址**：
   ```bash
   sudo netstat -tlnp | grep :80
   # 应该显示 0.0.0.0:80（监听所有地址）
   ```
4. **ping 测试**：在其他电脑上执行 `ping 服务器IP`，如果不通说明网络物理连接有问题

### 绑定域名（可选）

如果学校有内网 DNS 服务器，可以配置一个域名（如 `school.local`）指向这台服务器的 IP。这样用户可以通过 `http://school.local` 访问。

如果没有内网 DNS，可以在每台需要访问的电脑上修改 hosts 文件：
- Windows: `C:\Windows\System32\drivers\etc\hosts`
- Mac/Linux: `/etc/hosts`

添加一行：
```
192.168.1.100  school.xinhang.local
```

---

## 15. HTTPS 配置（可选）

### 内网自签名证书

如果只在内网使用，可以用自签名证书：

```bash
# 生成证书（有效期10年）
sudo mkdir -p /etc/nginx/ssl
sudo openssl req -x509 -nodes -days 3650 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/xinhang.key \
  -out /etc/nginx/ssl/xinhang.crt \
  -subj "/CN=school.xinhang.local"

# 修改 Nginx 配置
sudo tee /etc/nginx/sites-available/xinhang << 'EOF'
server {
    listen 80;
    server_name _;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name _;

    ssl_certificate /etc/nginx/ssl/xinhang.crt;
    ssl_certificate_key /etc/nginx/ssl/xinhang.key;

    client_max_body_size 10M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

sudo nginx -t && sudo systemctl reload nginx
```

> 注：自签名证书会让浏览器显示"不安全"警告，点击"高级"→"继续访问"即可。

---

## 16. 日常维护与故障排查

### 备份数据库（建议每天执行）

```bash
# 创建备份目录
sudo mkdir -p /data/backups

# 手动备份
pg_dump -h localhost -U xinhang_user -d xinhang > /data/backups/xinhang_$(date +%Y%m%d).sql

# 设置每天凌晨3点自动备份
sudo tee /etc/cron.d/xinhang-backup << 'EOF'
0 3 * * * root pg_dump -h localhost -U xinhang_user -d xinhang > /data/backups/xinhang_$(date +\%Y\%m\%d).sql && find /data/backups/ -name "*.sql" -mtime +30 -delete
EOF
```

> 自动备份会保留最近30天的数据，超过30天的自动删除。

### 备份上传的照片

```bash
# 照片存储在这里
ls /data/xinhang-app/backend/uploads/photos/

# 备份照片
cp -r /data/xinhang-app/backend/uploads /data/backups/uploads_$(date +%Y%m%d)
```

### 常见故障排查

| 问题 | 检查方法 | 解决办法 |
|------|----------|----------|
| 网站打不开 | `sudo systemctl status xinhang` | 如果 failed，看日志 `journalctl -u xinhang -n 50` |
| 数据库连接失败 | `sudo systemctl status postgresql` | 重启：`sudo systemctl restart postgresql` |
| Redis 连接失败 | `redis-cli ping` | 重启：`sudo systemctl restart redis-server` |
| 注册收不到验证码 | 查看后端日志中的邮件错误 | 检查 .env 中 SMTP 配置 |
| 照片上传失败 | `ls -la /data/xinhang-app/backend/uploads/` | 创建目录并设置权限 |
| 磁盘空间不足 | `df -h` | 清理旧备份或日志 |
| RAID 磁盘故障 | `cat /proc/mdstat` | 看下面的 RAID 维护说明 |

### RAID5 维护

```bash
# 查看 RAID 状态
cat /proc/mdstat
sudo mdadm --detail /dev/md0

# 正常状态应显示：State : active
# 如果显示 degraded，说明有一块盘坏了

# 更换坏盘步骤：
# 1. 确认哪块盘坏了
sudo mdadm --detail /dev/md0 | grep -i fail

# 2. 移除坏盘（假设是 /dev/sdc）
sudo mdadm /dev/md0 --remove /dev/sdc

# 3. 物理更换硬盘后，添加新盘
sudo mdadm /dev/md0 --add /dev/sdc

# 4. 等待重建完成
watch cat /proc/mdstat
```

### 更新网站代码

当开发团队提供新版本的代码时：

```bash
# 1. 停止服务
sudo systemctl stop xinhang

# 2. 备份当前代码
cp -r /data/xinhang-app /data/xinhang-app-backup-$(date +%Y%m%d)

# 3. 将新代码解压到 /data/xinhang-app（覆盖旧文件）
cd /data
unzip -o new-version.zip

# 4. 重新构建前端
cd /data/xinhang-app
npm install
npm run build

# 5. 重新编译后端
cd backend
go mod download
go build -o xinhang-backend .

# 6. 重启服务
sudo systemctl start xinhang

# 7. 检查是否正常
sudo systemctl status xinhang
curl http://localhost:8080/
```

### 查看实时访问日志

```bash
# 后端日志
sudo journalctl -u xinhang -f

# Nginx 访问日志
sudo tail -f /var/log/nginx/access.log

# Nginx 错误日志
sudo tail -f /var/log/nginx/error.log
```

---

## 快速检查清单

部署完成后，逐项确认：

- [ ] 服务器能 ping 通
- [ ] `sudo systemctl status postgresql` 显示 active
- [ ] `sudo systemctl status redis-server` 显示 active
- [ ] `sudo systemctl status xinhang` 显示 active
- [ ] `sudo systemctl status nginx` 显示 active
- [ ] 在服务器上 `curl http://localhost` 返回 HTML
- [ ] 在内网其他电脑浏览器打开 `http://服务器IP` 能看到网站
- [ ] 能正常注册账号并收到验证码邮件
- [ ] 能正常登录
- [ ] 能正常提交报名表
- [ ] 新闻页面能显示文章

---

## 联系方式

如部署过程中遇到问题，请联系开发团队并提供：
1. 具体报错信息（截图或复制文字）
2. 执行到哪一步出的问题
3. `sudo journalctl -u xinhang -n 50` 的输出

---

*文档版本：3.0*
*更新日期：2026-06-26*
