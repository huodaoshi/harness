import { isAbsolute, resolve } from 'path';
import type { ParsedSource } from './types.ts';

/**
 * 从已解析的 source 提取 owner/repo（GitLab 可为 group/subgroup/repo），
 * 用于 lockfile 与遥测。
 * 本地路径或无法解析时返回 null。
 * 支持具备 owner/repo URL 结构的任意 Git 宿主（含 GitLab 子组）。
 */
export function getOwnerRepo(parsed: ParsedSource): string | null {
  if (parsed.type === 'local') {
    return null;
  }

  // Git SSH URL（如 git@gitlab.com:owner/repo.git、git@github.com:owner/repo.git）
  const sshMatch = parsed.url.match(/^git@[^:]+:(.+)$/);
  if (sshMatch) {
    let path = sshMatch[1]!;
    path = path.replace(/\.git$/, '');

    // 至少包含 owner/repo（一个斜杠）
    if (path.includes('/')) {
      return path;
    }
    return null;
  }

  // HTTP(S) URL
  if (!parsed.url.startsWith('http://') && !parsed.url.startsWith('https://')) {
    return null;
  }

  try {
    const url = new URL(parsed.url);
    // 取 pathname，去掉 leading slash 与尾部 .git
    let path = url.pathname.slice(1);
    path = path.replace(/\.git$/, '');

    // 至少包含 owner/repo（一个斜杠）
    if (path.includes('/')) {
      return path;
    }
  } catch {
    // 无效 URL
  }

  return null;
}

/**
 * 从 owner/repo 字符串解析 owner 与 repo。
 * 格式无效时返回 null。
 */
export function parseOwnerRepo(ownerRepo: string): { owner: string; repo: string } | null {
  const match = ownerRepo.match(/^([^/]+)\/([^/]+)$/);
  if (match) {
    return { owner: match[1]!, repo: match[2]! };
  }
  return null;
}

/**
 * 检查 GitHub 仓库是否为私有。
 * 私有返回 true，公开返回 false，无法判断返回 null。
 * 仅支持 GitHub（不支持 GitLab）。
 */
export async function isRepoPrivate(owner: string, repo: string): Promise<boolean | null> {
  try {
    const res = await fetch(`https://api.github.com/repos/${owner}/${repo}`);

    // 仓库不存在或无权限时，为安全起见视为无法判断
    if (!res.ok) {
      return null; // 无法确定
    }

    const data = (await res.json()) as { private?: boolean };
    return data.private === true;
  } catch {
    // 出错时返回 null 表示无法确定
    return null;
  }
}

/**
 * 清理 subpath，防止路径遍历攻击。
 * 拒绝包含 ".." 段、可能逃出仓库根的 subpath。
 * 返回清理后的 subpath；不安全时抛出错误。
 */
export function sanitizeSubpath(subpath: string): string {
  // 统一为正向斜杠
  const normalized = subpath.replace(/\\/g, '/');

  // 检查各段是否含 ".."
  const segments = normalized.split('/');
  for (const segment of segments) {
    if (segment === '..') {
      throw new Error(
        `Unsafe subpath: "${subpath}" contains path traversal segments. ` +
          `Subpaths must not contain ".." components.`
      );
    }
  }

  return subpath;
}

/** 判断字符串是否为本地文件系统路径 */
function isLocalPath(input: string): boolean {
  return (
    isAbsolute(input) ||
    input.startsWith('./') ||
    input.startsWith('../') ||
    input === '.' ||
    input === '..' ||
    // Windows 绝对路径，如 C:\ 或 D:\
    /^[a-zA-Z]:[/\\]/.test(input)
  );
}

/**
 * 将 source 字符串解析为结构化格式。
 * 支持：本地路径、GitHub/GitLab URL、GitHub 简写、well-known URL、直接 git URL。
 */
// Source 别名：常见简写 → 规范 source
const SOURCE_ALIASES: Record<string, string> = {
  'coinbase/agentWallet': 'coinbase/agentic-wallet-skills',
};

interface FragmentRefResult {
  inputWithoutFragment: string;
  ref?: string;
  skillFilter?: string;
}

function decodeFragmentValue(value: string): string {
  try {
    return decodeURIComponent(value);
  } catch {
    return value;
  }
}

function looksLikeGitSource(input: string): boolean {
  if (input.startsWith('github:') || input.startsWith('gitlab:') || input.startsWith('git@')) {
    return true;
  }

  if (input.startsWith('http://') || input.startsWith('https://')) {
    try {
      const parsed = new URL(input);
      const pathname = parsed.pathname;

      // 仅对 repo/tree URL 将 GitHub fragment 视为 ref
      if (parsed.hostname === 'github.com') {
        return /^\/[^/]+\/[^/]+(?:\.git)?(?:\/tree\/[^/]+(?:\/.*)?)?\/?$/.test(pathname);
      }

      // 仅对 repo/tree URL 将 gitlab.com fragment 视为 ref
      if (parsed.hostname === 'gitlab.com') {
        return /^\/.+?\/[^/]+(?:\.git)?(?:\/-\/tree\/[^/]+(?:\/.*)?)?\/?$/.test(pathname);
      }
    } catch {
      // 落入下方通用检查
    }
  }

  if (/^https?:\/\/.+\.git(?:$|[/?])/i.test(input)) {
    return true;
  }

  return (
    !input.includes(':') &&
    !input.startsWith('.') &&
    !input.startsWith('/') &&
    /^([^/]+)\/([^/]+)(?:\/(.+)|@(.+))?$/.test(input)
  );
}

