import { readdir, readFile, stat } from 'fs/promises';
import { join, basename, dirname, resolve, normalize, sep } from 'path';
import { parseFrontmatter } from './frontmatter.ts';
import { sanitizeMetadata } from './sanitize.ts';
import type { Skill } from './types.ts';
import { getPluginSkillPaths, getPluginGroupings } from './plugin-manifest.ts';

const SKIP_DIRS = ['node_modules', '.git', 'dist', 'build', '__pycache__'];

/**
 * 是否应安装内部（internal）技能。
 * 默认隐藏，除非设置 INSTALL_INTERNAL_SKILLS=1。
 */
export function shouldInstallInternalSkills(): boolean {
  const envValue = process.env.INSTALL_INTERNAL_SKILLS;
  return envValue === '1' || envValue === 'true';
}

async function hasSkillMd(dir: string): Promise<boolean> {
  try {
    const skillPath = join(dir, 'SKILL.md');
    const stats = await stat(skillPath);
    return stats.isFile();
  } catch {
    return false;
  }
}

export async function parseSkillMd(
  skillMdPath: string,
  options?: { includeInternal?: boolean }
): Promise<Skill | null> {
  try {
    const content = await readFile(skillMdPath, 'utf-8');
    const { data } = parseFrontmatter(content);

    if (!data.name || !data.description) {
      return null;
    }

    // 确保 name、description 为字符串（YAML 可能解析为数字、布尔等）
    if (typeof data.name !== 'string' || typeof data.description !== 'string') {
      return null;
    }

    // 除非满足以下之一，否则跳过 internal 技能：
    // 1. INSTALL_INTERNAL_SKILLS=1，或
    // 2. includeInternal 为 true（用户显式请求该技能）
    const isInternal = data.metadata?.internal === true;
    if (isInternal && !shouldInstallInternalSkills() && !options?.includeInternal) {
      return null;
    }

    return {
      name: sanitizeMetadata(data.name),
      description: sanitizeMetadata(data.description),
      path: dirname(skillMdPath),
      rawContent: content,
      metadata: data.metadata,
    };
  } catch {
    return null;
  }
}

async function findSkillDirs(dir: string, depth = 0, maxDepth = 5): Promise<string[]> {
  if (depth > maxDepth) return [];

  try {
    const [hasSkill, entries] = await Promise.all([
      hasSkillMd(dir),
      readdir(dir, { withFileTypes: true }).catch(() => []),
    ]);

    const currentDir = hasSkill ? [dir] : [];

    // 并行搜索子目录
    const subDirResults = await Promise.all(
      entries
        .filter((entry) => entry.isDirectory() && !SKIP_DIRS.includes(entry.name))
        .map((entry) => findSkillDirs(join(dir, entry.name), depth + 1, maxDepth))
    );

    return [...currentDir, ...subDirResults.flat()];
  } catch {
    return [];
  }
}

export interface DiscoverSkillsOptions {
  /** 包含 internal 技能（如用户按名称显式请求时） */
  includeInternal?: boolean;
  /** 即使根目录已有 SKILL.md 也搜索全部子目录 */
  fullDepth?: boolean;
}

/**
 * 校验解析后的 subpath 仍在 base 目录内。
 * 防止 subpath 含 ".." 逃出克隆仓库目录的路径遍历。
 */
export function isSubpathSafe(basePath: string, subpath: string): boolean {
  const normalizedBase = normalize(resolve(basePath));
  const normalizedTarget = normalize(resolve(join(basePath, subpath)));

  return normalizedTarget.startsWith(normalizedBase + sep) || normalizedTarget === normalizedBase;
}

