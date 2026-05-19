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

- [ ] 10/10 危机剧本触发 `crisis` 事件，无 `token` 流
- [ ] 测试中断言 LLM/ChatModel mock 调用次数为 0
- [ ] Graph 危机边在代码审查上可见「未挂载 ChatModel 节点」
- [ ] 危机路径 operational 日志不落用户原文（仅 gate 结果或 hash）
- [ ] 与 #01 共用同一 SSE 端点，不新建平行 API

## 阻塞于

- [01-spike-s1-stream-sse.md](./01-spike-s1-stream-sse.md)

## 覆盖的用户故事

#9、#10、#13、#17、#18、#25
