---
description: harness 仓库 Git 分支与提交规范（GitHub / master）
paths:
  - "cli/**"
  - "backend/**"
---

# Git 规范（harness）

适用于本仓 **`cli/`** 与未来 **`backend/`** 代码变更。仅改 `.scratch/`、`.cursor/rules/` 等文档时，不必强制拉 feature 分支，除非用户要求。

## 重要原则

- **不要自动提交 git 代码**，除非用户明确指示
- 提交前在对应目录跑验证（如 `cli/`：`pnpm test`、`pnpm type-check`）
- 提交信息简洁，变更尽量小且相关
- Windows 终端勿用 `&&` 连接命令，见 `windows-shell.mdc`

## 仓库事实

- **默认分支：** `master`（`origin` 为 GitHub：`github.com/huodaoshi/harness`）
- **无 `develop` 分支** — 勿写 `git checkout develop`

---

## 分支管理

### 命名规范

| 分支类型 | 命名格式 | 示例 |
| -------- | -------- | ---- |
| 功能分支 | `<用户名>/feat_<短语描述>` | `shimin/feat_session-stream` |
| 修复分支 | `<用户名>/fix_<短语描述>` | `shimin/fix_sse-close` |
| 发布分支 | `release/<版本>` | `release/v1.2.0` |
| 热修复 | `hotfix/<版本>-<描述>` | `hotfix/v1.1.1-crash` |

- `<用户名>`：`git config user.name`，去空格转小写
- `<短语描述>`：简短英文，`-` 连接

### 开发前分支检查

开始 **cli/ 或 backend/** 开发前：

1. `git branch --show-current`
2. 若在 **`master` 或 `main`** 上 → 从最新默认分支创建开发分支：
   ```text
   git checkout master
   git pull
   git checkout -b '<用户名>/feat_<短语描述>'
   ```
   向用户确认已切换后再改代码。
3. 若已在 `*/feat_*` / `*/fix_*` 等开发分支 → 直接工作

---

## 提交规范

### Commit Message

模板：`<type>(<scope>): <subject>`

- 冒号后有空格
- **scope** 常用：`cli`、`backend`、`docs`、`skills`

| type | 含义 |
| ---- | ---- |
| feat | 新功能 |
| fix | 修 bug |
| docs | 文档 |
| style | 格式（不改行为） |
| refactor | 重构 |
| perf | 性能 |
| test | 测试 |
| chore | 构建/工具 |
| revert | 回退 |
| build | 打包 |

多要点时用列表：

```text
feat(cli): add experimental sync flag

- Add --dry-run to sync command
- Update tests for node_modules scan
```

---

## 提交与推送（仅当用户明确要求）

### 1. 提交

- `git add` **仅**本次相关文件（避免无脑 `git add -A`）
- `git commit` 使用上述格式

### 2. 推送

- `git push -u origin <当前分支名>`

### 3. 创建 PR（GitHub）

1. `git remote get-url origin` — 含 `github.com` 时用 **`gh pr create`**
2. **合并目标：** `master`（以仓库默认分支为准）
3. PR 标题 ≈ commit subject；正文含：概述、主要文件、验证命令与结果

```text
gh pr create --base master --title "<subject>" --body "..."
```

4. 将 **PR URL** 返回给用户

**GitLab / 自建：** 使用 push 输出中的 MR 链接，或让用户在浏览器创建（本仓 primary 为 GitHub）。
