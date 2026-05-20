# Frontend Scaffold（NextChat Web 接入）

**上游：** [CONTEXT.md](../../CONTEXT.md)（术语 **Frontend Scaffold**）、[ADR-0001](../../docs/adr/0001-frontend-backend-split.md)（部分范围已收窄，见 PRD）、[frontend/MIGRATION.md](../../frontend/MIGRATION.md)

**目标：** 将 NextChat **浏览器 Web** 壳迁入 `frontend/`，经 Go 提供 **托管单模型** 聊天（临时路径 `/api/bytedance`，`ARK_*`），完成最小 E2E。不含 Tauri/多端、不含多厂商代理、Scaffold 期 `/api/*` 不过 SafetyGate。

**产品化**（洪峰首屏、关怀 UX）属后续 **Frontend P1**，不在本 epic。

**状态（2026-05-20）**：01–03、05–06 已实现；04（Mask 打磨）可选。请本地验证流式聊天。

## Issue 顺序

| # | 文件 | 说明 |
|---|------|------|
| 01 | [01-api-config-e2e.md](./issues/01-api-config-e2e.md) | `GET /api/config` 竖切 |
| 02 | [02-api-bytedance-proxy.md](./issues/02-api-bytedance-proxy.md) | `/api/bytedance` 透明代理 |
| 03 | [03-nextchat-chat-core-e2e.md](./issues/03-nextchat-chat-core-e2e.md) | 迁聊天核心 + 流式一问一答 |
| 04 | [04-nextchat-masks-locales.md](./issues/04-nextchat-masks-locales.md) | Mask / 多语言（不阻塞脚手架） |
| 05 | [05-remove-web-mvp.md](./issues/05-remove-web-mvp.md) | 删除 `web/` 与静态托管 |
| 06 | [06-docs-env-alignment.md](./issues/06-docs-env-alignment.md) | `.env.example`、MIGRATION、ADR 备注 |
