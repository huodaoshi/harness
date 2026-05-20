# Backend 测试目录

本模块测试与源代码**分目录**存放，结构镜像 `backend/` 包路径：

```text
backend/
  conf/config.go                              →  tests/conf/config_test.go
  modules/wellness/application/*.go           →  tests/modules/wellness/application/*_test.go
  modules/wellness/infra/store/*.go           →  tests/modules/wellness/infra/store/*_test.go
  api/*.go                                    →  tests/httpapi/*_test.go（包名 `apitest`）
```

- 包名一般为 `<pkg>_test`，通过 import 访问被测包。
- 在 `backend/` 下执行：`go test ./...`

**不要**在 `conf/`、`modules/`、`api/` 等源码目录旁新建 `*_test.go`。
