# P1-04：Stream Redis 限流

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

接入 one-eino 式 Redis 限流器，对 `POST /v1/sessions/stream` 按 user（含 `anon:*`）或 IP 限速。超限返回 429；SSE 错误帧文案可理解。

## 验收标准

- [ ] 超限稳定返回 429，测试可复现
- [ ] 未超限流不影响危机分支零 LLM 行为
- [ ] `docker compose` 或文档说明本地 Redis

## 阻塞于

- [03-auth-jwt-guest.md](./03-auth-jwt-guest.md)

## 覆盖的用户故事

#20
