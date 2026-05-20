# frontend

Harness 聊天 UI（从 [NextChat](../NextChat/) 迁移）。**仅 Web**；API 由 [backend](../backend/) 提供，本目录不包含 `app/api/`。

## 本地开发（Scaffold E2E）

1. 配置 backend：复制 `config/app/local.secrets.yaml.example` → `local.secrets.yaml` 并填写 `llm.api_key`、`llm.model`（见 `backend/config/app/README.md`）
2. 启动 backend：

```powershell
cd d:\harness\backend
go run ./cmd/server
```

3. 启动 frontend：

```powershell
cd d:\harness\frontend
yarn dev
```

打开 http://localhost:3000 。`/api/*` 经 Next **rewrites** 转到 backend `:8080`。

**关怀模式**（`/v1/*`）：见 `lib/harness/auth.ts`、`headers.ts`。

## 构建

```bash
yarn build
```

Scaffold 期为兼容上游 NextChat 代码，构建时暂时忽略部分 ESLint/TS 检查（见 `next.config.mjs`）。

## 文档

- 术语：`CONTEXT.md`（Frontend Scaffold）
- 迁移：`MIGRATION.md`
- ADR：`docs/adr/0001-frontend-backend-split.md`（部分范围已收窄）
