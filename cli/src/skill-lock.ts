import { readFile, writeFile, mkdir } from 'fs/promises';
import { join, dirname } from 'path';
import { homedir } from 'os';
import { createHash } from 'crypto';
import { execSync } from 'child_process';

const AGENTS_DIR = '.agents';
const LOCK_FILE = '.skill-lock.json';
const CURRENT_VERSION = 3; // 自 v2 升至 v3：支持文件夹哈希（GitHub tree SHA）

/** 全局 lock 文件中的单条已安装技能。 */
export interface SkillLockEntry {
  /** 规范化的 source 标识（如 "owner/repo"、"mintlify/bun.com"） */
  source: string;
  /** provider/来源类型（如 "github"、"mintlify"、"huggingface"、"local"） */
  sourceType: string;
  /** 安装时使用的原始 URL（用于拉取更新） */
  sourceUrl: string;
  /** 安装使用的分支或标签 ref（用于按 ref 更新） */
  ref?: string;
  /** 源仓库内的子路径（若有） */
  skillPath?: string;
  /**
   * 整个技能文件夹的 GitHub tree SHA。
   * 文件夹内任意文件变更都会改变此哈希。
   * 由遥测服务通过 GitHub Trees API 获取。
   */
  skillFolderHash: string;
  /** 首次安装的 ISO 时间戳 */
  installedAt: string;
  /** 最后更新的 ISO 时间戳 */
  updatedAt: string;
  /** 所属插件名称（若有） */
  pluginName?: string;
}

/** 记录已关闭的提示，避免再次显示。 */
export interface DismissedPrompts {
  /** 已关闭 find-skills 安装提示 */
  findSkillsPrompt?: boolean;
}

/** 技能 lock 文件结构。 */
export interface SkillLockFile {
  /** 模式版本，便于未来迁移 */
  version: number;
  /** 技能名 → lock 条目 */
  skills: Record<string, SkillLockEntry>;
  /** 已关闭的提示 */
  dismissed?: DismissedPrompts;
  /** 上次安装时选择的 agent */
  lastSelectedAgents?: string[];
}

/**
 * 获取全局技能 lock 文件路径。
 * 若设置 $XDG_STATE_HOME，则用 $XDG_STATE_HOME/skills/.skill-lock.json；
 * 否则回退到 ~/.agents/.skill-lock.json。
 */
export function getSkillLockPath(): string {
  const xdgStateHome = process.env.XDG_STATE_HOME;
  if (xdgStateHome) {
    return join(xdgStateHome, 'skills', LOCK_FILE);
  }
  return join(homedir(), AGENTS_DIR, LOCK_FILE);
}

/**
 * 读取技能 lock 文件。
 * 不存在时返回空结构；旧格式（version < CURRENT_VERSION）会清空 lock。
 */
export async function readSkillLock(): Promise<SkillLockFile> {
  const lockPath = getSkillLockPath();

  try {
    const content = await readFile(lockPath, 'utf-8');
    const parsed = JSON.parse(content) as SkillLockFile;

    // 校验版本 — 旧格式则清空
    if (typeof parsed.version !== 'number' || !parsed.skills) {
      return createEmptyLockFile();
    }

    // 旧版本不兼容，重新开始（v3 增加 skillFolderHash，需全新安装以填充）
    if (parsed.version < CURRENT_VERSION) {
      return createEmptyLockFile();
    }

    return parsed;
  } catch {
    // 文件不存在或无效
    return createEmptyLockFile();
  }
}

/**
 * 写入技能 lock 文件。
 * 若目录不存在则创建。
 */
export async function writeSkillLock(lock: SkillLockFile): Promise<void> {
  const lockPath = getSkillLockPath();

  // 确保目录存在
  await mkdir(dirname(lockPath), { recursive: true });

  // 格式化输出，便于人工阅读
  const content = JSON.stringify(lock, null, 2);
  await writeFile(lockPath, content, 'utf-8');
}

/** 计算内容的 SHA-256 哈希。 */
export function computeContentHash(content: string): string {
  return createHash('sha256').update(content, 'utf-8').digest('hex');
}

let _ghWarningShown = false;

/** 仅测试用：重置一次性警告标记。 */
export function resetGhAuthWarning(): void {
  _ghWarningShown = false;
}

/**
 * 从用户环境获取 GitHub token。
 * 依次尝试：
 * 1. 环境变量 GITHUB_TOKEN（静默）
 * 2. 环境变量 GH_TOKEN（静默）
 * 3. 若已安装 gh，则 `gh auth token`；调用前会向 stderr 打印一次性警告，
 *    因部分企业终端安全工具会将该子进程视为凭据提取。
 *    调用方应惰性调用（如未认证请求触发限流后再取 token），
 *    实践中很少走到此回退。
 *
 * @returns token 字符串，不可用则 null
 */
