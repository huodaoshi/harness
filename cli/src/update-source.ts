export interface UpdateSourceEntry {
  source: string;
  sourceUrl: string;
  ref?: string;
  skillPath?: string;
}

export interface LocalUpdateSourceEntry {
  source: string;
  ref?: string;
  skillPath?: string;
}

export function formatSourceInput(sourceUrl: string, ref?: string): string {
  if (!ref) {
    return sourceUrl;
  }
  return `${sourceUrl}#${ref}`;
}

/**
 * 从以 SKILL.md 结尾的 skillPath 推导技能所在文件夹路径。
 * 技能位于仓库根目录时返回 ''。
 */
function deriveSkillFolder(skillPath: string): string {
  let folder = skillPath;
  if (folder.endsWith('/SKILL.md')) {
    folder = folder.slice(0, -9);
  } else if (folder.endsWith('SKILL.md')) {
    folder = folder.slice(0, -8);
  }
  if (folder.endsWith('/')) {
    folder = folder.slice(0, -1);
  }
  return folder;
}

function appendFolderAndRef(source: string, skillPath: string, ref?: string): string {
  const folder = deriveSkillFolder(skillPath);
  const withFolder = folder ? `${source}/${folder}` : source;
  return ref ? `${withFolder}#${ref}` : withFolder;
}

/**
 * 在 update 流程中为 `skills add` 构建 source 参数。
 * 对路径定向的更新使用简写形式，避免分支/路径歧义。
 */
export function buildUpdateInstallSource(entry: UpdateSourceEntry): string {
  if (!entry.skillPath) {
    return formatSourceInput(entry.sourceUrl, entry.ref);
  }
  return appendFolderAndRef(entry.source, entry.skillPath, entry.ref);
}

/**
 * 在项目级 update 中为 `skills add` 构建 source 参数。
 * 本地 lock 条目不含 `sourceUrl`，无 `skillPath` 时回退为裸 `source` 标识符。
 */
export function buildLocalUpdateSource(entry: LocalUpdateSourceEntry): string {
  if (!entry.skillPath) {
    return formatSourceInput(entry.source, entry.ref);
  }
  return appendFolderAndRef(entry.source, entry.skillPath, entry.ref);
}
