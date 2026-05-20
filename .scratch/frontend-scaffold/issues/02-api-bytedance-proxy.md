# Scaffold-02：`/api/bytedance` 托管单模型代理

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

在 **Go `backend`** 实现 NextChat 兼容的 **`/api/bytedance/*`** 透明转发（对照 `NextChat/app/api/bytedance.ts`）：

- 使用 **`ARK_API_KEY`、`ARK_MODEL_ID`、`ARK_BASE_URL`** 请求火山方舟；请求体中的 `model` 可被服务端规范为配置的 model id（托管单模型）
- 支持流式响应（SSE/chunked）透传
- 处理 `OPTIONS` 预检
- **不经 SafetyGate**（危机拦截仍在 `/v1/sessions/stream`）
- 若配置了 `CODE`，与 NextChat 一致的访问口令校验（对照 NextChat `auth`）

不实现其它 `/api/{provider}`；未实现路径可返回 501。

可用集成测试或脚本对 `POST /api/bytedance/v1/chat/completions`（或 NextChat 实际子路径）发一条最小对话验证。

## 验收标准

- [ ] 配置有效 `ARK_*` 时，对代理路径发流式请求能收到上游增量内容
- [ ] 缺 key 时返回明确 4xx/5xx，不 panic
- [ ] 流式与非流式（若 NextChat 会发）至少覆盖主路径
- [ ] 未配置 `CODE` 时无需口令即可请求（与 #01 一致）
- [ ] `go test ./...` 全绿

## 阻塞于

- [01-api-config-e2e.md](./01-api-config-e2e.md)

## 覆盖的用户故事

Scaffold 技术验证（非 PRD 洪峰用户故事）
