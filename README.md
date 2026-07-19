# API 服务商导航平台

基于 [webstack-go](https://github.com/ch3nnn/webstack-go) 深度定制开发的 **AI / LLM API 服务商导航平台**，专为聚合、展示和评价大模型 API 中转服务商而设计。由 **新疆幻城网安科技有限责任公司** 维护运营。

> **算力支持**：本项目由 [新疆幻城网安公益大模型 API 中转站](https://api.iamhc.cn) 提供算力与接口测试支持。成品网站演示请看这里https://link.hcnsec.cn

> 本项目在 webstack-go v2（Go + Gin + GORM + SQLite）基础上进行了大量功能增强和前端重构，保留了单二进制部署、零依赖（无需 MySQL/Redis）的轻量化优势。

---

## 文档

- [部署指南](docs/DEPLOYMENT.md) — 二进制、源码、Docker、Nginx、systemd 完整部署说明
- [API 接口概览](docs/API_OVERVIEW.md) — 公开接口与管理员接口速查
- [二次开发与自定义指南](docs/CUSTOMIZATION.md) — 主题修改、新增配置字段、新增广告位、扩展后台页面

---

## 核心特性

### 前端展示
- **Apple 风格现代化界面**：半透明导航栏、精致卡片、圆角设计、明暗主题切换
- **网站详情页**：支持多 Tab（概览 / 定价 / 评论），展示运行状态、支持模型、支付方式、截图等
- **智能搜索**：首页实时过滤，支持按分类、评分、评论数、 newest 排序
- **SEO 全面优化**：Open Graph、Twitter Card、JSON-LD 结构化数据、Canonical URL、Robots 自定义
- **友情链接**：后台可配置，首页页脚自动展示
- **广告系统**：多广告位（首页顶部/底部、详情页、关于页、投稿页），支持 HTML/图片/文本三种类型，后台管理投放时段与开关

### 用户互动
- **匿名评分**：访客无需登录即可为服务商打分，自动计算并缓存平均分
- **匿名评论**：支持昵称 + 内容投稿，提交后进入审核状态，通过后前台展示
- **网站投稿**：访客可提交新的 API 服务商，管理员后台审核通过后自动创建网站记录

### 管理员后台
- **仪表盘**：SSE 实时数据推送，展示站点统计概览
- **网站管理**：
  - 分类管理：支持一级/二级分类，图标中文化
  - 网站列表：单条新增（自动抓取标题、Logo、描述）+ 批量新增 + 编辑 + 删除
  - 高级字段：富文本介绍、截图、支持模型、定价信息、健康状态、标签、支付方式
  - 精选 / 推荐标记
- **字典管理**：模型字典、标签字典、支付方式字典，支持排序与启用开关
- **评论审核**：后台列表查看、状态变更（通过/拒绝）、管理员备注
- **评分管理**：查看所有评分记录，支持修改分数、删除后自动重算站点平均分
- **用户管理**：查看管理员账号列表、状态管理
- **投稿审核**：查看用户提交的站点信息，一键审核通过并自动入库
- **系统设置**：
  - 网站基础信息（标题、关键词、描述、版权、Logo、Favicon、备案）
  - 关于页面（Markdown 编辑器，支持启用/禁用）
  - 网站语言、默认排序、匿名评分/评论/投稿开关
  - 健康检查配置（间隔、失败阈值、手动触发）
  - 首页布局（每行卡片数、精选/推荐区块显示开关）
  - SEO 设置（作者、OG 图片、Twitter Card、Canonical、Robots）
  - 友情链接（JSON 配置）
  - 广告管理（完整 CRUD，按广告位筛选）

### 技术特性
- **单二进制部署**：编译后仅需一个可执行文件 + SQLite 数据库文件
- **模板嵌入**：使用 `embed.FS` 将前端模板与静态资源打包进二进制
- **JWT 认证**：管理后台采用 Token 认证，支持 Header 与 Cookie 双通道
- **健康检查**：内置定时探测，检测站点在线状态、可用率、平均延迟
- **gorm-gen**：类型安全的查询构造，自动生成 DAO 层代码

---

## 运行环境

- Go 1.22+
- SQLite（内置，无需额外安装）

---

## 快速开始

### 1. 二进制运行

从 [Releases](https://github.com/ch3nnn/webstack-go/releases) 下载预编译二进制：

```bash
./webstack-go -conf config/prod.yml
```

> **macOS 提示**：若提示“无法检查是否包含恶意软件”，请前往“系统设置 > 隐私与安全性 > 安全性”点击“仍然允许”。

### 2. 源码编译运行

```bash
# 1. 安装依赖
go mod tidy

# 2. 编译
go build -o ./bin/server ./cmd/server

# 3. 首次启动（自动生成 SQLite 数据库与表结构）
./bin/server -conf=config/prod.yml
```

默认监听 `0.0.0.0:8000`，浏览器访问 http://127.0.0.1:8000。

### 3. Docker 运行

```bash
docker run -i -t --restart always -p 8000:8000 \
  --name api-nav \
  -v ./data/storage:/data/app/storage \
  ch3nnn/webstack-go:latest
```

或使用 Docker Compose：

```yaml
services:
  api-nav:
    stdin_open: true
    tty: true
    restart: always
    ports:
      - 8000:8000
    container_name: api-nav
    image: ch3nnn/webstack-go:latest
    volumes:
      - ./data/storage:/data/app/storage
```

---

## 默认账号

首次启动后，系统会自动创建默认管理员账号：

| 字段   | 值       |
| ------ | -------- |
| 用户名 | `admin`  |
| 密码   | `admin`  |

> **重要**：生产环境请务必修改默认密码！登录后台后前往“个人中心”修改。

---

## 项目结构

```
webstack-go/
├── api/v1/               # API 协议定义（Request/Response 结构体）
├── cmd/server/           # 主程序入口
├── config/               # 配置文件（dev/prod/test）
├── internal/
│   ├── dal/model/        # GORM 数据模型
│   ├── dal/query/        # gorm-gen 自动生成的查询代码
│   ├── dal/repository/   # 数据仓库与初始化逻辑
│   ├── handler/          # HTTP 请求处理器
│   │   ├── admin/        # 管理员后台 API
│   │   ├── adminext/     # 扩展后台 API（评论/评分/字典/广告等）
│   │   ├── publicapi/    # 公开 API（站点列表/详情/评分/评论/广告等）
│   │   └── user/         # 用户认证相关
│   ├── healthcheck/      # 站点健康检查定时任务
│   ├── i18n/             # 国际化文案
│   ├── middleware/       # Gin 中间件（JWT/权限/跨域等）
│   ├── server/           # HTTP 服务器与路由注册
│   └── service/          # 业务逻辑层
├── pkg/                  # 公共工具包（JWT/日志/加密等）
├── storage/              # SQLite 数据库文件存放目录
├── web/
│   ├── static/           # 静态资源（CSS/JS/图片）
│   └── templates/        # HTML 模板（Go template）
│       ├── index/        # 前台页面（首页/详情/关于/投稿）
│       ├── admin/        # 后台管理页面
│       ├── adminext/     # 扩展后台页面（设置/评论/评分等）
│       └── config/       # 基础配置页面
└── README.md
```

---

## 主要配置文件

`config/prod.yml` 示例：

```yaml
app:
  name: api-nav
  version: 2.0.0
  env: prod

http:
  addr: 0.0.0.0:8000

jwt:
  secret: your-secret-key-here   # 请修改！用于 JWT 签名
  expiration: 86400              # Token 有效期（秒）

log:
  level: info
  filename: storage/logs/app.log
```

---

## 二次开发指南

### 修改模板

本项目使用 Go 1.16+ 的 `embed.FS` 将模板打包进二进制，**修改模板后需要重新编译**才能生效：

```bash
go build -o ./bin/server ./cmd/server
```

### 修改数据模型

若新增或修改了 `internal/dal/model/` 下的模型文件，需要更新 gorm-gen 生成的查询代码：

```bash
go run ./cmd/gen
```

或手动同步 `internal/dal/query/*.gen.go` 中的字段映射。

### 添加新的广告位

1. 在 `web/templates/index/*.html` 中插入广告容器：
   ```html
   <div id="ad-your-slot" style="max-width:1200px;margin:0 auto;padding:24px;"></div>
   ```
2. 在同页面的 `<script>` 中调用：
   ```javascript
   loadAds('your_slot', 'ad-your-slot');
   ```
3. 后台“广告管理”中新建广告，选择广告位 `your_slot` 即可。

---

## 技术栈

| 层级       | 技术选型                                          |
| ---------- | ------------------------------------------------- |
| 后端框架   | Go 1.22 + Gin                                     |
| ORM        | GORM + gorm-gen                                   |
| 数据库     | SQLite（轻量级，零配置）                          |
| 前端样式   | 纯手写 CSS（Apple Design 风格）                   |
| 模板引擎   | Go html/template + embed.FS                       |
| 后台 UI    | Light-Year Admin（Bootstrap 4 中文后台模板）      |
| 认证       | JWT（github.com/golang-jwt/jwt/v5）               |
| 部署       | 单二进制 / Docker                                 |

---

## 维护方

**新疆幻城网安科技有限责任公司**

- 官网：[https://hcnsec.cn](https://hcnsec.cn)
- 本项目作为 API 服务商导航平台，旨在帮助开发者高效发现和对接大模型 API 资源。

---

## 许可证

本项目基于原 [webstack-go](https://github.com/ch3nnn/webstack-go) 开源项目进行深度定制，遵循原项目许可证。

---

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=ch3nnn/webstack-go&type=Date)](https://star-history.com/#ch3nnn/webstack-go&Date)
