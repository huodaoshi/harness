# 仓库内规则源（`rules/`）

本目录是 **Cursor / Claude Code 编辑器规则的规范源**（分端维护后缀：`.mdc` / `.md`）。  
安装目标仍为项目下的 `.cursor/rules/` 与 `.claude/rules/`（或由 `skills rules add --to` 指定）。

## 目录

| 路径 | 用途 |
|------|------|
| `rules/cursor/` | Cursor 规则（`*.mdc` 等），对应安装到 **`.cursor/rules/`** |
| `rules/claude/` | Claude Code 规则（`*.md` 等），对应安装到 **`.claude/rules/`** |

## 日常 workflow

1. **只改 `rules/cursor/` / `rules/claude/`**（规范源），不要两边各改一套导致漂移。
2. 安装到编辑器目录（在仓库根执行）：

```powershell
pnpm --dir cli dev rules add . -a cursor
pnpm --dir cli dev rules add . -a claude-code
```

3. 可提交 **`rules-lock.json`**（由 `rules add` 写入），以便他人用 `pnpm --dir cli dev rules experimental_install` 复现。

## 与 `.cursor/rules`、`.claude/rules` 的关系

- 当前仓库已用 **本目录内容** /bootstrap 初始一致；之后以 **`rules/` 为编辑源头**，需要同步到 IDE 时再跑上面的 `rules add .`。
- 若你更习惯**只保留一端**：也可以只维护 `rules/`，将 `.cursor/rules` 视为「生成物」并在文档里说明（团队需统一约定）。

详见 **`docs/cli-rules.md`**、**`docs/adr/0001-skills-cli-rules-install.md`**。

---

## 安装到其它项目（目标目录示例 `D:\one-eino`）

在 **要安装规则的仓库** 中，终端 **当前工作目录 = 该仓库根**（例如先 **`cd D:\one-eino`**），再指向 **harness** 作为规则包来源。

**前置：** harness 在本机路径 **`D:\harness`**（按你实际路径修改）。

```powershell
cd D:\one-eino

# 把 harness 仓库里的 rules/cursor 整树复制到 one-eino\.cursor\rules
pnpm --dir D:\harness\cli dev rules add D:\harness -a cursor

# 把 rules/claude 复制到 one-eino\.claude\rules
pnpm --dir D:\harness\cli dev rules add D:\harness -a claude-code
```

说明：**不写 `--to` 时**，规则写入 **当前 cwd** 下的 `.cursor/rules` / `.claude/rules`，因此必须先 **`cd` 到 one-eino**。

在 one-eino 根会生成 / 更新 **`rules-lock.json`**（记录来源为 `D:\harness`）；他人可在 one-eino 根执行 **`pnpm --dir D:\harness\cli dev rules experimental_install`** 复现（需能访问同一来源）。

---

## 从 GitHub 源安装规则

harness 已推送 GitHub 后，可把 **`D:\harness`** 换成远程（仍在 **one-eino 根** 执行示例）：

```powershell
cd D:\one-eino

pnpm --dir D:\harness\cli dev rules add huodaoshi/harness -a cursor
pnpm --dir D:\harness\cli dev rules add huodaoshi/harness -a claude-code

# 固定分支或标签
pnpm --dir D:\harness\cli dev rules add huodaoshi/harness#main -a cursor
```

将 **`huodaoshi/harness`** 换成你的实际 **`owner/repo`**。包内仍须包含 **`rules/cursor/`** 与 **`rules/claude/`** 子目录。