export function getGitHubToken(): string | null {
  // 先查环境变量（用户已显式配置）
  if (process.env.GITHUB_TOKEN) {
    return process.env.GITHUB_TOKEN;
  }
  if (process.env.GH_TOKEN) {
    return process.env.GH_TOKEN;
  }

  // 最后回退：启动 gh CLI；每进程仅警告一次
  if (!_ghWarningShown) {
    process.stderr.write(
      'warn: GitHub API rate limit reached; reading a token via `gh auth token`.\n' +
        '      Set GITHUB_TOKEN in your environment to skip this fallback.\n'
    );
    _ghWarningShown = true;
  }
  try {
    const token = execSync('gh auth token', {
      encoding: 'utf-8',
      stdio: ['pipe', 'pipe', 'pipe'],
    }).trim();
    if (token) {
      return token;
    }
  } catch {
    // gh 未安装或未登录
  }

  return null;
}

/**
 * 用 GitHub Trees API 获取技能文件夹的 tree SHA（文件夹哈希）。
 * 一次 API 调用拉取整棵仓库树，再提取指定技能文件夹的 SHA。
 *
 * @param ownerRepo - GitHub owner/repo（如 "vercel-labs/agent-skills"）
 * @param skillPath - 技能文件夹或 SKILL.md 路径（如 "skills/react-best-practices/SKILL.md"）
 * @param getToken - 可选惰性 token 解析器；未认证限流时才调用
 * @param ref - 可选分支/标签；默认尝试 main 再 master
 * @returns 技能文件夹的 tree SHA，未找到则 null
 */
export async function fetchSkillFolderHash(
  ownerRepo: string,
  skillPath: string,
  getToken?: (() => string | null) | null,
  ref?: string
): Promise<string | null> {
  const { fetchRepoTree, getSkillFolderHashFromTree } = await import('./blob.ts');
  const tree = await fetchRepoTree(ownerRepo, ref, getToken ?? undefined);
  if (!tree) return null;
  return getSkillFolderHashFromTree(tree, skillPath);
}

/** 在 lock 中新增或更新技能条目。 */
export async function addSkillToLock(
  skillName: string,
  entry: Omit<SkillLockEntry, 'installedAt' | 'updatedAt'>
): Promise<void> {
  const lock = await readSkillLock();
  const now = new Date().toISOString();

  const existingEntry = lock.skills[skillName];

  lock.skills[skillName] = {
    ...entry,
    installedAt: existingEntry?.installedAt ?? now,
    updatedAt: now,
  };

  await writeSkillLock(lock);
}

/** 从 lock 移除技能。 */
export async function removeSkillFromLock(skillName: string): Promise<boolean> {
  const lock = await readSkillLock();

  if (!(skillName in lock.skills)) {
    return false;
  }

  delete lock.skills[skillName];
  await writeSkillLock(lock);
  return true;
}

/** 从 lock 读取单条技能。 */
export async function getSkillFromLock(skillName: string): Promise<SkillLockEntry | null> {
  const lock = await readSkillLock();
  return lock.skills[skillName] ?? null;
}

/** 读取 lock 中全部技能。 */
export async function getAllLockedSkills(): Promise<Record<string, SkillLockEntry>> {
  const lock = await readSkillLock();
  return lock.skills;
}

/** 按 source 分组技能，用于批量 update。 */
export async function getSkillsBySource(): Promise<
  Map<string, { skills: string[]; entry: SkillLockEntry }>
> {
  const lock = await readSkillLock();
  const bySource = new Map<string, { skills: string[]; entry: SkillLockEntry }>();

  for (const [skillName, entry] of Object.entries(lock.skills)) {
    const existing = bySource.get(entry.source);
    if (existing) {
      existing.skills.push(skillName);
    } else {
      bySource.set(entry.source, { skills: [skillName], entry });
    }
  }

  return bySource;
}

/** 创建空 lock 结构。 */
function createEmptyLockFile(): SkillLockFile {
  return {
    version: CURRENT_VERSION,
    skills: {},
    dismissed: {},
  };
}

/** 检查某提示是否已关闭。 */
export async function isPromptDismissed(promptKey: keyof DismissedPrompts): Promise<boolean> {
  const lock = await readSkillLock();
  return lock.dismissed?.[promptKey] === true;
}

/** 将某提示标记为已关闭。 */
export async function dismissPrompt(promptKey: keyof DismissedPrompts): Promise<void> {
  const lock = await readSkillLock();
  if (!lock.dismissed) {
    lock.dismissed = {};
  }
  lock.dismissed[promptKey] = true;
  await writeSkillLock(lock);
}

/** 获取上次选择的 agent 列表。 */
export async function getLastSelectedAgents(): Promise<string[] | undefined> {
  const lock = await readSkillLock();
  return lock.lastSelectedAgents;
}

/** 将所选 agent 保存到 lock 文件。 */
export async function saveSelectedAgents(agents: string[]): Promise<void> {
  const lock = await readSkillLock();
  lock.lastSelectedAgents = agents;
  await writeSkillLock(lock);
}
