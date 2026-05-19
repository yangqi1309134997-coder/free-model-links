# API 接口概览

## 基础信息

- **Base URL**: `http://localhost:8000/api`
- **Content-Type**: `application/json`
- **认证方式**: JWT Token（Header: `Token: <jwt>` 或 Cookie: `_login_token_`）

---

## 公开接口（无需认证）

### 站点信息

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/site-config` | 获取网站基础配置（语言、版权、开关等） |
| GET | `/sites` | 站点列表（支持分页、分类、搜索、标签、模型筛选、排序） |
| GET | `/sites/:id` | 站点详情（含模型、定价、健康状态、标签、支付方式） |
| GET | `/models` | 支持的模型字典列表 |
| GET | `/tags` | 标签字典列表 |

### 互动

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/sites/:id/rating` | 匿名评分（需后台开启允许匿名评分） |
| POST | `/sites/:id/review` | 匿名评论（需后台开启允许匿名评论，审核后展示） |
| GET | `/sites/:id/reviews` | 获取某站点的已审核评论列表 |
| POST | `/submissions` | 提交新站点（需后台开启允许投稿） |

### 广告

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ads?slot=xxx` | 按广告位获取当前生效的广告列表 |

---

## 管理员接口（需 JWT 认证）

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 管理员登录（用户名 + MD5密码） |
| GET | `/admin/info` | 获取当前登录管理员信息 |
| POST | `/admin/logout` | 退出登录 |
| PATCH | `/admin/modify_password` | 修改密码 |

### 网站管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/category` | 分类列表（树形结构） |
| POST | `/admin/category` | 创建分类 |
| PUT | `/admin/category/:id` | 更新分类 |
| DELETE | `/admin/category/:id` | 删除分类 |
| GET | `/admin/site` | 网站列表 |
| POST | `/admin/site` | 新增网站 |
| PUT | `/admin/site/:id` | 更新网站 |
| DELETE | `/admin/site/:id` | 删除网站 |
| POST | `/admin/site/batch` | 批量新增网站 |

### 扩展管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/reviews` | 评论审核列表 |
| PUT | `/admin/reviews/:id` | 审核评论（通过/拒绝） |
| GET | `/admin/ratings` | 评分记录列表 |
| PUT | `/admin/ratings/:id` | 修改评分 |
| DELETE | `/admin/ratings/:id` | 删除评分 |
| GET | `/admin/submissions` | 用户投稿列表 |
| PUT | `/admin/submissions/:id` | 审核投稿（通过自动创建网站） |
| GET | `/admin/users` | 管理员账号列表 |
| PUT | `/admin/users/:id` | 更新用户状态 |

### 字典管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/models` | 模型字典列表 |
| POST | `/admin/models` | 创建模型 |
| PUT | `/admin/models/:id` | 更新模型 |
| DELETE | `/admin/models/:id` | 删除模型 |
| GET | `/admin/tags` | 标签字典列表 |
| POST | `/admin/tags` | 创建标签 |
| PUT | `/admin/tags/:id` | 更新标签 |
| DELETE | `/admin/tags/:id` | 删除标签 |
| GET | `/admin/payment-methods` | 支付方式字典列表 |
| POST | `/admin/payment-methods` | 创建支付方式 |
| PUT | `/admin/payment-methods/:id` | 更新支付方式 |
| DELETE | `/admin/payment-methods/:id` | 删除支付方式 |

### 系统设置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/config` | 获取网站基础配置 |
| PUT | `/admin/config` | 更新网站基础配置 |
| GET | `/admin/site-settings` | 获取站点开关设置 |
| PUT | `/admin/site-settings` | 更新站点开关设置 |
| GET | `/admin/health-config` | 获取健康检查配置 |
| PUT | `/admin/health-config` | 更新健康检查配置 |
| POST | `/admin/health-trigger` | 手动触发健康检查 |
| GET | `/admin/homepage-config` | 获取首页布局配置 |
| PUT | `/admin/homepage-config` | 更新首页布局配置 |
| GET | `/admin/seo-config` | 获取 SEO 配置 |
| PUT | `/admin/seo-config` | 更新 SEO 配置 |
| PUT | `/admin/friend-links` | 更新友情链接 |

### 广告管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/admin/ads` | 广告列表（支持按广告位筛选） |
| POST | `/admin/ads` | 创建广告 |
| PUT | `/admin/ads/:id` | 更新广告 |
| DELETE | `/admin/ads/:id` | 删除广告 |

---

## 请求示例

### 管理员登录

```bash
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"21232f297a57a5a743894a0e4a801fc3"}'
```

响应：

```json
{
  "code": 200,
  "msg": "ok",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### 获取站点列表

```bash
curl "http://localhost:8000/api/sites?page=1&page_size=20&category_id=1&sort_by=rating&sort_order=desc"
```

### 提交匿名评分

```bash
curl -X POST http://localhost:8000/api/sites/1/rating \
  -H "Content-Type: application/json" \
  -d '{"score":5}'
```

### 更新网站配置（需 Token）

```bash
curl -X PUT http://localhost:8000/api/admin/config \
  -H "Content-Type: application/json" \
  -H "Token: YOUR_JWT_TOKEN" \
  -d '{"site_title":"API 导航","site_copyright":"新疆幻城网安科技有限责任公司"}'
```
