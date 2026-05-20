# Backend 对齐 one-eino（P1）

**上游：** [ADR-0002](../../docs/adr/0002-backend-one-eino-alignment.md)、[CONTEXT.md](../../CONTEXT.md)

P1 目标：平台骨架 + auth + 知识库；**不改** `/v1/sessions/stream` 编排契约。P2（ADK 合流）另开 epic。

## Issue 顺序

| # | 文件 | 说明 |
|---|------|------|
| 01 | [01-conf-infra-skeleton.md](./issues/01-conf-infra-skeleton.md) | conf + infra + 启动接线 |
| 02 | [02-wellness-modules-migration.md](./issues/02-wellness-modules-migration.md) | wellness 迁入 modules/ + api/ |
| 03 | [03-auth-jwt-guest.md](./issues/03-auth-jwt-guest.md) | SMS/JWT/游客；替代 query user_id |
| 04 | [04-stream-rate-limit.md](./issues/04-stream-rate-limit.md) | Redis 限流 |
| 05 | [05-knowledge-ingest-e2e.md](./issues/05-knowledge-ingest-e2e.md) | 默认 Space + 精简入库 E2E |
| 06 | [06-knowledgeindexing-worker.md](./issues/06-knowledgeindexing-worker.md) | 独立 Worker |
| 07 | [07-client-auth-headers.md](./issues/07-client-auth-headers.md) | 前端/客户端鉴权头 |
