# NextChat → frontend 迁移清单

对照目录：`NextChat/`（只读）。术语与 Scaffold 范围见仓库根目录 `CONTEXT.md`（**Frontend Scaffold**）。

## Scaffold 验收（当前）

- [x] `GET /api/config`（Go + `hideUserApiKey` + 单模型 `customModels`）
- [x] `POST /api/bytedance/*` 透明代理（`ARK_*`，不经 SafetyGate）
- [x] 聊天主路径迁入，`yarn dev` 可打开 NextChat UI
- [x] 删除根目录 `web/` MVP
- [ ] Mask / locales 打磨（见 issue #04，不阻塞）

## 从 NextChat 复制

- [x] `app/components/`、`store/`、`client/`、`utils/`、`locales/`、`styles/`、`icons/`、`lib/`
- [x] `app/masks/` + `yarn mask`
- [x] `app/constant.ts`、`typing.ts`、`utils.ts`、`polyfill.ts`、`command.ts`
- [x] `public/`

## 不要复制

- `app/api/`（由 `backend/api/nextchat` 实现）
- `app/config/server.ts`
- `src-tauri/`、Vercel/Docker

## 改造要点

- `app/config/build.ts` / `client.ts`：无 Tauri 依赖
- `app/mcp/actions.ts`：Scaffold 期 MCP 占位（`isMcpEnabled() === false`）
- API 请求走同域 `/api`（`next.config.mjs` rewrites → backend :8080）
- 设置页多厂商入口由 `/api/config` 隐藏，**不删**相关源码

## 后续（非 Scaffold）

- 替换临时路径 `/api/bytedance`
- `/api/*` 挂 SafetyGate（若产品需要）
- 关怀模式 UI + `/v1/sessions/stream` 接入
- 收紧 `next.config.mjs` 中 `typescript.ignoreBuildErrors`（上游 TS 与 NextChat 组件命名）
