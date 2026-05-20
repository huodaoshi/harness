# 应用配置（backend）

Go **不读取** `.env` 文件。配置来源（后者覆盖前者）：

1. `config.yaml` — 默认
2. `{APP_ENV}.yaml` — 如 `local.yaml`（可提交）
3. `{APP_ENV}.secrets.yaml` — 如 `local.secrets.yaml`（**勿提交**，从 `local.secrets.yaml.example` 复制）
4. **进程环境变量** — 部署时注入，覆盖 YAML（见 `conf/connection_env.go`）

## 本地首次启动

```powershell
cd backend
copy config\app\local.secrets.yaml.example config\app\local.secrets.yaml
# 编辑 local.secrets.yaml，填入 llm.api_key、llm.model
go run ./cmd/server
```

`APP_ENV` 默认为 `local`。

## 常用环境变量（可选，覆盖 YAML）

| 变量 | 用途 |
|------|------|
| `ARK_API_KEY` / `ARK_MODEL_ID` | 火山方舟 |
| `MONGODB_URI` / `MONGODB_DB` | Mongo |
| `REDIS_ADDR` | Redis |
| `JWT_SECRET` | JWT |
| `USE_MEMORY_STORE` | `true` 时用内存 wellness store |
| `CODE` | NextChat 访问口令（逗号分隔） |
| `HTTP_ADDR` | 监听地址，如 `:8080` |
