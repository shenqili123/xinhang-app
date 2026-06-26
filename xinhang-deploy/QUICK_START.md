# 新航学校网站 — 现场部署快速操作手册

> 此手册供现场部署人员使用，按顺序执行即可。
> 遇到问题请联系远程技术支持，并提供屏幕截图。

---

## 前提条件

- 服务器已安装 Ubuntu（24.04 或 22.04）
- 你能登录服务器（有用户名和密码）
- U 盘中包含 `xinhang-deploy/` 文件夹

---

## 第一步：挂载 U 盘，复制文件

```bash
# 插入 U 盘后，查看设备名
lsblk

# 挂载（sdb1 根据实际情况修改）
sudo mkdir -p /mnt/usb
sudo mount /dev/sdb1 /mnt/usb

# 复制部署包到本地
sudo cp -r /mnt/usb/xinhang-deploy /tmp/
cd /tmp/xinhang-deploy

# 卸载 U 盘
sudo umount /mnt/usb
```

如果 U 盘是 NTFS 格式无法挂载：
```bash
sudo apt install ntfs-3g
sudo mount -t ntfs-3g /dev/sdb1 /mnt/usb
```

---

## 第二步：运行环境检查

```bash
sudo bash scripts/00-check-environment.sh
```

确认输出中没有严重错误，然后继续。

---

## 第三步：安装软件

```bash
sudo bash scripts/01-install-software.sh
```

等待完成，确认 PostgreSQL、Redis、Nginx 都显示 [✓]。

---

## 第四步：配置数据库

```bash
sudo bash scripts/02-setup-database.sh
```

按提示输入数据库密码（自己决定一个，**务必记住**）。

---

## 第五步：部署应用文件

```bash
sudo bash scripts/03-deploy-app.sh
```

---

## 第六步：配置环境变量

```bash
sudo bash scripts/04-configure-env.sh
```

按提示依次输入：
- 数据库密码（第四步设置的）
- 邮箱信息（没准备好可输入 skip 跳过）
- 验证PIN（直接回车用默认值）

---

## 第七步：启动服务

```bash
sudo bash scripts/05-start-services.sh
```

脚本结束后会显示服务器 IP 地址。
在其他电脑浏览器输入 `http://显示的IP` 看是否能打开网站。

---

## 第八步：导入新闻数据

```bash
sudo bash scripts/06-import-data.sh
```

按提示输入数据库密码，等待图片解压完成（约1~2分钟）。

---

## 第九步：验收检查

```bash
sudo bash scripts/07-verify-all.sh
```

所有项目显示 [✓] 即为部署成功。

---

## 部署完成后

### 创建管理员账号

1. 在网站上正常注册一个账号
2. 执行以下命令提升为管理员：

```bash
sudo -u postgres psql -d xinhang -c "UPDATE users SET role='admin' WHERE email='你的注册邮箱';"
```

### 常用命令速查

| 操作 | 命令 |
|------|------|
| 查看网站状态 | `sudo systemctl status xinhang` |
| 重启网站 | `sudo systemctl restart xinhang` |
| 查看日志 | `sudo journalctl -u xinhang -f` |
| 修改配置 | `sudo nano /data/xinhang-app/backend/.env` |
| 服务器IP | `hostname -I` |

---

## 紧急故障处理

**网站打不开？**
```bash
sudo systemctl restart xinhang
sudo systemctl restart nginx
```

**数据库问题？**
```bash
sudo systemctl restart postgresql
```

**查看错误信息：**
```bash
sudo journalctl -u xinhang -n 50
```

---

*操作过程中如有任何疑问，请立即联系远程技术支持。*
