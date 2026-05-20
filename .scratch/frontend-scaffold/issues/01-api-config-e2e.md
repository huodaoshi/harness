# Scaffold-01：`GET /api/config` 竖切

**Status:** ready-for-agent  
**类型：** enhancement  
**切片：** AFK  

## 父级

[PRD.md](../PRD.md)

## 要构建什么

在 **Go `backend`** 实现 NextChat 兼容的 **`GET /api/config`**（及 `POST` 若 NextChat 客户端会发），响应 JSON 字段与 NextChat `DANGER_CONFIG` 对齐，至少包含：

- `needCode`：未配置 `CODE` 时为 `false`
- `hideUserApiKey`：托管单模型场景为 `true`
- `customModels`：单条模型（如 `-all,+显示名@bytedance=<ARK_MODEL_ID>` 格式，与 NextChat 约定一致）
- `defaultModel`：与上条一致

配置从环境变量读取：**`ARK_*`**（及既有 `LLM_*` 覆盖）；**不**引入 `BYTEDANCE_*` 为正式变量名。

**`frontend/`** 在现有探针页或最小客户端逻辑中，经同域 rewrite 成功拉取并展示关键字段（或 NextChat `getClientConfig` 路径能消费）。

注册路由时注意与现有 `/v1/*`、静态 `web/` 不冲突；必要时 CORS 仅面向本地 dev。

## 验收标准

- [ ] `backend` 启动后 `curl http://localhost:8080/api/config` 返回 200 与上述字段
- [ ] 未设置 `CODE` 时 `needCode=false`
- [ ] 设置 `CODE` 后 `needCode=true`（口令校验可在 #02 与聊天代理一并接线，或本 issue 仅返回标志位）
- [ ] `frontend` `yarn dev`（:3000）下访问页面可见 config 探针成功（非 `fetch failed`）
- [ ] `go test` 新增/现有测试不回归

## 阻塞于

无——可立即开始

## 覆盖的用户故事

基础设施；为 #02、#03 铺路
