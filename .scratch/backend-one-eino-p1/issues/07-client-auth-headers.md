# P1-07：客户端鉴权头迁移

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

更新仓库内调用 Wellness API 的客户端（`web/`、`frontend/` 关怀路径等）：发送 `Authorization: Bearer` 或 `X-Anon-ID`；移除 query `user_id`。401/429 展示明确提示。

## 验收标准

- [ ] 本地联调：游客可洪峰聊天；登录用户 profile 隔离
- [ ] 文档/README 更新启动步骤（auth + redis）

## 阻塞于

- [03-auth-jwt-guest.md](./03-auth-jwt-guest.md)

## 覆盖的用户故事

#4、#5、#20、#24
