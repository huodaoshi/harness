# 仓库内技能源（`skills/`）

本目录用于放 **可由 `skills add` 从本地路径安装的**技能包（每个技能为独立子目录，内含 **`SKILL.md`**）。

已将 **`.agents/skills/`** 下各技能目录**复制**到此处，便于把同一套内容当「可安装的包」使用（与 `.agents/skills` 并行；日后以一侧为权威时请自行约定并去重）。

与 **`.agents/skills/`** 的关系（复制后）：

| 位置 | 典型用途 |
|------|----------|
| **`skills/` 本目录** | 与 `skills add <本地路径>`、拆分发行为独立仓库时的**工作副本** |
| **`.agents/skills/`** | 原先 harness 维护位置；Cursor / Agent 仍可能直接读此路径（视工具配置） |

## 子目录

各子目录均来自 **`.agents/skills/`** 的拷贝（如 `caveman`、`triage`、`grill-with-docs` 等）。

## 从本仓库根安装到当前项目

在仓库根：

```powershell
# 列出本仓库作为「包」时可发现的技能
pnpm --dir cli dev add . --list

# 安装指定技能（示例：caveman）
pnpm --dir cli dev add . --skill caveman -y
```

更多见 **`cli/README.md`**、官方 `skills` 文档。

---

## 安装到其它项目（目标目录示例 `D:\one-eino`）

在 **要接收技能的仓库里** 打开终端（**当前目录 = 目标项目根**），用 **harness 里的 CLI** 指向本仓库作为**技能包来源**。

**前置：** 已克隆 harness 到本机，例如 **`D:\harness`**（以下路径请按你机器修改）。

```powershell
cd D:\one-eino

# 列出 harness 仓库作为「包」时能发现的技能（不安装）
pnpm --dir D:\harness\cli dev add D:\harness --list

# 安装若干技能到当前项目（cwd = one-eino，技能装进本项目的 .agents/skills 等约定目录）
pnpm --dir D:\harness\cli dev add D:\harness --skill caveman --skill triage -y
```

说明：

- **第一个参数**是技能包根路径（此处为 **`D:\harness`**，会扫描其中的 `skills/` 与各 `SKILL.md`）。
- **`cd` 到 `D:\one-eino`** 再执行，安装目标才是 **one-eino**，而不是 harness。
- 若已将 `skills` CLI **全局 / 发布为 npm 包**，也可使用 `npx skills add D:\harness ...`（与上面等价，视你环境选择）。

---

## 从 GitHub 源安装（不上传 harness 整仓到 one-eino 时）

当技能包已推送到 GitHub，用 **`owner/repo`** 或 **HTTPS / SSH URL** 代替本地路径 `D:\harness`。

示例（将 `huodaoshi/harness` 换成你的仓库；**仍在目标项目根执行**）：

```powershell
cd D:\one-eino

pnpm --dir D:\harness\cli dev add huodaoshi/harness --skill caveman -y

# 指定分支 / 标签（fragment）
pnpm --dir D:\harness\cli dev add huodaoshi/harness#main --skill caveman -y
```

私有仓库需本机已配置 **`gh auth`**、**SSH key** 或 **HTTPS 凭据**，与平常 `git clone` 一致。

