# ADR-0001：NextChat 能力迁移与前后端分离

- **状态**：已接受
- **日期**：2026-05-19
- **决策者**：产品与工程（grill-me 共识）

## 背景

仓库内存在：

- `backend/`：Go + Hertz + Eino，提供 wellness 会话（`/v1/sessions/*`）、火山方舟豆包与安全门禁；
- `web/`：临时 MVP 静态页，计划删除；
- `NextChat/`：完整聊天 UI + Next.js 内置 `/api/*` 多厂商代理，**保持只读参考，不在此目录修改**。

目标：将 NextChat 能力迁移到与 `backend/` 同级的新前端 `frontend/`，API 由 Go 承接，实现前后端分离，并保留后续接入关怀模式（SafetyGate）的路径。

## 决策

### 1. 后端：扩展现有 `backend/`

不新建 Node BFF。在 Go 中实现与 NextChat 兼容的 `/api/*` 代理契约（路径与请求体不变）。

### 2. API 契约

前端（迁移后）继续调用 `/api/config`、`/api/openai`、`/api/bytedance` 等；密钥仅存在于 `backend/.env`，通过 `/api/config` 下发运行时开关（如 `hideUserApiKey`），不下发明文 Key。

### 3. 安全双轨（可开关）

| 模式 | 触发 | 聊天路径 | 存储 |
|------|------|----------|------|
| 通用（默认） | 无 header / 未开全局开关 | `/api/*` 透明代理 | 前端 IndexedDB（NextChat 同构） |
| 关怀（wellness） | `X-Harness-Mode: wellness` 或 `SAFETY_GATE_API=on` | `/v1/sessions/stream` + SafetyGate | Mongo（现有 backend） |

第一期（里程碑 B）关怀模式**仅埋点**（模式枚举 + header），不做设置 UI；第二期交付切换与 Profile。

### 4. 前端：新建 `frontend/`

- 使用 `create-next-app`（Next.js 14 App Router + TypeScript）空壳；
- 从 `NextChat/` **按目录迁移**业务代码，不复制 `app/api/`、`src-tauri/`、Vercel/Docker 等部署杂物；
- `NextChat/` 保留为 diff 对照。

### 5. 分期交付

**第一期（B）**

- 前端：聊天主路径 + Mask；
- Go：`/api/config`、鉴权、`/api/bytedance`、`/api/openai`、CORS；其余厂商 501。

**第二期**

- 其余 `/api/*` 厂商、插件、`/api/proxy`、MCP（Server Actions）、WebDAV/Upstash、关怀模式 UI、`/v1/profile`。

### 6. 仓库布局（目标态）

```text
harness/
  backend/     # Go API
  frontend/    # Next.js UI（迁移目标）
  NextChat/    # 只读参考
  web/         # 待删 MVP
```

## 后果

### 正面

- 统一技术栈（Go 后端 + 独立 Next 前端）；
- wellness 与通用聊天可并存，互不阻塞迁移；
- `NextChat/` 可继续跟踪上游而不污染产品目录。

### 负面 / 风险

- Go 需逐厂商 port 代理逻辑，第一期仅覆盖豆包 + OpenAI；
- MCP 依赖 Next Server Actions，第二期需保留 Node 端能力或另建服务；
- 开发期需同时起 `backend` 与 `frontend`（或配置反向代理）。

## 参考

- 迁移清单：`frontend/MIGRATION.md`
- 术语：`CONTEXT.md`
- NextChat 代理实现对照：`NextChat/app/api/`
