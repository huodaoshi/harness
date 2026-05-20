---
description: TypeScript 通用编码规范——严格类型、模块化、错误处理、命名
paths:
  - "**/*.ts"
  - "**/*.tsx"
---

# TypeScript 编码规范

## 严格模式

- `tsconfig.json` 启用 `strict: true`（含 `strictNullChecks`、`noImplicitAny` 等）
- 不使用 `any`，需要逃逸时用 `unknown` 并显式断言
- 不使用 `@ts-ignore` / `@ts-nocheck`；确需要时用 `@ts-expect-error` 并加注释说明原因

## 类型声明

- 公共 API 必须有显式返回类型
- 接口与类型别名：领域模型用 `interface`（可扩展），联合 / 交叉 / mapped 用 `type`
- 枚举：优先 `const` 对象 + `as const` 推导，避免 TypeScript `enum`（除非有运行时反向映射需求）
- **领域语义数值（角色、状态位等）**：**禁止**使用无命名的数字字面量联合类型（如 `role: 1 | 2`、`kind: 0 | 1`）作为对外 props、store 模型或公共 API 的类型。必须在 `src/constants/`（或该域专用 `*.ts`）声明 **`const Xxx = { … } as const`**，并导出 **`type XxxId = (typeof Xxx)[keyof typeof Xxx]`**（命名与 `naming.md` 一致：`MessageRoleId`、`ScanStatusId` 等）；业务代码只引用 **`Xxx.User`** 等与 **`XxxId`**，禁止在业务处散落裸 `1` / `2` 表示语义。
- 不导出的内部类型放在使用处，导出类型放在 `types.ts` / `types/` 目录

## 模块化

- 使用 ESM 语法（`import` / `export`）
- 不使用默认导出 + 命名导出混用（按项目约定二选一，团队内一致）
- 路径别名（如 `@/`）从 `tsconfig.json` 的 `paths` 配置，不裸 `../../../`

## 错误处理

- `try/catch` 中不吞错误；至少 log + rethrow 或转为业务错误
- 异步代码用 `async/await`，避免 Promise then-chain 嵌套
- 不在 Promise 中漏掉 `.catch()`（`no-floating-promises`）

## 命名

- 类型 / 接口：`PascalCase`
- 变量 / 函数：`camelCase`
- 常量：`UPPER_SNAKE_CASE`（仅模块级真正常量）
- React 组件文件：`PascalCase.tsx`；其他文件：`kebab-case.ts` 或 `camelCase.ts`（按项目约定）
- 缩写词首字母大写、其余小写：`sessionId`、`spaceId`、`userId`——不写成 `sessionID`、`spaceID`（Go 风格）或 `session_id`（snake_case）；见 `rules/naming.md`

## 不可变性

- 默认 `const`，仅在需要重新赋值时 `let`
- 不使用 `var`
- 数组 / 对象的浅拷贝用 `[...arr]` / `{...obj}`，避免 mutation

## Lint / 格式化

- 使用 ESLint 或 Biome（按项目约定）
- 使用 Prettier 或 Biome 格式化（按项目约定）
- 缩进：**4 空格**（不使用 Tab）
- import 排序：标准库 / 第三方 / 路径别名 / 相对路径，每组之间空行
