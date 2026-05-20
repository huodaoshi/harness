---
description: Rust 通用编码规范——所有权、错误处理、命名、并发
paths:
  - "**/*.rs"
---

# Rust 编码规范

## 风格

- 使用 `rustfmt` 默认配置
- 行长度上限以 `rustfmt.toml` 为准（默认 100）
- 文件编码 UTF-8

## 所有权与借用

- 优先借用而非所有权转移
- 不必要的 `clone()` 视为代码异味
- 生命周期标注：能省略则省略；需要时显式且最小化

## 错误处理

- 库代码用自定义 `Error` 枚举（配合 `thiserror`）
- 应用代码用 `anyhow::Result`（按项目约定）
- 不在库中 `panic!` / `unwrap()` / `expect()`（除非确认不可能）
- 错误传播用 `?`

## 命名

- 模块 / crate：`snake_case`
- 类型 / trait：`PascalCase`
- 函数 / 变量：`snake_case`
- 常量 / 静态：`UPPER_SNAKE_CASE`
- 文件名：`snake_case.rs`

## 模块化

- 公共 API 通过 `pub` 显式导出
- 一个模块单一职责
- `mod.rs` 与 `module/` 目录一致

## 并发

- 优先 `tokio` / `async-std`（按项目约定）
- `Send` / `Sync` 边界明确
- 共享状态用 `Arc<Mutex<T>>` / `RwLock`，但优先 message passing

## Trait

- 公共 trait 应实现常见标准 trait（`Debug`, `Clone`, `PartialEq` 视情况）
- 不为外部类型实现外部 trait（孤儿规则）
- 使用 `#[derive(...)]` 减少样板代码

## 测试

- 单元测试放在 `#[cfg(test)] mod tests`
- 集成测试放在 `tests/` 目录
- 文档测试用三反引号代码块

## 工具

- 格式化：`cargo fmt`
- Lint：`cargo clippy --all-targets -- -D warnings`
- 测试：`cargo test`
- 构建：`cargo build` / `cargo build --release`
