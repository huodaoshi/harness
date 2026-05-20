# backend/conf — 配置加载（Go 包，无 YAML）

本目录**仅**存放 `package conf` 源码：结构体、`Load()`、环境变量覆盖。

静态 YAML 在 [`../config/`](../config/README.md)（应用见 `config/app/`，Wellness 见 `config/*.yaml` 等）。

```go
cfg, err := conf.Load() // 读 config/app/config.yaml + config/app/{APP_ENV}.yaml
```
