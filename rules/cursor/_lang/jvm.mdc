---
description: Java/Kotlin 通用编码规范——OOP、不可变性、错误处理、并发
paths:
  - "**/*.java"
  - "**/*.kt"
---

# JVM (Java / Kotlin) 编码规范

## 风格

- Java：遵循 Google Java Style 或项目约定
- Kotlin：遵循 Kotlin 官方编码规范（`ktlint` / `detekt`）
- 缩进 4 空格（Java）/ 4 空格（Kotlin）
- 行长度上限以项目 lint 配置为准（typical: 100-120）

## 命名

- 包：`lowercase`，反向域名（`com.example.module`）
- 类 / 接口：`PascalCase`
- 方法 / 变量：`camelCase`
- 常量：`UPPER_SNAKE_CASE`
- 文件名：与 public 类名一致

## OOP

- 优先组合而非继承
- 接口与实现分离
- 使用依赖注入（Spring / Guice / Dagger，按项目约定）
- 不导出可变内部状态

## 不可变性

- Java：字段优先 `final`；DTO 用 `record`（17+）
- Kotlin：优先 `val`，仅在需要重新赋值时 `var`
- 集合优先不可变（`List.of()` / `listOf()`）

## 错误处理

- 不吞异常；至少 log + rethrow 或转为业务异常
- Checked vs Unchecked：业务可恢复用 checked，编程错误用 unchecked
- Kotlin：避免抛出 checked exception；用 `Result<T>` / 密封类表达
- 资源管理用 `try-with-resources`（Java）/ `use {}`（Kotlin）

## 并发

- 优先 `CompletableFuture` / Kotlin Coroutines（按项目约定）
- 不直接创建 `Thread`，用 `ExecutorService`
- 共享状态用 `ConcurrentHashMap` / `AtomicXxx` / `synchronized`

## 测试

- 单元测试用 JUnit 5 + Mockito / MockK
- 测试方法命名清晰描述场景与期望
- 集成测试用专用 profile / source set

## 工具

- 构建：Maven / Gradle（按项目约定）
- 格式化：`google-java-format` / `ktlint`
- Lint：`checkstyle` / `spotbugs` / `detekt`
- 依赖管理：BOM / version catalog
