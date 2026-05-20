import { resolve } from 'path';

/**
 * CLI 安装类子命令的「目标项目根」。
 * 顺序：显式 `--cwd` → npm/pnpm 的 INIT_CWD（调用方目录）→ process.cwd()
 *
 * 使用 `pnpm --dir path/to/cli dev …` 时 Node 的 cwd 常为 cli 包目录，INIT_CWD 仍可为业务仓库根。
 */
export function resolveCliProjectRoot(explicitCwd?: string): string {
  const trimmed = explicitCwd?.trim();
  if (trimmed) {
    return resolve(trimmed);
  }
  const init = process.env.INIT_CWD?.trim();
  if (init) {
    return resolve(init);
  }
  return process.cwd();
}
