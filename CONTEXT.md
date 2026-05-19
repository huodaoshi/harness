# Harness 领域上下文

## 术语表

| 术语 | 含义 |
|------|------|
| **backend** | Go + Hertz 服务，托管 `/v1/*`（wellness）与 `/api/*`（NextChat 兼容代理） |
| **frontend** | 从 NextChat 迁移的 Next.js 聊天 UI，不含 `app/api` |
| **NextChat** | 仓库内只读参考副本，不直接修改 |
| **通用模式** | 默认聊天：走 `/api/*`，历史存浏览器本地（IndexedDB） |
| **关怀模式** | wellness：请求头 `X-Harness-Mode: wellness`，走 `/v1/sessions/*` 与 SafetyGate |
| **里程碑 B** | 第一期：聊天主路径 + Mask + Go 代理（config / bytedance / openai） |
| **里程碑 C** | 第二期：插件、MCP、同步、全厂商代理、关怀模式 UI |

## 架构决策

见 `docs/adr/0001-frontend-backend-split.md`。

## 本地开发

```text
backend/   → :8080  （Go API）
frontend/  → :3000  （Next.js；/api 可经 rewrite 转发到 backend）
```
