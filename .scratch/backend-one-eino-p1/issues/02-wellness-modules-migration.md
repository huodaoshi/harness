# P1-02：Wellness 迁入 modules 与 api 层

**Status:** ready-for-human  
**类型：** enhancement  
**切片：** AFK  

## 父级

[ADR-0002](../../../docs/adr/0002-backend-one-eino-alignment.md)

## 要构建什么

将现有 `internal/session`、`internal/safety`、`internal/chatmodel`、`internal/store` 与 HTTP 处理逻辑迁入 **`modules/wellness/`**（domain / application / infra）及 **`api/`** 路由注册。`cmd/server` 仅负责组装依赖与挂路由。

对外路径与 SSE 事件形状**不变**；危机/医疗零 LLM 回归测试仍通过。

## 验收标准

- [x] 目录符合 `modules/<name>/{domain,application,infra}` + `api/*.go`
- [x] `go test ./...` 全绿（含 crisis、medical、profile inject）
- [x] 迁出的 `*_test.go` 置于 `tests/` 镜像目录（见 `.cursor/rules/_lang/go.md`），不在 `internal/` 旁保留
- [x] 手工 `curl` stream 洪峰/危机用例与迁移前一致
- [x] 根目录无新增业务逻辑于 `internal/`（仅保留测试辅助或删除已迁代码）

## 阻塞于

- [01-conf-infra-skeleton.md](./01-conf-infra-skeleton.md)

## 覆盖的用户故事

#1–#13、#17–#18、#21–#25（Wellness 路径）
