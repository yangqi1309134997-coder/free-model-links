# 二次开发与自定义指南

## 修改品牌信息

进入管理员后台 **"网站设置"** 页面（`/admin/config`），可直接修改：

- 网站标题
- 站点关键词与描述
- 网站版权信息（默认为“新疆幻城网安科技有限责任公司”）
- 网站 Logo / Favicon（支持上传图片，自动转 Base64 嵌入）
- 备案信息与备案管理地址

修改后立即生效，无需重启服务。

---

## 新增广告位

### 1. 前台页面添加广告容器

在 `web/templates/index/*.html` 的合适位置插入：

```html
<!-- 广告位：自定义位置 -->
<div id="ad-custom-slot" style="max-width:1200px;margin:0 auto;padding:24px;"></div>
```

### 2. 加载广告脚本

在同页面 `<script>` 中调用：

```javascript
loadAds('custom_slot', 'ad-custom-slot');
```

### 3. 后台创建广告

进入 **"系统设置" → "广告管理"**，新建广告时选择广告位 `custom_slot` 即可。

---

## 修改主题样式

前端主题样式位于 `web/static/index/css/apple-design.css`，主要 CSS 变量：

```css
:root {
  --accent: #007aff;           /* 主色调 */
  --accent-light: #007aff15;   /* 浅色背景 */
  --text-primary: #1d1d1f;     /* 主文字 */
  --text-secondary: #86868b;   /* 次要文字 */
  --bg-card: #ffffff;          /* 卡片背景 */
  --border: rgba(0,0,0,0.08);  /* 边框色 */
  --radius-lg: 20px;           /* 大圆角 */
  --radius-md: 12px;           /* 中圆角 */
}
```

暗色模式变量在 `[data-theme="dark"]` 下定义，修改即可同步影响明暗两套主题。

> 注意：修改静态资源后需重新编译才能生效（`go build`）。

---

## 添加新的系统配置字段

若需在 `sys_config` 表中新增配置项，按以下步骤操作：

### 1. 修改数据模型

编辑 `internal/dal/model/sys_config.gen.go`，在 `SysConfig` 结构体新增字段：

```go
NewField string `gorm:"column:new_field;type:varchar(100)" json:"new_field"`
```

### 2. 更新 gorm-gen 查询代码

编辑 `internal/dal/query/sys_config.gen.go`，在 `sysConfig` 结构体、构造函数、`updateTableName`、`fillFieldMap` 中同步添加新字段。

### 3. 更新 API 协议

在 `api/v1/config.go` 的 `ConfigResp` 和 `ConfigUpdateReq` 中新增对应字段。

### 4. 更新业务逻辑

在 `internal/service/config/config.go` 和 `internal/service/config/update.go` 中处理新字段的读取与保存。

### 5. 更新前端模板

在 `web/templates/config/conf_web.html` 添加表单控件，在 `web/templates/index/index.html` 中使用该配置。

### 6. 重新编译

```bash
go build -o ./bin/server ./cmd/server
```

---

## 添加新的后台管理页面

### 1. 创建 HTML 模板

在 `web/templates/adminext/` 下新建页面，例如 `admin_logs.html`。建议复用现有后台的 Bootstrap 4 + jQuery 技术栈，保持风格一致。

### 2. 注册路由

编辑 `internal/server/http.go`，在 `render` 路由组中添加：

```go
render.GET("logs", func(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "admin_logs.html", nil)
})
```

### 3. 添加菜单入口

编辑 `web/templates/admin/admin_index.html`，在侧边栏菜单中增加入口链接：

```html
<a class="nav-link" href="/admin/logs">
    <i class="mdi mdi-file-document"></i> <span>日志管理</span>
</a>
```

### 4. 实现后端 API（如需）

在 `internal/handler/adminext/handler.go` 中新增 Handler 方法，并在 `strictAuthRouter` 中注册对应路由。

---

## 健康检查扩展

当前健康检查仅检测 HTTP 可达性。如需增加更复杂的探测逻辑（如检查 API 响应格式、认证状态等），编辑 `internal/healthcheck/checker.go`：

```go
func (c *Checker) checkSite(site model.StSite) HealthResult {
    // 自定义探测逻辑
    resp, err := http.Get(site.URL)
    // ...
}
```

健康检查结果会写入 `st_site.health_status`（JSON 格式），前端详情页自动展示。

---

## 数据库迁移

本项目使用 GORM 的 `AutoMigrate`，首次启动会自动建表。后续若手动修改表结构，可编写迁移脚本：

```go
// cmd/migrate/main.go
package main

import (
    "github.com/ch3nnn/webstack-go/internal/dal/model"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    db, _ := gorm.Open(sqlite.Open("storage/webstack-go.db"), &gorm.Config{})
    db.AutoMigrate(&model.SysConfig{})
    // 执行 ALTER TABLE 等操作
}
```

> 生产环境建议使用更严谨的迁移方案（如 [golang-migrate/migrate](https://github.com/golang-migrate/migrate)）。
