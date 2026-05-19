# L2 云内容安全（fail-closed）

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md) — Spike S6

## 要构建什么

实现 `ModerationClient` 抽象与一家国内云审供应商适配。SafetyGate 在 L1 `pass` 后调用 L2；**staging/prod 配置 fail-closed**（服务不可用则拒绝洪峰会话或全拒绝，行为写进配置文档）。

dev 可用 fake moderation 默认 pass。

## 验收标准

- [ ] pass 消息在云审通过后进入 ChatModel 路径
- [ ] 云审返回 block 时 SSE `error` 或 block，不调 LLM
- [ ] fail-closed：mock 超时/500 时不静默放行（staging 实测）
- [ ] 集成测试使用 fake client，不依赖外网
- [ ] 不回归 #02 危机零 LLM

## 阻塞于

- [05-safety-medical-block-audit.md](./05-safety-medical-block-audit.md)

## 覆盖的用户故事

#12
