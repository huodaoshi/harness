---
name: code-review
description: |
  为 React 19、Vue 3、Angular 17+、Svelte 5、Rust、TypeScript、Java、Python、Django、Go、C#/.NET、Kotlin、NestJS、C/C++ 等提供全面代码审查指引。
  帮助发现 bug、提升代码质量、给出建设性反馈。
  Use when: reviewing pull requests, conducting PR reviews, code review, reviewing code changes,
  establishing review standards, mentoring developers, architecture reviews, security audits,
  checking code quality, finding bugs, giving feedback on code,
  审查 PR、代码审查、架构审查、安全审计、代码质量检查、找 bug、代码反馈。
allowed-tools:
  - Read
  - Grep
  - Glob
  - Bash      # 运行 lint/test/build 命令验证代码质量
  - WebFetch  # 查阅最新文档和最佳实践
---

# 卓越代码审查

通过建设性反馈、系统分析与协作改进，把 code review 从「挡路」变成「知识共享」。

## 何时使用本技能

- 审查 pull request 与代码变更
- 为团队建立 code review 标准
- 通过 review 带 junior 开发者
- 做架构审查
- 编写 review checklist 与指南
- 改善团队协作
- 缩短 review 周期
- 维持代码质量标准

## 核心原则

### 1. Review 心态

**Code review 的目标：**
- 发现 bug 与边界情况
- 保证可维护性
- 在团队内共享知识
- 执行编码标准
- 改进设计与架构
- 建设团队文化

**不是目标：**
- 炫技
- 抠格式（用 linter）
- 无故阻塞进度
- 按你的偏好重写

### 2. 有效反馈

**好反馈应：**
- 具体、可执行
- 教育性，非评判性
- 针对代码，非针对人
- 平衡（也表扬做得好的）
- 分优先级（critical vs nice-to-have）

```markdown
❌ 差："This is wrong."
✅ 好："多用户同时访问可能竞态。此处 Consider using a mutex."

❌ 差："Why didn't you use X pattern?"
✅ 好："是否考虑过 Repository pattern？更易测。示例：[link]"

❌ 差："Rename this variable."
✅ 好："[nit] Consider `userCount` instead of `uc` for clarity. 不阻塞，保留也可。"
```

### 3. Review 范围

**要审的：**
- 逻辑正确性与边界情况
- 安全漏洞
- 性能影响
- 测试覆盖与质量
- 错误处理
- 文档与注释
- API 设计与命名
- 架构契合度

**勿手工审的：**
- 代码格式（用 Prettier、Black 等）
- import 排序
- lint 违规
- 简单 typo

## 审查流程

### Phase 1：收集上下文（2–3 分钟）

深入代码前先理解：
1. 读 PR 描述与关联 issue
2. 看 PR 规模（>400 行？要求拆分）
3. 看 CI/CD 状态（测试是否过）
4. 理解业务需求
5. 留意相关架构决策

### Phase 2：高层审查（5–10 分钟）

1. **架构与设计** — 方案是否匹配问题？
   - 重大变更见 [Architecture Review Guide](reference/architecture-review-guide.md)
   - 查：SOLID、耦合/内聚、反模式
2. **性能评估** — 有无性能顾虑？
   - 性能敏感代码见 [Performance Review Guide](reference/performance-review-guide.md)
   - 查：算法复杂度、N+1、内存占用
3. **文件组织** — 新文件位置是否合理？
4. **测试策略** — 有无覆盖边界情况的测试？

### Phase 3：逐行审查（10–20 分钟）

对每个文件检查：
- **逻辑与正确性** — 边界、off-by-one、null、竞态
- **安全** — 输入校验、注入、XSS、敏感数据
- **性能** — N+1、多余循环、内存泄漏
- **可维护性** — 清晰命名、单一职责、注释
- **复用** — 接受新代码前，搜现有 utility/helper 能否替代；查相邻文件与共享模块的类似模式。反模式见 [Universal Quality Guide](reference/code-quality-universal.md)（参数膨胀、leaky abstraction、嵌套条件、stringly-typed、TOCTOU、no-op update 等）

### Phase 4：总结与决策（2–3 分钟）

1. 总结主要顾虑
2. 指出做得好的地方
3. 明确决策：
   - ✅ Approve
   - 💬 Comment（小建议）
   - 🔄 Request Changes（必须处理）
4. 复杂时 offer pair

## 审查技巧

### 技巧 1：Checklist 法

用 checklist 保证一致性。安全清单见 [Security Review Guide](reference/security-review-guide.md)。

