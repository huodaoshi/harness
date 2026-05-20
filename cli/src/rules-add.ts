import * as p from '@clack/prompts';
import pc from 'picocolors';
import { existsSync } from 'fs';
import { cp, rm, mkdir } from 'fs/promises';
import { join, resolve, dirname } from 'path';
import { resolveCliProjectRoot } from './project-root.ts';
import { parseSource } from './source-parser.ts';
import { cloneRepo, cleanupTempDir } from './git.ts';
import { agents } from './agents.ts';
import {
  writeRulesLock,
  readRulesLock,
  type RulesLockAgent,
  type RulesLockInstall,
} from './rules-lock.ts';

export const RULES_AGENT_TYPES: readonly RulesLockAgent[] = ['cursor', 'claude-code'] as const;

function isRulesAgentType(a: string): a is RulesLockAgent {
  return a === 'cursor' || a === 'claude-code';
}

export function getRulesPackSubdir(agent: RulesLockAgent): 'rules/cursor' | 'rules/claude' {
  return agent === 'cursor' ? 'rules/cursor' : 'rules/claude';
}

export interface RulesAddOptions {
  agent?: string[];
  global?: boolean;
  to?: string;
  yes?: boolean;
  copy?: boolean;
  /** 从 lock 恢复时不再次写锁 */
  skipWriteLock?: boolean;
  /**
   * 安装目标「项目根」（写入 `.cursor/rules`、`rules-lock.json`）。
   * 与 `process.cwd()` 不同：用 `pnpm --dir path/to/cli` 时 cwd 常为 cli 包目录，应传递此字段或使用 INIT_CWD。
   */
  cwd?: string;
}

/** 规则/锁文件与 `resolveCliProjectRoot` 使用同一套项目根解析。 */
export function resolveRulesProjectRoot(explicitCwd?: string): string {
  return resolveCliProjectRoot(explicitCwd);
}

export function parseRulesAddOptions(args: string[]): { source: string[]; options: RulesAddOptions } {
  const options: RulesAddOptions = {};
  const source: string[] = [];

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];

    if (arg === '-g' || arg === '--global') {
      options.global = true;
    } else if (arg === '-y' || arg === '--yes') {
      options.yes = true;
    } else if (arg === '--copy') {
      options.copy = true;
    } else if (arg === '--to') {
      i++;
      const next = args[i];
      if (next) options.to = next;
    } else if (arg === '--cwd' || arg === '-C') {
      i++;
      const next = args[i];
      if (next) options.cwd = next;
    } else if (arg === '-a' || arg === '--agent') {
      options.agent = options.agent || [];
      i++;
      while (i < args.length && args[i] && !args[i]!.startsWith('-')) {
        options.agent!.push(args[i]!);
        i++;
      }
      i--;
    } else if (arg && !arg.startsWith('-')) {
      source.push(arg);
    }
  }

  return { source, options };
}

function resolveRulesDestDir(
  agent: RulesLockAgent,
  options: { global?: boolean; cwd?: string; to?: string }
): string {
  if (options.to) {
    return resolve(options.to);
  }

  const cfg = agents[agent];
  const cwd = options.cwd || process.cwd();

  if (options.global) {
    if (!cfg.globalRulesDir) {
      throw new Error(`${cfg.displayName} does not support global rules installation`);
    }
    return cfg.globalRulesDir;
  }

  if (!cfg.rulesDir) {
    throw new Error(`${cfg.displayName} does not support project rules installation`);
  }

  return join(cwd, cfg.rulesDir);
}

function repoRootFromParsed(
  basePath: string,
  parsed: { subpath?: string }
): string {
  if (parsed.subpath) {
    return join(basePath, parsed.subpath);
  }
  return basePath;
}

async function copyRulesTree(
  packPath: string,
  destPath: string
): Promise<void> {
  try {
    await rm(destPath, { recursive: true, force: true });
  } catch {
    // ignore
  }

  await mkdir(dirname(destPath), { recursive: true });
  await cp(packPath, destPath, {
    recursive: true,
    force: true,
    errorOnExist: false,
  });
}

/**
 * 安装规则包：从 `rules/cursor` 或 `rules/claude` 复制到目标目录。
 */
