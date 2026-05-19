# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## Spike S1–S3（当前）

- `GET /v1/profile?user_id=` — 关系档案（无档案时返回空对象）
- `PUT /v1/profile?user_id=` — 全量 upsert（`self`、`people[]`、`current_issue`）
- `POST /v1/sessions/stream` — SSE：`token` + `done`（普聊）；`crisis` / `medical` / `error`（门禁支路，零 LLM）
- Eino Graph：`safety_gate` → 分支 → `crisis_branch` | `medical_branch` | `block_branch` | `profile_inject` → `fake_chat`
- MongoDB：`relationship_profiles`、`session_summaries`（Level 1–2 记忆）
- 配置：`config/safety_rules_v1.yaml`、`config/crisis_templates/zh-CN.json`、`config/boundary_templates/zh-CN.json`
- 审计：`audit_event` 结构化日志（`gate_result`、`session_id`、`message_hash`，不落用户原文）

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
