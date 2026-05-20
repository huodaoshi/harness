import * as p from '@clack/prompts';
import pc from 'picocolors';
import { readLocalLock } from './local-lock.ts';
import { runAdd, parseAddOptions } from './add.ts';
import { runSync, parseSyncOptions } from './sync.ts';
import { getUniversalAgents } from './agents.ts';
import { resolveCliProjectRoot } from './project-root.ts';

/**
 * 从本地 skills-lock.json 安装全部技能。
 * 按 source 分组后对每组调用 `runAdd`。
 *
 * 仅安装到 `.agents/skills/`（通用 agent）——项目级规范位置；
 * 不安装到各 agent 专属目录。
 *
 * node_modules 中的技能通过 experimental_sync 处理。
 */
export async function runInstallFromLock(args: string[]): Promise<void> {
  const { options: addOpts } = parseAddOptions(args);
  const projectRoot = resolveCliProjectRoot(addOpts.cwd);
  const lock = await readLocalLock(projectRoot);
  const skillEntries = Object.entries(lock.skills);

  if (skillEntries.length === 0) {
    p.log.warn('No project skills found in skills-lock.json');
    p.log.info(
      `Add project-level skills with ${pc.cyan('npx skills add <package>')} (without ${pc.cyan('-g')})`
    );
    return;
  }

  // 仅安装到 .agents/skills/（通用 agent）
  const universalAgentNames = getUniversalAgents();

  // 区分 node_modules 技能与远程技能
  const nodeModuleSkills: string[] = [];
  const bySource = new Map<string, { sourceType: string; skills: string[] }>();

  for (const [skillName, entry] of skillEntries) {
    if (entry.sourceType === 'node_modules') {
      nodeModuleSkills.push(skillName);
      continue;
    }

    const installSource = entry.ref ? `${entry.source}#${entry.ref}` : entry.source;
    const existing = bySource.get(installSource);
    if (existing) {
      existing.skills.push(skillName);
    } else {
      bySource.set(installSource, {
        sourceType: entry.sourceType,
        skills: [skillName],
      });
    }
  }

  const remoteCount = skillEntries.length - nodeModuleSkills.length;
  if (remoteCount > 0) {
    p.log.info(
      `Restoring ${pc.cyan(String(remoteCount))} skill${remoteCount !== 1 ? 's' : ''} from skills-lock.json into ${pc.dim('.agents/skills/')}`
    );
  }

  // 按 source 分组安装远程技能
  for (const [source, { skills }] of bySource) {
    try {
      await runAdd([source], {
        skill: skills,
        agent: universalAgentNames,
        yes: true,
        cwd: projectRoot,
      });
    } catch (error) {
      p.log.error(
        `Failed to install from ${pc.cyan(source)}: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }

  // 通过 sync 处理 node_modules 技能
  if (nodeModuleSkills.length > 0) {
    p.log.info(
      `${pc.cyan(String(nodeModuleSkills.length))} skill${nodeModuleSkills.length !== 1 ? 's' : ''} from node_modules`
    );
    try {
      const { options: syncOptions } = parseSyncOptions(args);
      await runSync(args, { ...syncOptions, yes: true, agent: universalAgentNames, cwd: projectRoot });
    } catch (error) {
      p.log.error(
        `Failed to sync node_modules skills: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }
}