export async function discoverSkills(
  basePath: string,
  subpath?: string,
  options?: DiscoverSkillsOptions
): Promise<Skill[]> {
  const skills: Skill[] = [];
  const seenNames = new Set<string>();

  // 校验 subpath 不逃出 basePath（防路径遍历）
  if (subpath && !isSubpathSafe(basePath, subpath)) {
    throw new Error(
      `Invalid subpath: "${subpath}" resolves outside the repository directory. Subpath must not contain ".." segments that escape the base path.`
    );
  }

  const searchPath = subpath ? join(basePath, subpath) : basePath;

  // 获取插件分组，将技能映射到父插件
  // 从 base 搜索路径查找插件定义
  const pluginGroupings = await getPluginGroupings(searchPath);

  // 若有插件名则写入 skill
  const enhanceSkill = (skill: Skill) => {
    const resolvedPath = resolve(skill.path);
    if (pluginGroupings.has(resolvedPath)) {
      skill.pluginName = pluginGroupings.get(resolvedPath);
    }
    return skill;
  };

  // 若直接指向某技能目录，加入该技能（除非设置了 fullDepth 则继续搜）
  if (await hasSkillMd(searchPath)) {
    let skill = await parseSkillMd(join(searchPath, 'SKILL.md'), options);
    if (skill) {
      skill = enhanceSkill(skill);
      skills.push(skill);
      seenNames.add(skill.name);
      // 未设置 fullDepth 时提前返回
      if (!options?.fullDepth) {
        return skills;
      }
    }
  }

  // 优先搜索常见技能目录
  const prioritySearchDirs = [
    searchPath,
    join(searchPath, 'skills'),
    join(searchPath, 'skills/.curated'),
    join(searchPath, 'skills/.experimental'),
    join(searchPath, 'skills/.system'),
    join(searchPath, '.agents/skills'),
    join(searchPath, '.claude/skills'),
    join(searchPath, '.cline/skills'),
    join(searchPath, '.codebuddy/skills'),
    join(searchPath, '.codex/skills'),
    join(searchPath, '.commandcode/skills'),
    join(searchPath, '.continue/skills'),

    join(searchPath, '.github/skills'),
    join(searchPath, '.goose/skills'),
    join(searchPath, '.iflow/skills'),
    join(searchPath, '.junie/skills'),
    join(searchPath, '.kilocode/skills'),
    join(searchPath, '.kiro/skills'),
    join(searchPath, '.mux/skills'),
    join(searchPath, '.neovate/skills'),
    join(searchPath, '.opencode/skills'),
    join(searchPath, '.openhands/skills'),
    join(searchPath, '.pi/skills'),
    join(searchPath, '.qoder/skills'),
    join(searchPath, '.roo/skills'),
    join(searchPath, '.trae/skills'),
    join(searchPath, '.windsurf/skills'),
    join(searchPath, '.zencoder/skills'),
  ];

  // 加入插件 manifest 声明的技能路径
  prioritySearchDirs.push(...(await getPluginSkillPaths(searchPath)));

  for (const dir of prioritySearchDirs) {
    try {
      const entries = await readdir(dir, { withFileTypes: true });

      for (const entry of entries) {
        if (entry.isDirectory()) {
          const skillDir = join(dir, entry.name);
          if (await hasSkillMd(skillDir)) {
            let skill = await parseSkillMd(join(skillDir, 'SKILL.md'), options);
            if (skill && !seenNames.has(skill.name)) {
              skill = enhanceSkill(skill);
              skills.push(skill);
              seenNames.add(skill.name);
            }
          }
        }
      }
    } catch {
      // 目录不存在
    }
  }

  // 未找到任何技能或设置了 fullDepth 时，回退到递归搜索
  if (skills.length === 0 || options?.fullDepth) {
    const allSkillDirs = await findSkillDirs(searchPath);

    for (const skillDir of allSkillDirs) {
      let skill = await parseSkillMd(join(skillDir, 'SKILL.md'), options);
      if (skill && !seenNames.has(skill.name)) {
        skill = enhanceSkill(skill);
        skills.push(skill);
        seenNames.add(skill.name);
      }
    }
  }

  return skills;
}

export function getSkillDisplayName(skill: Skill): string {
  return skill.name || basename(skill.path);
}

/**
 * 按用户输入过滤技能（大小写不敏感、直接匹配）。
 * 多词技能名须在命令行加引号。
 */
export function filterSkills(skills: Skill[], inputNames: string[]): Skill[] {
  const normalizedInputs = inputNames.map((n) => n.toLowerCase());

  return skills.filter((skill) => {
    const name = skill.name.toLowerCase();
    const displayName = getSkillDisplayName(skill).toLowerCase();

    return normalizedInputs.some((input) => input === name || input === displayName);
  });
}
