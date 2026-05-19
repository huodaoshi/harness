# Spike S2：SafetyGate 危机支路 + 零 LLM 调用

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — W1 门禁 S2

## 要构建什么

在 RelationshipSessionGraph 中接入 **SafetyGate（L1 规则）**：对自伤、家暴类输入走 **CrisisBranch**，SSE 发出 `event: crisis`（含 `template_id`、`body`），**编译期不连接 ChatModel**。

提供外置 `safety_rules_v1.yaml` 与 `crisis_templates/zh-CN.json`（占位热线可标「待核实」）。自动化脚本跑 **10 条危机剧本**，断言：SSE 为 `crisis`；假 ChatModel **调用次数 = 0**。

## 验收标准

- [x] 10/10 危机剧本触发 `crisis` 事件，无 `token` 流
- [x] 测试中断言 LLM/ChatModel mock 调用次数为 0（`FakeChatCallCounter`）
- [x] Graph：`safety_gate` → branch → `crisis_branch` | `fake_chat`（无 crisis→chat 边）
- [x] 危机路径 `safety.Audit` 不落用户原文
- [x] 与 #01 共用 `POST /v1/sessions/stream`

## 评论

### 2026-05-19 · 实现完成

- `backend/internal/safety`、`session.Executor`、危机 SSE
- 测试：`crisis_test.go` + `httpserver` 10 条剧本 E2E

## 阻塞于

- [01-spike-s1-stream-sse.md](./01-spike-s1-stream-sse.md)

## 覆盖的用户故事

#9、#10、#13、#17、#18、#25
