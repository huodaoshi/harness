# family-wellness-ai backend

关系情绪 AI MVP 后端（Go + Hertz + Eino）。

## Spike S1（当前）

- `POST /v1/sessions/stream` — SSE：`token` + `done`
- Eino Graph 单节点假 ChatModel（无外部 LLM）

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
```
