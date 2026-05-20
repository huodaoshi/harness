---
description: Vue 通用编码规范——组合式 API、SFC、性能、可访问性
paths:
  - "**/*.vue"
---

# Vue 编码规范

## 组件

- Vue 3 + 组合式 API（`<script setup>`）优先
- SFC 文件命名 PascalCase（`UserCard.vue`）
- Props / Emits 使用 TypeScript 类型声明
- 单文件组件结构顺序：`<script setup>` → `<template>` → `<style>`（若存在）  
- **`frontend/` 例外**：当 `.claude/rules/frontend-h5-ui.md` 命中时，**业务样式不得写在 SFC 的 `<style>` 中**，须独立 `.scss`/`.css` 并在 `<script setup>` 中 `import`；详见该规则 **§3.1**。此类 SFC 可无 `<style>` 块。

## 组合式 API

- 副作用放在 `onMounted` / `onUnmounted` / `watchEffect`
- 派生状态用 `computed`
- 自定义 hook 命名 `use*`，放在 `composables/` 目录
- ref / reactive 的选择：基础值用 `ref`，对象用 `reactive`（按项目约定可统一为 `ref`）

## 模板

- 列表渲染必须 `:key`，不用 index
- 避免 `v-if` 与 `v-for` 同元素
- 复杂表达式提取为 `computed`
- 事件命名 `kebab-case`（emit）/ `camelCase`（props）

## 状态管理

- 局部状态用组合式 API
- 跨组件状态用 Pinia（Vue 3 推荐），不用 Vuex
- 服务端状态用专用库（按 skills 中描述）

## 性能

- 大组件用 `defineAsyncComponent` 异步加载
- 列表大数据用虚拟滚动
- v-show vs v-if：频繁切换用 `v-show`，条件初始化用 `v-if`

## 样式

- **默认**（非 `frontend/` 或未受 `frontend-h5-ui` 约束的路径）：`<style scoped>`；全局样式独立文件。  
- **`frontend/`**：以 **`frontend-h5-ui.md`** 为准——**禁止**在 `.vue` 里写业务 `<style>`，须独立样式文件 + `import`。  
- 不在 `<style>` 中使用 `!important`（除非边界场景）；**外置样式文件**同样遵守 **frontend-h5-ui** 与 **本节** 对 `!important` 的约束。  
- CSS 方案按项目约定  
- **SFC 内 `<style>` / `<style scoped>` / `<style lang="scss">`（及 `less` 等）**：选择器与声明块的**每一级嵌套均为 4 个空格缩进**，**不使用 Tab**；与 **§ 格式化** 及工程内 Prettier（如 `matechat-h5/.prettierrc.cjs` 的 `tabWidth: 4`）一致。

## 可访问性

- 交互元素用语义化标签
- 表单控件关联 `label`
- 图片必须 `alt`

## 格式化

- **缩进**：**每级 4 个空格**（**不使用 Tab**），与 **`typescript.md` §Lint/格式化** 一致。  
- **适用范围**：同一 `.vue` 文件内的 **`<template>`**、**`<script setup>`**、**`<style>`**（含 `scoped`、`lang="scss"` / `less`、第二段非 scoped `<style>`）均遵守上述缩进；嵌套越深，每多一级增加 **4 空格**，禁止 2 空格与 Tab 混用。
