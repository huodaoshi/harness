# `skills rules`：项目规则包与锁文件

本页说明 `cli/`（`skills` 命令）中 **`skills rules add`** 与 **`skills rules experimental_install`** 的约定，对应决策见 `docs/adr/0001-skills-cli-rules-install.md`。

## 规则包目录（Git 仓库内）

安装时按 **`-a` 选择的编辑器** 只读取对应子目录：

| `-a` 参数（Agent） | 仓库内路径 | 安装到项目内默认目录（未使用 `--to` 时） |
|-------------------|------------|----------------------------------------|
| `cursor` | `rules/cursor/` | `.cursor/rules/` |
| `claude-code` | `rules/claude/` | `.claude/rules/` |

- 子目录内允许任意层级（例如 `_lang/go.mdc`）；安装时**整棵子树**复制到目标目录。
- **不要**在安装时自动改扩展名（`.md` ↔ `.mdc`）；分端两套文件由包维护者分别维护。
- `--to <path>`：将**该子树以内**的文件铺到 `<path>`，**不**再套一层 `rules/cursor`。

MVP **仅**支持上表两个 Agent；其他 Agent 即使存在于 `agents` 定义中也不接受 `skills rules`。

## `rules-lock.json`（项目根）

与 `skills-lock.json` **分立**；用于从固定来源恢复规则目录。

- **路径：** 项目根目录 `rules-lock.json`（与 `skills-lock.json` 同级）。
- **MVP 形态：** 单次安装只保留**一条**活跃记录 `install`；再次 `rules add` 会覆盖该记录。

### 字段（草案）

| 字段 | 类型 | 说明 |
|------|------|------|
| `version` | `number` | 固定为 `1` |
| `install` | `object \| null` | 无安装或未写锁时为 `null` |
| `install.agent` | `"cursor" \| "claude-code"` | 与 `-a` 一致 |
| `install.source` | `string` | 安装时使用的来源字符串（如 `owner/repo` 或本地路径） |
| `install.ref` | `string`（可选） | Git ref / tag / branch |
| `install.sourceType` | `string` | 与内部 `ParsedSource.type` 一致（如 `github`、`local`） |
| `install.copy` | `boolean`（可选） | 预留；当前实现以**递归复制**为主 |

恢复命令 **`skills rules experimental_install`** 只读此文件，**不**读写 `skills-lock.json`。

---

## 安装到其它项目（示例 `D:\one-eino`）与 GitHub 源

- **目标项目不是 harness 时**：在 **目标仓库根**打开终端（如 `cd D:\one-eino`），再执行  
  `skills rules add <来源> -a cursor|claude-code`，规则会写入**当前目录**下的 `.cursor/rules` / `.claude/rules`，并在该目录生成 **`rules-lock.json`**。  
  来源可为 **本地 harness 路径**（如 `D:\harness`）或 **GitHub `owner/repo`**（仓库内需含 `rules/cursor/`、`rules/claude/`）。  
  详见仓库内 **`rules/README.md`**（含 `pnpm --dir D:\harness\cli dev ...` 示例）。

---

## 非目标（当前版本）

- `rules list` / `rules remove`
- 单次命令同时为两个 Agent 安装
- 与 `skills-lock.json` 合并
