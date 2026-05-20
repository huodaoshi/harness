# Web 壳（family-wellness-ai #04）

静态 PWA 友好页面，由 `backend` 同一进程托管。

## 本地运行

```powershell
cd d:\harness\backend
docker compose up -d
# 复制 .env.example → .env，设置 JWT_SECRET（见 backend/config/app/local.yaml）
go run ./cmd/server
```

浏览器打开 http://localhost:8080/

**依赖**：Mongo + **Redis**（鉴权与 stream 限流）。仅内存 wellness 数据时可设 `USE_MEMORY_STORE=true`，Redis 仍须可用。

## 鉴权（P1-07）

所有 `/v1/profile`、`/v1/sessions/*` 请求由 `js/user.js` 自动附带：

- 已登录：`Authorization: Bearer <token>`（`localStorage` 键 `fwa_access_token`）
- 游客：`X-Anon-ID: <uuid>`（键 `fwa_anon_id`，自旧版 `fwa_user_id` 迁移）

不再发送 query `user_id` 或 `X-User-Id`。401/429 会在页面顶部错误条显示中文提示。

## 验收要点

1. 勾选免责声明后才能进入聊天
2. 「我现在很难受」→ `mode=distress`；「先聊聊」→ `mode=normal`
3. 发送消息后可见流式打字，直至 `done`
4. 输入危机关键词（如「不想活了」）应出现危机卡片并禁用输入
5. 移动端单列布局，输入区 sticky 在底部
6. 「编辑关系档案」可增删重要他人；保存后下一条聊天应带上档案上下文
7. 去掉鉴权头或伪造 `X-Anon-ID` → 401；短时间大量 stream → 429

## 文件

- `index.html` — 欢迎页 + 聊天页
- `css/app.css` — 移动端样式
- `js/user.js` — 游客 UUID / Bearer 与请求头
- `js/sse.js` — POST SSE 解析
- `js/app.js` — 会话状态与 UI
- `js/profile.js` — 关系档案
