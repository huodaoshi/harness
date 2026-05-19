# 内测鉴权 + Redis 限流

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

为 API 增加 **内测鉴权**（如 Bearer token 或魔法链接派生的 `user_id`），确保多用户隔离。Redis 实现按 IP 或 user 的 **速率限制**；超限返回 429 与 SSE `error` 可理解文案。

## 验收标准

- [ ] 无 token 请求 stream/profile 返回 401
- [ ] 两 user 无法互读 profile/session
- [ ] 超限触发 429，测试可稳定复现
- [ ] Web 在 401/429 时展示明确提示
- [ ] Redis 可本地 compose 启动

## 阻塞于

- [04-web-distress-chat-shell.md](./04-web-distress-chat-shell.md)

## 覆盖的用户故事

#20
