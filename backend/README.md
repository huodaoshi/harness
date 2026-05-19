# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## Spike S1–S3（当前）

- `POST /v1/sessions/stream` — SSE：`token` + `done`（普聊）或 `crisis`（自伤/家暴 L1 规则）
- Eino Graph：`safety_gate` → 分支 → `crisis_branch` | `profile_inject` → `fake_chat`
- MongoDB：`relationship_profiles`、`session_summaries`（Level 1–2 记忆）
- 配置：`config/safety_rules_v1.yaml`、`config/crisis_templates/zh-CN.json`

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
```
