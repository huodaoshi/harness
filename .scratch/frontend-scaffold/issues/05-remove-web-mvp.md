# Scaffold-05：删除 `web/` MVP 静态壳

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

在 #03 验收通过后：

- 删除仓库根目录 **`web/`**
- 移除 `backend` 中 **`RegisterWebStatic`** 及对 `configpaths.WebRoot()` 的挂载
- 更新相关测试（如 `web_static_test.go`）
- 本地开发文档改为以 **`frontend` :3000** 为 Web 入口

Wellness 能力仍通过 `/v1/*` API 提供；仅去掉旧静态 MVP 页面。

## 验收标准

- [ ] `web/` 目录不存在
- [ ] `go test ./...` 全绿
- [ ] `go run ./cmd/server` 不再托管原 `web/index.html`；`:8080/` 不再提供旧 MVP（可 404 或重定向说明，以实现为准）
- [ ] `frontend` dev 仍可完成聊天 E2E

## 阻塞于

- [03-nextchat-chat-core-e2e.md](./03-nextchat-chat-core-e2e.md)

## 覆盖的用户故事

替换 family-wellness-ai #04 的长期入口（新入口为 NextChat 壳，产品化后续）
