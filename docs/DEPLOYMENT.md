# 部署指南

## 环境要求

- Go 1.22 或更高版本（源码编译时）
- 或直接使用预编译的二进制文件 / Docker 镜像
- 操作系统：Linux / macOS / Windows

---

## 方式一：二进制部署（推荐生产环境）

### 1. 下载二进制

从 [Releases](https://github.com/ch3nnn/webstack-go/releases) 下载对应平台的二进制文件。

### 2. 准备配置文件

创建 `config/prod.yml`：

```yaml
app:
  name: api-nav
  version: 2.0.0
  env: prod

http:
  addr: 0.0.0.0:8000

jwt:
  secret: YOUR_STRONG_SECRET_HERE   # 务必修改！
  expiration: 86400

log:
  level: info
  filename: storage/logs/app.log
```

> **安全提示**：`jwt.secret` 必须使用强随机字符串，否则存在 token 被伪造的风险。

### 3. 启动服务

```bash
chmod +x webstack-go
./webstack-go -conf config/prod.yml
```

首次启动会自动创建 SQLite 数据库文件（默认位于 `storage/webstack-go.db`）及表结构。

### 4. 使用 systemd 守护（Linux）

创建 `/etc/systemd/system/api-nav.service`：

```ini
[Unit]
Description=API Nav Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/api-nav
ExecStart=/opt/api-nav/webstack-go -conf config/prod.yml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动：

```bash
sudo systemctl daemon-reload
sudo systemctl enable api-nav
sudo systemctl start api-nav
```

---

## 方式二：源码编译部署

```bash
git clone https://github.com/ch3nnn/webstack-go.git
cd webstack-go

# 安装依赖
go mod tidy

# 编译
go build -o ./bin/server ./cmd/server

# 启动
./bin/server -conf=config/prod.yml
```

---

## 方式三：Docker 部署

### 使用官方镜像

```bash
docker run -d \
  --name api-nav \
  --restart always \
  -p 8000:8000 \
  -v $(pwd)/data/storage:/data/app/storage \
  ch3nnn/webstack-go:latest
```

### 使用 Docker Compose

```yaml
version: "3.8"

services:
  api-nav:
    image: ch3nnn/webstack-go:latest
    container_name: api-nav
    restart: always
    ports:
      - "8000:8000"
    volumes:
      - ./data/storage:/data/app/storage
    environment:
      - TZ=Asia/Shanghai
```

启动：

```bash
docker-compose up -d
```

### 本地构建 Docker 镜像

```bash
make docker
```

---

## Nginx 反向代理（推荐）

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 数据备份

SQLite 数据库为单文件，直接备份即可：

```bash
# 热备份（SQLite 支持读取时备份）
cp storage/webstack-go.db backup/webstack-go-$(date +%Y%m%d).db
```

建议配合 `cron` 每日自动备份数据库文件。

---

## 常见问题

### 1. 端口被占用

```bash
# 查找占用 8000 端口的进程
lsof -i :8000
# 或修改配置文件中的 http.addr 为其他端口
```

### 2. 权限不足

确保运行用户对 `storage/` 目录有读写权限：

```bash
chmod -R 755 storage/
```

### 3. 模板修改后不生效

本项目使用 Go `embed.FS` 打包模板，**修改模板后必须重新编译**。

### 4. 如何修改默认管理员密码

登录后台后，点击右上角头像 → "修改密码"，或直接在数据库中修改 `sys_user` 表的 `password` 字段（MD5 哈希值）。
