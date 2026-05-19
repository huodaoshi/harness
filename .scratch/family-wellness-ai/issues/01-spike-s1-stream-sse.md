# Spike S1：Hertz → Eino → 假 ChatModel → SSE 流式

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — W1 门禁 S1

## 要构建什么

搭建 `backend/` 最小骨架：Hertz 暴露 `POST /v1/sessions/stream`，请求体含 `message`、`mode`；内部经 Eino `compose.Graph` 调用**假 ChatModel**（固定或逐字吐出 token），响应为 `text/event-stream`，至少发出 `token` 与 `done` 事件。

本切片不要求真实 LLM、Mongo、SafetyGate 分支；目标是证明**编排 + SSE 管道可跑通**（可用 curl 或最小 HTML 验证）。

## 验收标准

- [x] `POST /v1/sessions/stream` 返回 `Content-Type: text/event-stream`
- [x] 连续收到 `event: token` 与最终 `event: done`（含 `session_id`）
- [x] 使用假 ChatModel，无外部 API 依赖即可在 dev 启动
- [x] 集成测试或脚本可自动化断言 SSE 序列（`backend/internal/httpserver` E2E）
- [x] Eino v0.8.13 / Hertz v0.10.4 在 `go.mod` 中 pin；符合 `.cursor/rules/eino.md`

## 评论

### 2026-05-19 · 实现完成

- 代码：`backend/`（`go test ./...` 通过）
- 运行：`go run ./cmd/server` → `curl` 见 `backend/README.md`

## 阻塞于

无——可立即开始

## 覆盖的用户故事

#2、#17
