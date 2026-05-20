# skills CLI：规则安装扩展

## 背景

在保留现有「技能」安装能力的前提下，让同一套 `skills` CLI 能从 Git 规则包将 **Cursor / Claude Code** 各自规则安装到项目约定目录，并支持 **`rules-lock.json`** 复现。

## 决策依据

见 `docs/adr/0001-skills-cli-rules-install.md`（独立锁、分端目录、`rules/cursor`与`rules/claude`、MVP 不做 list/remove）。

## 范围

- **做（MVP）：** `skills rules add`、基于 `rules-lock.json` 的恢复命令、`-a` 单端、`--to` 覆盖、与 `skills add` 同源的 clone/安装语义（含 symlink/`--copy`）。
- **不做（本阶段）：** `rules list` / `rules remove`、单次命令双 Agent、与 `skills-lock.json` 合并。

## Issue 索引

见 `issues/` 下编号文件；实现顺序建议按编号升序。
