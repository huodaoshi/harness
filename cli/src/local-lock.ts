import { readFile, writeFile, readdir, stat } from 'fs/promises';
import { join, relative } from 'path';
import { createHash } from 'crypto';

const LOCAL_LOCK_FILE = 'skills-lock.json';
const CURRENT_VERSION = 1;

/**
 * 本地（项目）lock 文件中的单条技能条目。
 *
 * 有意保持精简、无时间戳，以减少合并冲突。
 * 不同分支添加不同技能时 JSON 键不重叠，git 可自动合并。
 */
export interface LocalSkillLockEntry {
  /** 技能来源：npm 包名、owner/repo、本地路径等 */
  source: string;
  /** 安装时使用的分支或标签 ref */
  ref?: string;
  /** provider/来源类型（如 "github"、"node_modules"、"local"） */
  sourceType: string;
  /**
   * 源仓库内 SKILL.md 路径（如 "skills/pdf/SKILL.md"）。
   * update 时仅重装该技能所必需；否则 update 会拉取源仓库全部技能。
   * 对旧 lock 可选；node_modules/本地路径等非仓库来源可省略。
   */
  skillPath?: string;
  /**
   * 对技能文件夹内全部文件计算的 SHA-256。
   * 与全局 lock 使用 GitHub tree SHA 不同，本地 lock 基于磁盘内容。
   */
  computedHash: string;
}

/**
 * 本地（项目级）技能 lock 文件结构。
 * 设计为纳入版本控制。
 *
 * 写入时按技能名排序，输出确定、减少合并冲突。
 */
export interface LocalSkillLockFile {
  /** 模式版本，便于未来迁移 */
  version: number;
  /** 技能名 → lock 条目（按字母序写入） */
  skills: Record<string, LocalSkillLockEntry>;
}

/** 获取项目的本地技能 lock 文件路径。 */
export function getLocalLockPath(cwd?: string): string {
  return join(cwd || process.cwd(), LOCAL_LOCK_FILE);
}

/**
 * 读取本地技能 lock 文件。
 * 文件不存在或损坏（如含合并冲突标记）时返回空结构。
 */
export async function readLocalLock(cwd?: string): Promise<LocalSkillLockFile> {
  const lockPath = getLocalLockPath(cwd);

  try {
    const content = await readFile(lockPath, 'utf-8');
    const parsed = JSON.parse(content) as LocalSkillLockFile;

    if (typeof parsed.version !== 'number' || !parsed.skills) {
      return createEmptyLocalLock();
    }

    if (parsed.version < CURRENT_VERSION) {
      return createEmptyLocalLock();
    }

    return parsed;
  } catch {
    return createEmptyLocalLock();
  }
}

/**
 * 写入本地技能 lock 文件。
 * 技能按名称字母序排序，保证输出确定。
 */
export async function writeLocalLock(lock: LocalSkillLockFile, cwd?: string): Promise<void> {
  const lockPath = getLocalLockPath(cwd);

  // 按字母序排序，便于 diff / 合并
  const sortedSkills: Record<string, LocalSkillLockEntry> = {};
  for (const key of Object.keys(lock.skills).sort()) {
    sortedSkills[key] = lock.skills[key]!;
  }

  const sorted: LocalSkillLockFile = { version: lock.version, skills: sortedSkills };
  const content = JSON.stringify(sorted, null, 2) + '\n';
  await writeFile(lockPath, content, 'utf-8');
}

/**
 * 对技能目录内全部文件计算 SHA-256。
 * 递归读取、按相对路径排序后拼接内容再哈希。
 */
export async function computeSkillFolderHash(skillDir: string): Promise<string> {
  const files: Array<{ relativePath: string; content: Buffer }> = [];
  await collectFiles(skillDir, skillDir, files);

  // 按相对路径排序，保证哈希确定
  files.sort((a, b) => a.relativePath.localeCompare(b.relativePath));

  const hash = createHash('sha256');
  for (const file of files) {
    // 路径纳入哈希，以便检测重命名
    hash.update(file.relativePath);
    hash.update(file.content);
  }

  return hash.digest('hex');
}

async function collectFiles(
  baseDir: string,
  currentDir: string,
  results: Array<{ relativePath: string; content: Buffer }>
): Promise<void> {
  const entries = await readdir(currentDir, { withFileTypes: true });

  await Promise.all(
    entries.map(async (entry) => {
      const fullPath = join(currentDir, entry.name);

      if (entry.isDirectory()) {
        // 跳过技能目录内的 .git 与 node_modules
        if (entry.name === '.git' || entry.name === 'node_modules') return;
        await collectFiles(baseDir, fullPath, results);
      } else if (entry.isFile()) {
        const content = await readFile(fullPath);
        const relativePath = relative(baseDir, fullPath).split('\\').join('/');
        results.push({ relativePath, content });
      }
    })
  );
}

/** 在本地 lock 中新增或更新技能条目。 */
export async function addSkillToLocalLock(
  skillName: string,
  entry: LocalSkillLockEntry,
  cwd?: string
): Promise<void> {
  const lock = await readLocalLock(cwd);
  lock.skills[skillName] = entry;
  await writeLocalLock(lock, cwd);
}

/** 从本地 lock 移除技能。 */
export async function removeSkillFromLocalLock(skillName: string, cwd?: string): Promise<boolean> {
  const lock = await readLocalLock(cwd);

  if (!(skillName in lock.skills)) {
    return false;
  }

  delete lock.skills[skillName];
  await writeLocalLock(lock, cwd);
  return true;
}

function createEmptyLocalLock(): LocalSkillLockFile {
  return {
    version: CURRENT_VERSION,
    skills: {},
  };
}
