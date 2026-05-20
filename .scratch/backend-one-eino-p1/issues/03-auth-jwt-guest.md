# P1-03：Auth（SMS + JWT + 游客）

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

从 one-eino **复制改写** `modules/auth` 与 `api` 鉴权中间件。P1 启用：短信验证码登录、JWT refresh、**游客 `X-Anon-ID`** 访问 stream/profile。微信/扫码等配置项默认关闭。

废弃 query `?user_id=` / `X-User-Id` 作为租户键。游客 `user_id` 为 `anon:{uuid}`，**可读写关系档案**。

## 验收标准

- [x] 无 Authorization 且无合法 `X-Anon-ID` 时 stream/profile 返回 401
- [x] JWT 用户与游客 A 无法读取 B 的 session/profile
- [x] SMS 本地 mock（或 local 配置）可完成登录并拿 token
- [x] 现有集成测试改为 JWT/游客；`go test ./...` 绿
- [x] 对齐 [.scratch/family-wellness-ai/issues/12-auth-and-rate-limit.md](../../family-wellness-ai/issues/12-auth-and-rate-limit.md) 的隔离要求（限流见 04）

## 阻塞于

- [02-wellness-modules-migration.md](./02-wellness-modules-migration.md)

## 覆盖的用户故事

#5、#12、#14、#20
