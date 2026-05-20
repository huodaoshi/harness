import { readFile, writeFile } from 'fs/promises';
import { join } from 'path';

export const RULES_LOCK_FILE = 'rules-lock.json';
const CURRENT_VERSION = 1;

export type RulesLockAgent = 'cursor' | 'claude-code';

export interface RulesLockInstall {
  agent: RulesLockAgent;
  /** 安装时使用的来源（与 CLI 参数一致，如 owner/repo 或本地绝对路径） */
  source: string;
  ref?: string;
  sourceType: string;
  /** 预留；安装实现以递归复制为主 */
  copy?: boolean;
}

export interface RulesLockFile {
  version: number;
  /** MVP：单条活跃安装记录；再次 add 时覆盖 */
  install: RulesLockInstall | null;
}

export function getRulesLockPath(cwd?: string): string {
  return join(cwd || process.cwd(), RULES_LOCK_FILE);
}

export function createEmptyRulesLock(): RulesLockFile {
  return { version: CURRENT_VERSION, install: null };
}

export async function readRulesLock(cwd?: string): Promise<RulesLockFile> {
  const lockPath = getRulesLockPath(cwd);

  try {
    const content = await readFile(lockPath, 'utf-8');
    const parsed = JSON.parse(content) as RulesLockFile;

    if (typeof parsed.version !== 'number' || parsed.version !== CURRENT_VERSION) {
      return createEmptyRulesLock();
    }

    if (parsed.install !== null && typeof parsed.install !== 'object') {
      return createEmptyRulesLock();
    }

    return parsed;
  } catch {
    return createEmptyRulesLock();
  }
}

export async function writeRulesLock(lock: RulesLockFile, cwd?: string): Promise<void> {
  const lockPath = getRulesLockPath(cwd);
  const normalized: RulesLockFile = {
    version: CURRENT_VERSION,
    install: lock.install,
  };
  const content = JSON.stringify(normalized, null, 2) + '\n';
  await writeFile(lockPath, content, 'utf-8');
}