### 技巧 2：提问法

少陈述问题，多提问：

```markdown
❌ "This will fail if the list is empty."
✅ "What happens if `items` is an empty array?"

❌ "You need error handling here."
✅ "How should this behave if the API call fails?"
```

### 技巧 3：建议，非命令

用协作语气：

```markdown
❌ "You must change this to use async/await"
✅ "Suggestion: async/await 可能更易读。你觉得呢？"

❌ "Extract this into a function"
✅ "这段逻辑出现 3 次。是否值得抽取？"
```

### 技巧 4：区分严重度

用标签标优先级：

- 🔴 `[blocking]` — 合并前必改
- 🟡 `[important]` — 应改，有分歧可讨论
- 🟢 `[nit]` — 可选，不阻塞
- 💡 `[suggestion]` — 可考虑的替代方案
- 📚 `[learning]` — 教学性评论，无需动作
- 🎉 `[praise]` — 做得好，保持

## 语言专项指南

按审查代码语言查阅对应详细指南：

| 语言/框架 | 参考文件 | 要点 |
|-----------|----------|------|
| **React** | [React Guide](reference/react.md) | Hooks、useEffect、React 19 Actions、RSC、Suspense、TanStack Query v5 |
| **Vue 3** | [Vue Guide](reference/vue.md) | Composition API、响应性系统、Props/Emits、Watchers、Composables |
| **Angular 17+** | [Angular Guide](reference/angular.md) | Signals、Standalone 组件、RxJS、Zoneless 变更检测、模板优化 |
| **Rust** | [Rust Guide](reference/rust.md) | 所有权/借用、Unsafe 审查、异步代码、取消安全性、错误处理 |
| **TypeScript** | [TypeScript Guide](reference/typescript.md) | 类型安全、async/await、不可变性 |
| **Python** | [Python Guide](reference/python.md) | 可变默认参数、异常处理、类属性 |
| **Django / DRF** | [Django Guide](reference/django.md) | 安全审查、N+1 查询、Serializer 反模式、ViewSet、异步视图 |
| **Java** | [Java Guide](reference/java.md) | Java 17/21 新特性、Spring Boot 3、虚拟线程、Stream/Optional |
| **C# / .NET** | [C# Guide](reference/csharp.md) | C# 12 特性、异步编程、EF Core 性能、ASP.NET Core、LINQ |
| **Go** | [Go Guide](reference/go.md) | 错误处理、goroutine/channel、context、接口设计 |
| **Kotlin / Android** | [Kotlin Guide](reference/kotlin.md) | 协程、Flow、Jetpack Compose、空安全、内存泄漏、架构模式 |
| **NestJS** | [NestJS Guide](reference/nestjs.md) | 依赖注入、分层架构、DTO 验证、Guard/Interceptor、循环依赖 |
| **Svelte / SvelteKit** | [Svelte Guide](reference/svelte.md) | Runes、Load 函数、Form Actions、Store 迁移、SSR/CSR 边界 |
| **C** | [C Guide](reference/c.md) | 指针/缓冲区、内存安全、UB、错误处理 |
| **C++** | [C++ Guide](reference/cpp.md) | RAII、生命周期、Rule of 0/3/5、异常安全 |
| **CSS/Less/Sass** | [CSS Guide](reference/css-less-sass.md) | 变量规范、!important、性能优化、响应式、兼容性 |
| **Qt** | [Qt Guide](reference/qt.md) | 对象模型、信号/槽、内存管理、线程安全、性能 |

## 横切指南

与语言无关、适用于所有 review 的模式：

| 主题 | 参考文件 | 要点 |
|------|----------|------|
| **通用质量** | [Universal Quality Guide](reference/code-quality-universal.md) | 复用审计、参数膨胀、leaky abstraction、嵌套条件、stringly-typed、TOCTOU、no-op update、冗余状态 |

## 其他资源

- [Architecture Review Guide](reference/architecture-review-guide.md) — 架构设计审查（SOLID、反模式、耦合度）
- [Performance Review Guide](reference/performance-review-guide.md) — 性能审查（Web Vitals、N+1、复杂度）
- [Common Bugs Checklist](reference/common-bugs-checklist.md) — 按语言分类的常见错误清单
- [Security Review Guide](reference/security-review-guide.md) — 安全审查指南
- [Code Review Best Practices](reference/code-review-best-practices.md) — 代码审查最佳实践
- [PR Review Template](assets/pr-review-template.md) — PR 审查评论模板
- [Review Checklist](assets/review-checklist.md) — 快速参考清单