export async function runRulesAdd(
  sourceArgs: string[],
  passedOptions?: RulesAddOptions
): Promise<void> {
  const { source, options: parsedOpts } = parseRulesAddOptions(sourceArgs);
  const options: RulesAddOptions = { ...parsedOpts, ...passedOptions };
  const projectRoot = resolveRulesProjectRoot(options.cwd);

  if (source.length === 0) {
    p.log.error(pc.red('Missing source. Example: skills rules add vercel-labs/foo -a cursor'));
    process.exit(1);
  }

  const agentsList = options.agent?.filter(Boolean) ?? [];
  if (agentsList.length !== 1) {
    p.log.error(
      pc.red(
        'Exactly one agent is required: -a cursor | -a claude-code (MVP does not support multiple agents per invocation)'
      )
    );
    process.exit(1);
  }

  const agentArg = agentsList[0]!;
  if (!isRulesAgentType(agentArg)) {
    p.log.error(
      pc.red(
        `Unsupported agent "${agentArg}". For rules, only: ${RULES_AGENT_TYPES.join(', ')}`
      )
    );
    process.exit(1);
  }

  const agent: RulesLockAgent = agentArg;

  const input = source[0]!;
  const parsed = parseSource(input);
  let tempDir: string | null = null;

  if (!options.skipWriteLock) {
    p.intro(`${pc.bgCyan(pc.black(pc.bold(' rules ')))} Install`);
  }

  try {
    let repoRoot: string;

    if (parsed.type === 'local') {
      const lp = parsed.localPath!;
      if (!existsSync(lp)) {
        p.outro(pc.red(`Path does not exist: ${lp}`));
        process.exit(1);
      }
      repoRoot = repoRootFromParsed(lp, parsed);
    } else {
      const spinner = p.spinner();
      spinner.start('Cloning repository...');
      try {
        tempDir = await cloneRepo(parsed.url, parsed.ref);
        spinner.stop('Repository cloned');
      } catch (e) {
        spinner.stop('Clone failed');
        throw e;
      }
      repoRoot = repoRootFromParsed(tempDir, parsed);
    }

    const packRel = getRulesPackSubdir(agent);
    const packPath = join(repoRoot, packRel);

    if (!existsSync(packPath)) {
      p.outro(
        pc.red(
          `Rules pack not found at ${packRel}/ (resolved: ${packPath}). ` +
            `Add that folder to the repository or use a different -a.`
        )
      );
      process.exit(1);
    }

    const dest = resolveRulesDestDir(agent, {
      global: options.global,
      cwd: projectRoot,
      to: options.to,
    });

    if (projectRoot !== process.cwd()) {
      p.log.message(pc.dim(`Target project: ${projectRoot}`));
    }

    const spinner = p.spinner();
    spinner.start(`Copying rules to ${pc.dim(dest)}...`);
    await copyRulesTree(packPath, dest);
    spinner.stop('Rules installed');

    if (!options.skipWriteLock) {
      const installRecord: RulesLockInstall = {
        agent,
        source: input,
        ref: parsed.ref,
        sourceType: parsed.type,
        copy: options.copy ?? true,
      };
      const existing = await readRulesLock(projectRoot);
      await writeRulesLock({ version: existing.version, install: installRecord }, projectRoot);
      p.log.message(pc.dim(`Wrote ${pc.cyan('rules-lock.json')}`));
    }

    p.outro(pc.green('Done'));
  } finally {
    if (tempDir) {
      await cleanupTempDir(tempDir).catch(() => {});
    }
  }
}

/**
 * 根据项目根 `rules-lock.json` 恢复规则（仅一条 install 记录）。
 */
export async function runRulesInstallFromLock(args: string[]): Promise<void> {
  const { options: parseOpts } = parseRulesAddOptions(args);
  const projectRoot = resolveRulesProjectRoot(parseOpts.cwd);
  const lock = await readRulesLock(projectRoot);

  if (!lock.install) {
    p.log.warn('No rules install entry in rules-lock.json');
    p.log.info(`Add rules with ${pc.cyan('npx skills rules add <source> -a cursor|-a claude-code')}`);
    return;
  }

  const { agent, source } = lock.install;
  if (!isRulesAgentType(agent)) {
    p.log.error(pc.red(`Invalid lock agent: ${agent}`));
    process.exit(1);
  }

  p.intro(`${pc.bgCyan(pc.black(pc.bold(' rules ')))} Restore from lock`);

  await runRulesAdd([source, '-a', agent, '-y'], {
    skipWriteLock: true,
    cwd: projectRoot,
  });
}
