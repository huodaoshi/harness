---
description: Go 编码通用规范——分层、错误处理、命名、并发、依赖管理
paths:
  - "**/*.go"
---

# Go 编码规范

## 分层架构

服务/应用代码应遵循单向依赖的分层结构（具体层名以项目实际命名为准，常见组合：Controller / Handler / API → Service / UseCase → Repo / DAO → Model）。

- 顶层：参数绑定、认证 / 上下文提取、调用业务层、返回响应
- 业务层：业务逻辑编排、事务管理、事件通知
- 数据层：数据库 CRUD、缓存管理、查询构建
- 模型层：结构体定义、表映射、自定义类型

> **禁止**：顶层直接操作 `*sql.DB` / `*gorm.DB`；业务层直接写 HTTP 响应；数据层包含业务判断逻辑。

## 接口优先

业务层与数据层应定义公开接口 + 私有实现：

```go
type Example interface {
    DoSomething(ctx context.Context, in *DoSomethingInput) error
}

type example struct { /* deps */ }

func New(deps ...) Example { return &example{...} }
```

## 错误处理

- 不忽略 error 返回值（`errcheck` lint）
- 错误传递时用 `fmt.Errorf("...: %w", err)` 包装上下文
- 数据层错误应统一包装，向上层屏蔽底层细节
- 用户面错误与内部错误应有不同的返回方式（具体方案以项目 skills 中描述为准）

## 命名约定

- 包名：小写、单词、避免下划线
- 接口：动词或 `er` 后缀（`Reader`、`Writer`）；非 `er` 命名时（如领域接口）首字母大写
- 导出符号：首字母大写
- 缩写：保持全大写（`HTTP`、`URL`、`API`、`ID`、`UID`）——`sessionID` 不写成 `sessionId` 或 `session_id`
- 文件名：小写下划线分隔
- JSON 序列化字段用 `snake_case` tag（`json:"session_id"`），Go 字段名保持 `SessionID`；两者不相同是正确的，见 `rules/naming.md`

## 上下文与并发

- 函数第一个参数总是 `ctx context.Context`
- 启动 goroutine 必须考虑取消与回收
- 使用 `errgroup`、`sync.WaitGroup` 管理并发
- 共享状态用 `sync.Mutex` 或 channel 保护

## 依赖管理

- 使用 go modules
- 不在 `init()` 中执行业务初始化（除注册 driver / encoding 之类）
- 单例使用 `sync.Once` 模式

## 格式化

- 使用 `goimports`（统一 import 顺序：标准库 / 第三方 / 内部）
- 缩进：**Tab 字符**（`gofmt` 强制，不使用空格），Tab 显示宽度 = 4
- 行长度上限以项目 lint 配置为准（typical: 120-200）
- 函数长度上限以项目 lint 配置为准（typical: 80-300）
- 圈复杂度上限以项目 lint 配置为准（typical: 15-30）

## 测试

- **目录**：测试文件放在 **`tests/`** 下，目录结构与源代码**镜像对应**，**不要**与 `.go` 源文件同目录并列。
  - 例：`internal/session/graph.go` → `tests/internal/session/graph_test.go`
  - 例：`conf/config.go` → `tests/conf/config_test.go`
- **包名**：对外部包测试使用 `package <pkg>_test`（如 `session_test`），`import` 被测包；仅测未导出细节时可在 `tests/...` 下同包，但须单独子目录且不与 `internal/` 源混放。
- 文件命名仍用 `*_test.go`；表驱动测试优先。
- 集成 / E2E：`tests/integration/` 或 `tests/e2e/`（可再加 build tag `//go:build integration`）。
- 运行：`go test ./...` 仍从模块根（如 `backend/`）执行，须能收集到 `tests/**` 下用例。
