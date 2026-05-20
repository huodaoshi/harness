# PRD：Frontend Scaffold（NextChat Web 接入）

**Status:** ready-for-agent  
**特性：** frontend-scaffold  
**版本：** Scaffold v1  
**最后更新：** 2026-05-20  
**上游：** grill-with-docs 共识（见 [CONTEXT.md](../../CONTEXT.md)）

---

## 问题陈述

仓库需将 NextChat 聊天 UI 与 Go `backend` 分离，但不宜一次性复制多厂商代理与全部产品能力。需要先搭 **可演示的 Web 聊天架子**，密钥与模型由服务端托管，前端不配置多模型。

## 解决方案

- 新建/充实 `frontend/`：从 `NextChat/` 按目录迁移（不复制 `app/api/`、`src-tauri/`）。
- `backend` 新增最小 NextChat 兼容 API：`GET /api/config`、`POST /api/bytedance`（及 OPTIONS），使用 **`ARK_*`** 转发火山方舟；未配置 `CODE` 时 `needCode=false`。
- 本地开发：`frontend` `:3000`，`/api` rewrite 到 `backend` `:8080`。
- 验收：浏览器能流式完成一问一答；`/api/config` 锁定单模型并 `hideUserApiKey`。
- **不做（本 epic）**：多厂商 `/api/{provider}`、SafetyGate 挂在 `/api/*`、Tauri/桌面、关怀洪峰 UI、删除 NextChat 设置相关代码（仅隐藏）。

## 非目标

- 替换 `/api/bytedance` 为统一 `/v1/chat`（后续 epic）
- 与 ADR-0001 原「里程碑 B 全量 OpenAI + 豆包」对齐（已收窄为单模型）

## 参考

- 迁移清单：`frontend/MIGRATION.md`
- NextChat 代理对照：`NextChat/app/api/bytedance.ts`、`NextChat/app/api/config/route.ts`
