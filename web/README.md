# Web 壳（family-wellness-ai #04）

静态 PWA 友好页面，由 `backend` 同一进程托管。

## 本地运行

```powershell
cd d:\harness\backend
$env:USE_MEMORY_STORE = "true"
go run ./cmd/server
```

浏览器打开 http://localhost:8080/

## 验收要点

1. 勾选免责声明后才能进入聊天
2. 「我现在很难受」→ `mode=distress`；「先聊聊」→ `mode=normal`
3. 发送消息后可见流式打字，直至 `done`
4. 输入危机关键词（如「不想活了」）应出现危机卡片并禁用输入
5. 移动端单列布局，输入区 sticky 在底部
6. 「编辑关系档案」可增删重要他人；保存后下一条聊天应带上档案上下文

## 文件

- `index.html` — 欢迎页 + 聊天页
- `css/app.css` — 移动端样式
- `js/sse.js` — POST SSE 解析
- `js/app.js` — 会话状态与 UI
