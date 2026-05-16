# AGENTS.md

本文件为在 `skills` CLI 代码库上工作的 AI 编程助手提供指引。

## 项目概览

`skills` 是开放 Agent 技能生态系统的命令行工具。

## 命令

| 命令                          | 说明                                           |
| ----------------------------- | ---------------------------------------------- |
| `skills`                      | 显示横幅与可用命令                             |
| `skills add <pkg>`            | 从 Git 仓库、URL 或本地路径安装技能            |
| `skills experimental_install` | 从 skills-lock.json 恢复技能                   |
| `skills experimental_sync`    | 将 node_modules 中的技能同步到各 Agent 目录    |
| `skills list`                 | 列出已安装技能（别名：`ls`）                   |
| `skills update [skills...]`   | 将技能更新到最新版本                           |
| `skills init [name]`          | 创建新的 SKILL.md 模板                         |

别名：`skills a` 等同于 `add`。`skills i`、`skills install`（无参数）从 `skills-lock.json` 恢复。`skills ls` 等同于 `list`。`skills experimental_install` 从 `skills-lock.json` 恢复。`skills experimental_sync` 扫描 `node_modules` 中的技能。

## 架构

```
src/
├── cli.ts           # 主入口：命令路由、init/check/update
├── add.ts           # add 命令核心逻辑
├── constants.ts      # 共享常量
├── find.ts           # find/搜索命令
├── list.ts          # 列出已安装技能
├── remove.ts         # remove 命令实现
├── agents.ts        # Agent 定义与检测
├── installer.ts     # 技能安装（符号链接/复制）+ listInstalledSkills
├── skills.ts        # 技能发现与解析
├── skill-lock.ts    # 全局锁文件（~/.agents/.skill-lock.json）
├── local-lock.ts    # 项目锁文件（skills-lock.json，可提交）
├── sync.ts          # sync 命令：扫描 node_modules 中的技能
├── source-parser.ts # 解析 git URL、GitHub 简写、本地路径
├── git.ts           # Git 克隆
├── telemetry.ts     # 匿名使用统计
├── types.ts         # TypeScript 类型
├── mintlify.ts      # Mintlify 技能拉取（遗留）
├── plugin-manifest.ts # 插件清单发现
├── prompts/         # 交互式提示辅助
│   └── search-multiselect.ts
├── providers/       # 远程技能提供方（GitHub、HuggingFace、Mintlify）
│   ├── index.ts
│   ├── registry.ts
│   ├── types.ts
│   ├── huggingface.ts
│   ├── mintlify.ts
│   └── wellknown.ts
tests/
├── test-utils.ts            # CLI 子进程测试辅助
├── add.test.ts              # add 命令测试
├── add-prompt.test.ts       # add 提示行为测试
├── cli.test.ts              # CLI 测试
├── init.test.ts             # init 命令测试
├── list.test.ts             # list 命令测试
├── remove.test.ts           # remove 命令测试
├── update-source.test.ts    # 更新 source 格式化测试
├── source-parser-gitlab.test.ts # GitLab / git URL 解析测试
├── cross-platform-paths.test.ts # 跨平台路径规范化
├── full-depth-discovery.test.ts # --full-depth 技能发现测试
├── openclaw-paths.test.ts       # OpenClaw 专用路径测试
├── plugin-manifest-discovery.test.ts # 插件清单技能发现
├── sanitize-name.test.ts     # sanitizeName 测试（防路径遍历）
├── skill-matching.test.ts    # filterSkills 测试（多词技能名匹配）
├── source-parser.test.ts     # URL/路径解析测试
├── installer-symlink.test.ts # 符号链接安装测试
├── list-installed.test.ts    # 已安装技能列表测试
├── skill-path.test.ts        # 技能路径处理测试
├── wellknown-provider.test.ts # well-known 提供方测试
├── xdg-config-paths.test.ts   # XDG 全局路径测试
└── dist.test.ts               # 构建产物分发测试
```

## 更新检查机制

### `skills check` 与 `skills update` 如何工作

1. 读取 `~/.agents/.skill-lock.json` 中的已安装技能
2. 筛选 GitHub 来源且同时具备 `skillFolderHash` 与 `skillPath` 的技能
3. 对每个技能调用 `fetchSkillFolderHash(source, skillPath, token)`。可选认证令牌来自 `GITHUB_TOKEN`、`GH_TOKEN` 或 `gh auth token`，用于提高速率限制
4. `fetchSkillFolderHash` 直接调用 GitHub Trees API（先 `main` 分支 `/git/trees/<branch>?recursive=1`，失败则回退 `master`）
5. 将最新目录树 SHA 与锁文件中的 `skillFolderHash` 比较；不一致表示有更新
6. `skills update` 通过直接调用当前 CLI 入口（`node <repo>/bin/cli.mjs add <source-tree-url> -g -y`）重装变更的技能，避免嵌套 npm exec/npx

### 锁文件兼容性

锁文件格式为 v3。关键字段：`skillFolderHash`（技能目录的 GitHub tree SHA）。

若读取到旧版锁文件，会被清空。用户需重新安装技能以填充新格式。

## 关键集成点

| 功能                       | 实现                                                          |
| -------------------------- | ------------------------------------------------------------- |
| `skills add`               | `src/add.ts` — 完整实现                                       |
| `skills experimental_sync` | `src/sync.ts` — 扫描 node_modules                             |
| `skills check`             | `src/cli.ts` + `src/skill-lock.ts` 中的 `fetchSkillFolderHash` |
| `skills update`            | `src/cli.ts` 直接比较 hash + 通过 `skills add` 重装           |

## 开发

```bash
# 安装依赖
pnpm install

# 构建
pnpm build

# 本地测试
pnpm dev add vercel-labs/agent-skills --list
pnpm dev experimental_sync
pnpm dev check
pnpm dev update
pnpm dev init my-skill

# 运行全部测试
pnpm test

# 运行指定测试文件
pnpm test tests/sanitize-name.test.ts
pnpm test tests/skill-matching.test.ts tests/source-parser.test.ts

# 类型检查
pnpm type-check

# 格式化代码
pnpm format

# 仅检查格式
pnpm format:check

# 校验并同步 Agent 元数据/文档
pnpm run -C scripts validate-agents.ts
pnpm run -C scripts sync-agents.ts
```

## 代码风格

本项目使用 Prettier 格式化代码。**提交前请运行 `pnpm format`**，以保持格式一致。

```bash
# 格式化所有文件
pnpm format

# 仅检查、不修改
pnpm format:check
```

CI 会在格式不正确时失败。

## 发布

```bash
# 1. 在 package.json 中 bump 版本
# 2. 构建
pnpm build
# 3. 发布
npm publish
```

## 添加新 Agent

1. 在 `src/agents.ts` 中添加 Agent 定义
2. 运行 `pnpm run -C scripts validate-agents.ts` 校验
3. 运行 `pnpm run -C scripts sync-agents.ts` 更新 README.md 与 package keywords
