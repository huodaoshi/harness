# backend/config — 静态配置文件（无 Go 源码）

本目录**仅**存放 YAML/JSON 等配置文件，由 `conf` 包或 `modules/wellness/infra/configpaths` 读取。

## 布局

| 路径 | 用途 | 加载方 |
|------|------|--------|
| `app/config.yaml` | 应用默认配置（端口、Mongo、Redis、JWT、MQ 等） | `conf.Load()` |
| `app/{APP_ENV}.yaml` | 环境叠加（如 `local.yaml`） | `conf.Load()` |
| `safety_rules_v1.yaml` | SafetyGate 规则 | wellness `configpaths` |
| `crisis_templates/` | 危机文案模板 | wellness `configpaths` |
| `boundary_templates/` | 边界/拦截文案 | wellness `configpaths` |

## 约定

- **勿**在本目录添加 `.go` 文件；配置加载逻辑在 `backend/conf/`。
- 在 `backend/` 目录启动服务，路径相对模块根解析。
- 敏感项用 `${ENV_VAR}`，由 `conf/connection_env.go` 与环境变量覆盖。
