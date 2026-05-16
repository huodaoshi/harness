import simpleGit from 'simple-git';
import { join, normalize, resolve, sep } from 'path';
import { mkdtemp, rm } from 'fs/promises';
import { tmpdir } from 'os';

const DEFAULT_CLONE_TIMEOUT_MS = 300_000; // 5 分钟
const CLONE_TIMEOUT_MS = (() => {
  const raw = process.env.SKILLS_CLONE_TIMEOUT_MS;
  if (!raw) return DEFAULT_CLONE_TIMEOUT_MS;
  const parsed = Number.parseInt(raw, 10);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : DEFAULT_CLONE_TIMEOUT_MS;
})();

/** 克隆时安全传给 git 的环境变量（避免 simple-git 拦截 EDITOR 等）。 */
const GIT_CLONE_ENV_KEYS = [
  'PATH',
  'HOME',
  'USERPROFILE',
  'HOMEDRIVE',
  'HOMEPATH',
  'SYSTEMROOT',
  'WINDIR',
  'TEMP',
  'TMP',
  'APPDATA',
  'LOCALAPPDATA',
  'PROGRAMFILES',
  'HTTP_PROXY',
  'HTTPS_PROXY',
  'ALL_PROXY',
  'NO_PROXY',
  'http_proxy',
  'https_proxy',
  'all_proxy',
  'no_proxy',
  'GIT_SSH_COMMAND',
  'SSH_AUTH_SOCK',
  // 不传 GIT_ASKPASS：Cursor/IDE 会设置它，simple-git 3.36+ 需 allowUnsafeAskPass；
  // 公开仓库克隆用 GIT_TERMINAL_PROMPT=0 即可，无需交互式凭据。
] as const;

function gitCloneEnv(): NodeJS.ProcessEnv {
  const env: NodeJS.ProcessEnv = {
    GIT_TERMINAL_PROMPT: '0',
    GIT_LFS_SKIP_SMUDGE: '1',
  };
  for (const key of GIT_CLONE_ENV_KEYS) {
    const value = process.env[key];
    if (value !== undefined) env[key] = value;
  }
  return env;
}

export class GitCloneError extends Error {
  readonly url: string;
  readonly isTimeout: boolean;
  readonly isAuthError: boolean;

  constructor(message: string, url: string, isTimeout = false, isAuthError = false) {
    super(message);
    this.name = 'GitCloneError';
    this.url = url;
    this.isTimeout = isTimeout;
    this.isAuthError = isAuthError;
  }
}

export async function cloneRepo(url: string, ref?: string): Promise<string> {
  const tempDir = await mkdtemp(join(tmpdir(), 'skills-'));
  // `env` 不是构造函数选项 — 请在实例上使用 .env()（见 SimpleGit.env）。
  const git = simpleGit({
    timeout: { block: CLONE_TIMEOUT_MS },
    // simple-git 3.36+ 默认禁止 -c filter.*，需显式开启。
    // 下列值由本 CLI 设置（空值 = 禁用 LFS 过滤器），非来自远程仓库
    // — 见 @simple-git/argv-parser 中的 allowUnsafeFilter。
    unsafe: {
      allowUnsafeFilter: true,
    },
    // 未安装 git-lfs 时，GIT_LFS_SKIP_SMUDGE 无效 —
    // git 在 .gitattributes 中看到 `filter=lfs` 会尝试运行
    // `git-lfs filter-process`，checkout 失败，例如：
    //   git-lfs filter-process: git-lfs: command not found
    //   fatal: the remote end hung up unexpectedly
    //   warning: Clone succeeded, but checkout failed.
    // 在命令级覆盖 filter.lfs.* 可完全禁用过滤器，使 checkout 成功，
    // 与是否安装 git-lfs 无关。LFS 文件会保留为约 130 字节的指针文件，
    // 技能安装器不会读取它们（技能为纯文本 HTML/MD/JSON，不走 LFS）。
    //
    // 下游反馈：heygen-com/hyperframes#407。
    config: [
      'filter.lfs.required=false',
      'filter.lfs.smudge=',
      'filter.lfs.clean=',
      'filter.lfs.process=',
    ],
  }).env(gitCloneEnv());
  const cloneOptions = ref ? ['--depth', '1', '--branch', ref] : ['--depth', '1'];

  try {
    await git.clone(url, tempDir, cloneOptions);
    return tempDir;
  } catch (error) {
    // 失败时清理临时目录
    await rm(tempDir, { recursive: true, force: true }).catch(() => {});

    const errorMessage = error instanceof Error ? error.message : String(error);
    const isTimeout = errorMessage.includes('block timeout') || errorMessage.includes('timed out');
    const isAuthError =
      errorMessage.includes('Authentication failed') ||
      errorMessage.includes('could not read Username') ||
      errorMessage.includes('Permission denied') ||
      errorMessage.includes('Repository not found');

    if (isTimeout) {
      const seconds = Math.round(CLONE_TIMEOUT_MS / 1000);
      throw new GitCloneError(
        `Clone timed out after ${seconds}s. Common causes:\n` +
          `  - Large repository: raise the timeout with SKILLS_CLONE_TIMEOUT_MS=600000 (10m)\n` +
          `  - Slow network: retry, or clone manually and pass the local path to 'skills add'\n` +
          `  - Private repo without credentials: ensure auth is configured\n` +
          `      - For SSH: ssh-add -l (to check loaded keys)\n` +
          `      - For HTTPS: gh auth status (if using GitHub CLI)`,
        url,
        true,
        false
      );
    }

    if (isAuthError) {
      throw new GitCloneError(
        `Authentication failed for ${url}.\n` +
          `  - For private repos, ensure you have access\n` +
          `  - For SSH: Check your keys with 'ssh -T git@github.com'\n` +
          `  - For HTTPS: Run 'gh auth login' or configure git credentials`,
        url,
        false,
        true
      );
    }

    throw new GitCloneError(`Failed to clone ${url}: ${errorMessage}`, url, false, false);
  }
}

export async function cleanupTempDir(dir: string): Promise<void> {
  // 校验目录在 tmpdir 内，防止删除任意路径
  const normalizedDir = normalize(resolve(dir));
  const normalizedTmpDir = normalize(resolve(tmpdir()));

  if (!normalizedDir.startsWith(normalizedTmpDir + sep) && normalizedDir !== normalizedTmpDir) {
    throw new Error('Attempted to clean up directory outside of temp directory');
  }

  await rm(dir, { recursive: true, force: true });
}
