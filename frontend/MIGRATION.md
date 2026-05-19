# NextChat → frontend 迁移清单

对照目录：`NextChat/`（只读）。架构决策见 `docs/adr/0001-frontend-backend-split.md`。

## 第一期（里程碑 B）

### 从 NextChat 复制（保持路径在 `app/` 下）

- [ ] `app/components/`
- [ ] `app/store/`
- [ ] `app/client/`（将 `path()` 中的 baseUrl 改为 `getApiBaseUrl()` + `/api/...`）
- [ ] `app/masks/` + `scripts`：`mask` / `mask:watch`（见 `package.json`）
- [ ] `app/locales/`
- [ ] `app/utils/`
- [ ] `app/constant.ts`、`app/typing.ts`、`app/utils.ts`、`app/polyfill.ts`
- [ ] `app/styles/`
- [ ] `public/`

### 不要复制

- [ ] `app/api/`（已由 `backend/` 实现）
- [ ] `app/config/server.ts`（密钥改读 backend 环境变量）
- [ ] `src-tauri/`、`vercel.json`、根目录 Docker（除非单独决策）

### 依赖与配置

- [ ] 合并 `NextChat/package.json` 中运行时依赖
- [ ] `next.config.mjs`：SVG（`@svgr/webpack`）、Sass、必要时 `rewrites` → backend
- [ ] `app/config/client.ts`：仅客户端构建配置；运行时拉取 `GET {API_BASE}/api/config`

### 改造点

- [ ] 所有 `ApiPath.*` 请求前缀使用 `lib/harness/api-base.ts`
- [ ] `lib/harness/headers.ts`：预留 `X-Harness-Mode`（第一期无 UI）
- [ ] 删除对 `getServerSideConfig()` 的前端密钥依赖

### 验收

- [ ] `backend` + `frontend` 同时运行，能打开聊天页
- [ ] 豆包流式回复正常
- [ ] 切换 Mask 正常
- [ ] 设置页能读取 `/api/config`

## 第二期

- [ ] 其余 `/api/{provider}`
- [ ] `app/mcp/`（Server Actions）
- [ ] 插件、`/api/proxy`、WebDAV、Upstash
- [ ] 关怀模式设置项 + `/v1/profile` + `/v1/sessions/stream`
- [ ] 删除 `web/`
