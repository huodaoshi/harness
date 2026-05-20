# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## 配置

**不使用 `.env` 文件。** 见 [`config/app/README.md`](config/app/README.md)。

```powershell
copy config\app\local.secrets.yaml.example config\app\local.secrets.yaml
# 编辑 llm.api_key、llm.model
```

## 目录（ADR-0002）

```text
backend/
  api/                          # HTTP 路由与 handler
  modules/wellness/
    domain/                     # Store 接口与模型
    application/                # 会话图、Executor
    infra/{store,safety,chatmodel,configpaths}/
  tests/                        # 唯一允许的 *_test.go 位置（镜像源码路径）
  scripts/check_test_placement/ # 校验：tests/ 外不得有 *_test.go
  conf/          # Go：conf.Load()
  config/        # 静态 YAML（app/ + wellness）
  cmd/server/
```

## Spike S1–S3（当前）

- `GET /v1/sessions/:id` — 会话消息列表（需 JWT 或 `X-Anon-ID`）
- `POST /v1/sessions/end` — 结束会话并生成 `summary3`
- `GET /v1/profile` — 关系档案（无档案时返回空对象）
- `PUT /v1/profile` — 全量 upsert（`self`、`people[]`、`current_issue`）
- `POST /v1/sessions/stream` — SSE：`token` + `done`（普聊）；`crisis` / `medical` / `error`（门禁支路，零 LLM）
- NextChat 兼容：`GET /api/config`、`/api/bytedance/*`（配置见 `config/app` + `api/nextchat`）

### LLM（YAML `llm:` + 可选环境变量）

| 字段 / 变量 | 说明 |
| ----------- | ---- |
| `llm.provider` | `fake` 或 `ark` |
| `llm.api_key` / `ARK_API_KEY` | 火山方舟 API Key（勿入库） |
| `llm.model` / `ARK_MODEL_ID` | 接入点模型 ID |
| `llm.request_timeout` | 默认 `90s` |
| `llm.first_token_target_ms` | 首 token 目标（日志），默认 `3000` |

### 本地 Mongo

```text
cd backend
docker compose up -d
```

启动 **Mongo + Redis**（鉴权、限流依赖 Redis）。

环境变量与 YAML 字段详见 `config/app/README.md`。

```text
cd backend
go test ./...
go run ./cmd/server
```

配置（ADR-0002 P1-01）：

- 分层 YAML：`config/app/config.yaml` + `config/app/{APP_ENV}.yaml` + 可选 `{APP_ENV}.secrets.yaml`
- 须在 **`backend/` 目录**下启动；`conf/` 为 Go 包，`config/` 仅静态文件（见 `config/app/README.md`）
- `HTTP_ADDR` 环境变量可覆盖 `app.port`；Mongo/Redis 等见 `conf/connection_env.go`

## 知识库（P1-05 / P1-06）

- 默认 Space `space_id=1`（doc_type 1/2/3），启动时 bootstrap
- Admin 入库：`POST /v1/admin/spaces/:space_id/ingest`（inline markdown，`source_type=1`）
- 检索验收：`GET /v1/admin/spaces/:space_id/search?q=...`（需 Redis Stack / RediSearch）
- 流水线：`resolve → parse → chunk → assignID → index`（**无** LLM score/tag）
- **MQ**：
  - `mq.provider=local`：`cmd/server` **内嵌**消费 `knowledge-ingest`（同进程）
  - `mq.provider=rocketmq`：独立运行 `go run ./cmd/knowledgeindexing`
- 本地 embedding 默认 `fake`（`config/app/local.yaml`）；生产可设 `embedding.provider=ark`

## 鉴权（P1-03）

- Wellness：`Authorization: Bearer <token>` **或** `X-Anon-ID: <uuid>`（游客 `user_id`=`anon:{uuid}`）
- 已废弃：query `?user_id=` / `X-User-Id`
- 短信：`POST /v1/auth/sms/send`、`POST /v1/auth/sms/verify`（`sms.provider=local` 时验证码打日志）
- 本地需 **Redis** 与 **jwt.secret**（`local.yaml` 或 `JWT_SECRET` 环境变量）

## 验证

```text
curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
  -H "X-Anon-ID: 11111111-1111-4111-8111-111111111111" ^
  -d "{\"message\":\"hello\",\"mode\":\"distress\"}"
```
