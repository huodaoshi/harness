# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## Spike S1–S3（当前）

- `GET /v1/sessions/:id?user_id=` — 会话消息列表
- `POST /v1/sessions/end` — 结束会话并生成 `summary3`
- `GET /v1/profile?user_id=` — 关系档案（无档案时返回空对象）
- `PUT /v1/profile?user_id=` — 全量 upsert（`self`、`people[]`、`current_issue`）
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

环境变量：

| 变量 | 默认 |
| ---- | ---- |
| `MONGODB_URI` | `mongodb://localhost:27017` |
| `MONGODB_DB` | `family_wellness` |
| `USE_MEMORY_STORE` | 未设置 → 连 Mongo；`true` → 仅内存（测试用） |

## 开发

```text
cd backend
go test ./...
go run ./cmd/server
```

环境变量：

- `HTTP_ADDR` — 默认 `:8080`

## 验证

```text
curl -N -X POST http://localhost:8080/v1/sessions/stream ^
  -H "Content-Type: application/json" ^
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
