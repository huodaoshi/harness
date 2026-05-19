# SafetyGate 补全：医疗模板、block、审计事件

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

在 L1 规则上补全：**medical_boundary**（固定模板、不调 LLM）、**block**（拒绝、不调 LLM）。可选写入 `audit_events`：`gate_result`、`session_id`、`timestamp`，**不存用户原文**。

SSE：`medical`/`block` 走 `error` 或专用事件（与 PRD 一致并在实现中写清契约）。扩展规则表与模板配置，供合规更新热线文案。

## 验收标准

- [x] 用药/诊断类输入返回固定边界话术，ChatModel 调用次数 0
- [x] 违禁类输入被 block，用户可见明确拒绝
- [x] 审计记录无完整用户消息正文
- [x] 表驱动测试覆盖 medical + block 各 ≥3 条
- [x] 与 #02 危机用例共存，无回归

## SSE 契约（#05 实现）

| event | payload | 说明 |
| ----- | ------- | ---- |
| `medical` | `{ template_id, body }` | 医疗边界固定模板 |
| `error` | `{ code: "content_blocked", message }` | 违禁 block（PRD 阻断形态） |
| `crisis` | 不变 | #02 |

## 阻塞于

- [02-spike-s2-crisis-zero-llm.md](./02-spike-s2-crisis-zero-llm.md)

## 覆盖的用户故事

#11、#12、#13、#16、#25
