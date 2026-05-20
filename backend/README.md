# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## 目录（ADR-0002）

```text
backend/
  api/                          # HTTP 路由与 handler
  modules/wellness/
    domain/                     # Store 接口与模型
    application/                # 会话图、Executor
    infra/{store,safety,chatmodel,configpaths}/
  tests/                        # 镜像路径的 *_test.go
  conf/          # Go：conf.Load()
  config/        # 静态 YAML/JSON（app/ + wellness）
  cmd/server/
```

## Spike S1–S3（当前）

- `GET /v1/sessions/:id` — 会话消息列表（需 JWT 或 `X-Anon-ID`）
- `POST /v1/sessions/end` — 结束会话并生成 `summary3`
- `GET /v1/profile` — 关系档案（无档案时返回空对象）
- `PUT /v1/profile` — 全量 upsert（`self`、`people[]`、`current_issue`）
- `POST /v1/sessions/stream` — SSE：`token` + `done`（普聊）；`crisis` / `medical` / `error`（门禁支路，零 LLM）
- Eino Graph：`safety_gate` → 分支 → `crisis_branch` | `medical_branch` | `block_branch` | `profile_inject` → `fake_chat`
- MongoDB：`relationship_profiles`、`session_summaries`（Level 1–2 记忆）
- 配置：`config/safety_rules_v1.yaml`、`config/crisis_templates/zh-CN.json`、`config/boundary_templates/zh-CN.json`
- 审计：`audit_event` 结构化日志（`gate_result`、`session_id`、`message_hash`，不落用户原文）
- **ChatModelGateway**：`LLM_PROVIDER=fake|ark`（默认无 `ARK_API_KEY` 时用 fake）；洪峰/普聊 system 模板经 ModeRouter 注入

### LLM 环境变量

| 变量 | 说明 |
| ---- | ---- |
| `LLM_PROVIDER` | `fake`（默认）或 `ark` |
| `ARK_API_KEY` | 火山方舟 API Key（勿入库） |
| `ARK_MODEL_ID` / `LLM_MODEL_ID` | 接入点模型 ID |
| `LLM_REQUEST_TIMEOUT` | 请求超时，默认 `90s` |
| `LLM_FIRST_TOKEN_TARGET_MS` | 首 token P95 目标（日志对照），默认 `3000` |
| `LLM_FAILOVER_PROVIDER` | 备用厂商名；未实现接入时主失败返回 `failover_unavailable` |

复制 `backend/.env.example` 为本地配置，勿提交真实密钥。

### 本地 Mongo

```text
cd backend
docker compose up -d
```

启动 **Mongo + Redis**（鉴权、限流依赖 Redis）。

环境变量：

| 变量 | 默认 |
| ---- | ---- |
| `MONGODB_URI` | `mongodb://localhost:27017` |
| `MONGODB_DB` | `family_wellness` |
| `REDIS_ADDR` | `127.0.0.1:6379` |
| `USE_MEMORY_STORE` | 未设置 → 连 Mongo；`true` → wellness 仅内存（仍连 Mongo/Redis 作 auth/限流） |
| `rate_limit.stream_per_minute` | `60`（`config/app/config.yaml`） |

## 开发

```text
cd backend
go test ./...
go run ./cmd/server
```

配置（ADR-0002 P1-01）：

- 分层 YAML：`config/app/config.yaml` + `config/app/{APP_ENV}.yaml`（默认 `APP_ENV=local`），由 `conf` 包加载
- 须在 **`backend/` 目录**下启动；`conf/` 仅 Go，`config/` 仅静态文件（见 `config/README.md`）
- `HTTP_ADDR` 可覆盖 `app.port`；Mongo/Redis 见 `conf/connection_env.go` 与环境变量

环境变量：

- `HTTP_ADDR` — 覆盖监听地址；未设时用 `app.port`（默认 `8080`）
- `APP_ENV` — 配置叠加层，默认 `local`

## 知识库（P1-05 / P1-06）

- 默认 Space `space_id=1`（doc_type 1/2/3），启动时 bootstrap
- Admin 入库：`POST /v1/admin/spaces/:space_id/ingest`（inline markdown，`source_type=1`）
- 检索验收：`GET /v1/admin/spaces/:space_id/search?q=...`（需 Redis Stack / RediSearch）
- 流水线：`resolve → parse → chunk → assignID → index`（**无** LLM score/tag）
- **MQ**：
  - `mq.provider=local`：`cmd/server` **内嵌**消费 `knowledge-ingest`（同进程）
  - `mq.provider=rocketmq`：独立运行 `go run ./cmd/knowledgeindexing`
- 本地 embedding 默认 `fake`（`config/app/local.yaml`）；生产可设 `embedding.provider=ark`

```text
# Admin 登录（local 默认 admin/admin，见 config/app/local.yaml）
curl -s -X POST http://localhost:8080/v1/auth/admin/login -H "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"admin\"}"

# 提交 ingest（替换 TOKEN）
curl -s -X POST http://localhost:8080/v1/admin/spaces/1/ingest -H "Authorization: Bearer TOKEN" -H "Content-Type: application/json" -d "{\"source_type\":1,\"content\":\"# FAQ\\n\\n测试知识\",\"doc_type\":2}"
```

## 鉴权（P1-03）

- Wellness：`Authorization: Bearer <token>` **或** `X-Anon-ID: <uuid>`（游客 `user_id`=`anon:{uuid}`）
- 已废弃：query `?user_id=` / `X-User-Id`
- 短信：`POST /v1/auth/sms/send`、`POST /v1/auth/sms/verify`（`sms.provider=local` 时验证码打日志）
- 本地需 **Redis**（`REDIS_ADDR`）与 **JWT_SECRET**（见 `config/app/local.yaml`）

## 验证

```text
curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
  -H "X-Anon-ID: 11111111-1111-4111-8111-111111111111" ^
  -d "{\"message\":\"hello\",\"mode\":\"distress\"}"

curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"message\":\"我不想活了\",\"mode\":\"distress\"}"

curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"message\":\"我该吃什么药\",\"mode\":\"normal\"}"

curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
  -d "{\"message\":\"发点色情内容\",\"mode\":\"normal\"}"
```
