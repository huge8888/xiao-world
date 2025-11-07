# Xiao-World: 多平台内容发布系统

基于 [xiaohongshu-mcp](https://github.com/xpzouying/xiaohongshu-mcp) 扩展的多平台内容发布系统。

## 新功能

本项目在原有小红书 MCP 功能基础上，增加了**多平台内容发布**能力：

- ✅ 从小红书获取内容
- ✅ 自动翻译为英文（支持 Google Translate API）
- ✅ 发布到多个国际社交媒体平台：
  - **Twitter/X** - 支持文本和图片
  - **TikTok** - 支持视频
  - **Facebook** - 支持文本、图片和视频
  - **YouTube** - 支持视频
- ✅ 定时发布功能
- ✅ 内容自动适配各平台要求

## 新增 MCP 工具

在原有 12 个小红书工具基础上，新增了 8 个多平台发布工具：

### 13. `publish_to_twitter`
将小红书笔记内容发布到 Twitter/X（自动翻译为英文）

**参数：**
- `feed_id` - 小红书笔记ID
- `xsec_token` - 访问令牌

### 14. `publish_to_tiktok`
将小红书视频内容发布到 TikTok（自动翻译为英文）

### 15. `publish_to_facebook`
将小红书内容发布到 Facebook（自动翻译为英文）

### 16. `publish_to_youtube`
将小红书视频内容发布到 YouTube（自动翻译为英文）

### 17. `publish_to_all_platforms`
同时发布到多个平台

**参数：**
- `feed_id` - 小红书笔记ID
- `xsec_token` - 访问令牌
- `platforms` - 平台列表（可选）：`["twitter", "tiktok", "facebook", "youtube"]`

### 18. `schedule_publish`
创建定时发布任务

**参数：**
- `feed_id` - 小红书笔记ID
- `xsec_token` - 访问令牌
- `platforms` - 平台列表
- `scheduled_at` - 定时时间（格式：`2006-01-02 15:04:05`）

### 19. `list_scheduled_jobs`
查看所有定时发布任务及其状态

### 20. `cancel_scheduled_job`
取消指定的定时发布任务

**参数：**
- `job_id` - 任务ID

## 配置

### 环境变量配置

```bash
# Twitter/X 配置
export TWITTER_ENABLED=true
export TWITTER_BEARER_TOKEN="your_bearer_token"

# TikTok 配置
export TIKTOK_ENABLED=true
export TIKTOK_ACCESS_TOKEN="your_access_token"
export TIKTOK_CLIENT_KEY="your_client_key"
export TIKTOK_CLIENT_SECRET="your_client_secret"

# Facebook 配置
export FACEBOOK_ENABLED=true
export FACEBOOK_ACCESS_TOKEN="your_page_access_token"
export FACEBOOK_PAGE_ID="your_page_id"

# YouTube 配置
export YOUTUBE_ENABLED=true
export YOUTUBE_ACCESS_TOKEN="your_access_token"
export YOUTUBE_REFRESH_TOKEN="your_refresh_token"
export YOUTUBE_CLIENT_ID="your_client_id"
export YOUTUBE_CLIENT_SECRET="your_client_secret"

# Google Translate API（可选，不设置则使用免费服务）
export GOOGLE_TRANSLATE_API_KEY="your_api_key"
```

### 配置文件（可选）

也可以使用 JSON 配置文件：

```bash
./xiao-world -config /path/to/config
```

配置文件示例：

**twitter.json:**
```json
{
  "enabled": true,
  "bearer_token": "your_bearer_token"
}
```

**tiktok.json:**
```json
{
  "enabled": true,
  "access_token": "your_access_token",
  "client_key": "your_client_key",
  "client_secret": "your_client_secret"
}
```

**facebook.json:**
```json
{
  "enabled": true,
  "access_token": "your_page_access_token",
  "page_id": "your_page_id"
}
```

**youtube.json:**
```json
{
  "enabled": true,
  "access_token": "your_access_token",
  "refresh_token": "your_refresh_token",
  "client_id": "your_client_id",
  "client_secret": "your_client_secret"
}
```

## 使用示例

### 1. 发布到单个平台

```
用户：请将这篇小红书笔记发布到 Twitter
{
  "feed_id": "abc123",
  "xsec_token": "xyz..."
}
```

### 2. 发布到所有平台

```
用户：请将这篇笔记同时发布到所有平台
{
  "feed_id": "abc123",
  "xsec_token": "xyz...",
  "platforms": ["twitter", "facebook"]
}
```

### 3. 定时发布

```
用户：请在明天下午3点将这篇笔记发布到 Twitter 和 Facebook
{
  "feed_id": "abc123",
  "xsec_token": "xyz...",
  "platforms": ["twitter", "facebook"],
  "scheduled_at": "2025-11-08 15:00:00"
}
```

## 平台特性

### Twitter/X
- 文本限制：280 字符（Premium 用户 4000 字符）
- 支持最多 4 张图片
- 自动下载并上传图片

### TikTok
- 仅支持视频内容
- 标题限制：2200 字符
- 自动下载并上传视频

### Facebook
- 无文本限制
- 支持多张图片和视频
- 自动处理图片上传

### YouTube
- 仅支持视频内容
- 标题限制：100 字符
- 描述限制：5000 字符
- 支持标签

## 工作原理

1. **内容获取**：从小红书获取笔记详情（文本、图片、视频）
2. **内容翻译**：使用 Google Translate API 将中文翻译为英文
3. **内容适配**：根据各平台限制调整内容格式
4. **媒体处理**：下载图片/视频并上传到目标平台
5. **发布**：调用各平台 API 发布内容

## 技术架构

```
xiaohongshu-mcp (基础层)
    ↓
pkg/translator (翻译层)
    ↓
pkg/processor (内容处理层)
    ↓
pkg/publishers/* (发布层)
    ├── twitter
    ├── tiktok
    ├── facebook
    └── youtube
    ↓
pkg/scheduler (调度层)
```

## Docker 部署

```bash
docker run -d \
  -e TWITTER_ENABLED=true \
  -e TWITTER_BEARER_TOKEN="..." \
  -e FACEBOOK_ENABLED=true \
  -e FACEBOOK_ACCESS_TOKEN="..." \
  -e FACEBOOK_PAGE_ID="..." \
  -p 18060:18060 \
  xiao-world:latest
```

## API 凭证获取

### Twitter/X
1. 访问 [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard)
2. 创建应用获取 Bearer Token

### TikTok
1. 访问 [TikTok for Developers](https://developers.tiktok.com/)
2. 注册应用并申请 Content Posting API 权限

### Facebook
1. 访问 [Facebook Developers](https://developers.facebook.com/)
2. 创建应用并获取 Page Access Token

### YouTube
1. 访问 [Google Cloud Console](https://console.cloud.google.com/)
2. 启用 YouTube Data API v3
3. 创建 OAuth 2.0 凭据

## 原始功能

本项目保留了 xiaohongshu-mcp 的所有原始功能，包括：

1. ✅ 登录和检查登录状态
2. ✅ 获取登录二维码
3. ✅ 删除 cookies
4. ✅ 发布小红书图文内容
5. ✅ 获取 Feeds 列表
6. ✅ 搜索内容
7. ✅ 获取笔记详情
8. ✅ 获取用户主页
9. ✅ 发表评论
10. ✅ 发布视频
11. ✅ 点赞笔记
12. ✅ 收藏笔记

详细文档请参考原项目：[xiaohongshu-mcp](https://github.com/xpzouying/xiaohongshu-mcp)

## 许可证

MIT License

## 致谢

感谢原项目 [xiaohongshu-mcp](https://github.com/xpzouying/xiaohongshu-mcp) 提供的基础能力。
