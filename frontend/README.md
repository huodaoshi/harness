# frontend

Harness 聊天 UI（从 [NextChat](../NextChat/) 迁移）。API 由 [backend](../backend/) 提供，本目录**不包含** `app/api/`。

## 本地开发

1. 配置 `backend/.env`（见 `backend/.env.example`）
2. 启动 backend：`cd backend && go run ./cmd/server`（默认 `:8080`）
3. 本目录：

```bash
cp .env.example .env.local
yarn dev
```

打开 http://localhost:3000 。默认通过 Next **rewrites** 将 `/api`、`/v1` 转发到 backend，无需 CORS。

**关怀模式鉴权**（`lib/harness/auth.ts` + `headers.ts`）：调用 `/v1/*` 时 `withHarnessHeaders()` 在 `wellness` 模式下自动附带 `Authorization: Bearer` 或 `X-Anon-ID`。

> 若系统环境变量 `NODE_ENV=windows`，请在终端先执行 `$env:NODE_ENV="development"`，或使用 VS Code 调试配置「Frontend: 开发服务器」。

## 文档

- 架构：[`docs/adr/0001-frontend-backend-split.md`](../docs/adr/0001-frontend-backend-split.md)
- 迁移清单：[`MIGRATION.md`](./MIGRATION.md)
