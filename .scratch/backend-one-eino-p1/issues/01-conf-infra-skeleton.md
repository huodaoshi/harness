# P1-01：配置与共享基础设施骨架

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

在 `harness/backend` 引入 one-eino 式 **`conf/`**（分层 YAML + 连接环境变量）与 **`infra/`**（MongoDB、Redis 客户端单例）。`cmd/server` 启动时加载配置、校验 Redis/Mongo 连通；HTTP 端口取自 `app.port`（仍兼容 `HTTP_ADDR` 覆盖）。

现有 Wellness API（`/v1/sessions/*`、`/v1/profile`）行为与测试**保持不变**；`internal/*` 尚未迁移。

## 验收标准

- [ ] `conf/config.yaml` + `conf/local.yaml` 存在；`APP_ENV=local` 可加载
- [ ] `MONGODB_URI` / `REDIS_ADDR` 等环境变量覆盖与 ADR 一致
- [ ] `go test ./...` 全绿
- [ ] `go run ./cmd/server` 启动日志含配置端口与 infra ping 结果
- [ ] 未设置 Redis 时启动失败有明确错误（或 local.yaml 默认本机 Redis）

## 阻塞于

无——可立即开始

## 覆盖的用户故事

基础设施；为 #12（鉴权）与知识库 issue 铺路
