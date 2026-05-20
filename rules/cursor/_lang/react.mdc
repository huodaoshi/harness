---
description: React 通用编码规范——函数组件、Hooks、性能、可访问性
paths:
  - "**/*.tsx"
  - "**/*.jsx"
---

# React 编码规范

## 组件

- 函数组件优先；不使用 `class` 组件（除非边界场景如 ErrorBoundary）
- 函数组件使用 `const` 箭头函数 + 默认导出（按项目约定）
- 组件 Props 用 `interface` 或 `type` 显式声明
- 一个文件一个组件；强相关的小组件可同文件，但导出一个主组件

## Hooks 规则

- Hook 调用必须在组件顶层，不在条件 / 循环中调用
- 自定义 Hook 命名 `use*`
- 副作用放在 `useEffect`，不在 render 期间执行 IO
- 依赖数组完整声明，不使用 `// eslint-disable-next-line react-hooks/exhaustive-deps`

## 性能

- 列表渲染必须设置稳定的 `key`（不用 index）
- 大组件用 `React.memo` 包裹（仅在 profiling 后）
- 稳定函数引用：用 `useCallback` 或项目封装的等价 hook（如 ahooks `useMemoizedFn`）
- 派生状态用 `useMemo`，不在 render 中重复计算

## 状态管理

- 局部状态用 `useState` / `useReducer`
- 共享状态按项目选用：Context / Zustand / Redux / Jotai 等（按 skills 中描述）
- 服务端状态用专用库（React Query / SWR / ahooks `useRequest`，按 skills 中描述）

## 可访问性（A11y）

- 交互元素用语义化标签（`button` / `a` / `label`）
- 表单控件关联 `label`
- 图片必须 `alt`
- 键盘可达性：`tabIndex`、`onKeyDown` 配合 `onClick`

## 样式

- 不裸拼接 className 字符串；用 `clsx` / `cn` 工具函数
- CSS 方案按项目约定（Tailwind / CSS Modules / Emotion / styled-components）

## 文件结构

- 组件文件 PascalCase（如 `UserCard.tsx`）
- Hooks 文件 `use-*.ts` 或 `useXxx.ts`（按项目约定）
- 路由懒加载：`React.lazy` + `Suspense`
