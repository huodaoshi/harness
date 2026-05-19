# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## Spike S1 + S2（当前）

- `POST /v1/sessions/stream` — SSE：`token` + `done`（普聊）或 `crisis`（自伤/家暴 L1 规则）
- Eino Graph：`safety_gate` → 分支 → `crisis_branch` | `fake_chat`（危机不调 fake Chat）
- 配置：`config/safety_rules_v1.yaml`、`config/crisis_templates/zh-CN.json`

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