function parseFragmentRef(input: string): FragmentRefResult {
  const hashIndex = input.indexOf('#');
  if (hashIndex < 0) {
    return { inputWithoutFragment: input };
  }

  const inputWithoutFragment = input.slice(0, hashIndex);
  const fragment = input.slice(hashIndex + 1);

  // 仅对 git 类 source 将 URL fragment 视为 git ref，
  // 避免改变通用 well-known URL 的行为
  if (!fragment || !looksLikeGitSource(inputWithoutFragment)) {
    return { inputWithoutFragment: input };
  }

  const atIndex = fragment.indexOf('@');
  if (atIndex === -1) {
    return {
      inputWithoutFragment,
      ref: decodeFragmentValue(fragment),
    };
  }

  const ref = fragment.slice(0, atIndex);
  const skillFilter = fragment.slice(atIndex + 1);
  return {
    inputWithoutFragment,
    ref: ref ? decodeFragmentValue(ref) : undefined,
    skillFilter: skillFilter ? decodeFragmentValue(skillFilter) : undefined,
  };
}

function appendFragmentRef(input: string, ref?: string, skillFilter?: string): string {
  if (!ref) {
    return input;
  }
  return `${input}#${ref}${skillFilter ? `@${skillFilter}` : ''}`;
}

export function parseSource(input: string): ParsedSource {
  // 本地路径：绝对、相对或当前目录
  if (isLocalPath(input)) {
    const resolvedPath = resolve(input);
    // 即使路径不存在也返回 local — 在主流程中再校验
    return {
      type: 'local',
      url: resolvedPath, // 为一致性存入 url
      localPath: resolvedPath,
    };
  }

  const {
    inputWithoutFragment,
    ref: fragmentRef,
    skillFilter: fragmentSkillFilter,
  } = parseFragmentRef(input);
  input = inputWithoutFragment;

  // 解析前先应用 source 别名
  const alias = SOURCE_ALIASES[input];
  if (alias) {
    input = alias;
  }

  // 前缀简写：github:owner/repo → owner/repo（由现有简写逻辑处理）
  // 亦支持 github:owner/repo/subpath、github:owner/repo@skill
  const githubPrefixMatch = input.match(/^github:(.+)$/);
  if (githubPrefixMatch) {
    return parseSource(appendFragmentRef(githubPrefixMatch[1]!, fragmentRef, fragmentSkillFilter));
  }

  // 前缀简写：gitlab:owner/repo → https://gitlab.com/owner/repo
  const gitlabPrefixMatch = input.match(/^gitlab:(.+)$/);
  if (gitlabPrefixMatch) {
    return parseSource(
      appendFragmentRef(
        `https://gitlab.com/${gitlabPrefixMatch[1]!}`,
        fragmentRef,
        fragmentSkillFilter
      )
    );
  }

  // 带路径的 GitHub URL：https://github.com/owner/repo/tree/branch/path/to/skill
  const githubTreeWithPathMatch = input.match(/github\.com\/([^/]+)\/([^/]+)\/tree\/([^/]+)\/(.+)/);
  if (githubTreeWithPathMatch) {
    const [, owner, repo, ref, subpath] = githubTreeWithPathMatch;
    return {
      type: 'github',
      url: `https://github.com/${owner}/${repo}.git`,
      ref: ref || fragmentRef,
      subpath: subpath ? sanitizeSubpath(subpath) : subpath,
    };
  }

  // 仅分支的 GitHub URL：https://github.com/owner/repo/tree/branch
  const githubTreeMatch = input.match(/github\.com\/([^/]+)\/([^/]+)\/tree\/([^/]+)$/);
  if (githubTreeMatch) {
    const [, owner, repo, ref] = githubTreeMatch;
    return {
      type: 'github',
      url: `https://github.com/${owner}/${repo}.git`,
      ref: ref || fragmentRef,
    };
  }

  // GitHub URL：https://github.com/owner/repo
  const githubRepoMatch = input.match(/github\.com\/([^/]+)\/([^/]+)/);
  if (githubRepoMatch) {
    const [, owner, repo] = githubRepoMatch;
    const cleanRepo = repo!.replace(/\.git$/, '');
    return {
      type: 'github',
      url: `https://github.com/${owner}/${cleanRepo}.git`,
      ...(fragmentRef ? { ref: fragmentRef } : {}),
    };
  }

  // 带路径的 GitLab URL（任意实例）：https://gitlab.com/owner/repo/-/tree/branch/path
  // 关键标识为 GitLab 特有的 "/-/tree/" 路径模式
  // 用非贪婪匹配支持子组仓库路径
  const gitlabTreeWithPathMatch = input.match(
    /^(https?):\/\/([^/]+)\/(.+?)\/-\/tree\/([^/]+)\/(.+)/
  );
  if (gitlabTreeWithPathMatch) {
    const [, protocol, hostname, repoPath, ref, subpath] = gitlabTreeWithPathMatch;
    if (hostname !== 'github.com' && repoPath) {
      return {
        type: 'gitlab',
        url: `${protocol}://${hostname}/${repoPath.replace(/\.git$/, '')}.git`,
        ref: ref || fragmentRef,
        subpath: subpath ? sanitizeSubpath(subpath) : subpath,
      };
    }
  }

  // 仅分支的 GitLab URL（任意实例）：https://gitlab.com/owner/repo/-/tree/branch
  const gitlabTreeMatch = input.match(/^(https?):\/\/([^/]+)\/(.+?)\/-\/tree\/([^/]+)$/);
  if (gitlabTreeMatch) {
    const [, protocol, hostname, repoPath, ref] = gitlabTreeMatch;
    if (hostname !== 'github.com' && repoPath) {
      return {
        type: 'gitlab',
        url: `${protocol}://${hostname}/${repoPath.replace(/\.git$/, '')}.git`,
        ref: ref || fragmentRef,
      };
    }
  }

  // GitLab.com URL：https://gitlab.com/owner/repo 或 group/subgroup/repo
  // 仅官方 gitlab.com 域，便于用户使用
  // 支持嵌套子组（如 gitlab.com/group/subgroup1/subgroup2/repo）
  const gitlabRepoMatch = input.match(/gitlab\.com\/(.+?)(?:\.git)?\/?$/);
  if (gitlabRepoMatch) {
    const repoPath = gitlabRepoMatch[1]!;
    // 至少 owner/repo（一个斜杠）
    if (repoPath.includes('/')) {
      return {
        type: 'gitlab',
        url: `https://gitlab.com/${repoPath}.git`,
        ...(fragmentRef ? { ref: fragmentRef } : {}),
      };
    }
  }

  // GitHub 简写：owner/repo、owner/repo/path/to/skill、owner/repo@skill-name
  // 排除以 . 或 / 开头的路径，避免匹配本地路径
  // 先检查 @skill：owner/repo@skill-name
  const atSkillMatch = input.match(/^([^/]+)\/([^/@]+)@(.+)$/);
  if (atSkillMatch && !input.includes(':') && !input.startsWith('.') && !input.startsWith('/')) {
    const [, owner, repo, skillFilter] = atSkillMatch;
    return {
      type: 'github',
      url: `https://github.com/${owner}/${repo}.git`,
      ...(fragmentRef ? { ref: fragmentRef } : {}),
      skillFilter: fragmentSkillFilter || skillFilter,
    };
  }

  const shorthandMatch = input.match(/^([^/]+)\/([^/]+)(?:\/(.+?))?\/?$/);
  if (shorthandMatch && !input.includes(':') && !input.startsWith('.') && !input.startsWith('/')) {
    const [, owner, repo, subpath] = shorthandMatch;
    return {
      type: 'github',
      url: `https://github.com/${owner}/${repo}.git`,
      ...(fragmentRef ? { ref: fragmentRef } : {}),
      subpath: subpath ? sanitizeSubpath(subpath) : subpath,
      ...(fragmentSkillFilter ? { skillFilter: fragmentSkillFilter } : {}),
    };
  }

  // Well-known 技能：非 GitHub/GitLab 的任意 HTTP(S) URL
  // 最终回退：先查 /.well-known/agent-skills/index.json，
  // 再回退到 /.well-known/skills/index.json
  if (isWellKnownUrl(input)) {
    return {
      type: 'well-known',
      url: input,
    };
  }

  // 回退：视为直接 git URL
  return {
    type: 'git',
    url: input,
    ...(fragmentRef ? { ref: fragmentRef } : {}),
  };
}

/**
 * 判断 URL 是否可能为 well-known 技能端点。
 * 须为 HTTP(S)，且非已知 git 宿主（GitHub、GitLab）。
 * 亦排除形似 git 仓库的 URL（.git 后缀）。
 */
function isWellKnownUrl(input: string): boolean {
  if (!input.startsWith('http://') && !input.startsWith('https://')) {
    return false;
  }

  try {
    const parsed = new URL(input);

    // 排除已有专门处理的 git 宿主
    const excludedHosts = ['github.com', 'gitlab.com', 'raw.githubusercontent.com'];
    if (excludedHosts.includes(parsed.hostname)) {
      return false;
    }

    // 不匹配形似 git 仓库的 URL（应由 git 类型处理）
    if (input.endsWith('.git')) {
      return false;
    }

    return true;
  } catch {
    return false;
  }
}
