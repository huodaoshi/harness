---
description: 跨语言命名——JSON/Mongo snake_case，Go/TS 各自 camelCase 约定
paths:
  - "cli/**"
  - "backend/**"
  - "**/*.go"
---

# 跨语言命名规范（harness）

## 原则

**同一概念在不同层用该层惯用写法**；边界由序列化 / DTO 映射完成，不要强行统一成一种拼写。

| 层 / 语言 | 约定 | 示例 |
| --------- | ---- | ---- |
| JSON API / MongoDB / YAML | `snake_case` | `session_id`、`user_id` |
| Go | `camelCase`，缩写全大写 | `sessionID`、`userID` |
| TypeScript（`cli/`） | `camelCase`，缩写首字母大写其余小写 | `sessionId`、`userId` |

## 强制规则

- **Go：** `ID` / `URL` / `UID` 全大写 → `sessionID`，不用 `sessionId` 或 `session_id` 字段名
- **TypeScript：** `sessionId`，不用 `sessionID` 或 `session_id` 变量名
- **JSON / MongoDB：** 键名一律 `snake_case`

### JSON 负载（API、对外配置）

跨进程、跨语言的 JSON **键名一律 `snake_case`**（与 Mongo 文档一致）。

- **Go（`backend/`）：** 对外结构体必须写 `json:"snake_case"` tag，禁止用默认 PascalCase 作线上契约
- **TypeScript（`cli/`）：** 对外 API 若后端为 snake_case，在边界映射；仓库内 TS 代码用 `camelCase`（见 `cli/AGENTS.md`）

## 禁止混用

```go
// ❌ Go 里 snake_case 或 TS 风格
session_id string
SessionId  string
// ✅
SessionID  string `json:"session_id"`
```

```typescript
// ❌ TS 里 Go 风格或 snake_case
const sessionID = '';
const session_id = '';
// ✅
const sessionId = '';
```
