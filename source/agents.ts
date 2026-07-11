import { mkdir, readFile, writeFile } from 'node:fs/promises'
import { homedir } from 'node:os'
import { dirname, join } from 'node:path'
import type { AgentScope, AgentTool } from './types.js'
import { RunXError } from './errors.js'

type EmbeddedResources = { skill: string } | undefined

declare global {
  var __RUNX_EMBEDDED_RESOURCES__: EmbeddedResources
}

const managedStart = '<!-- BEGIN RUNX AGENT INSTRUCTIONS -->'
const managedEnd = '<!-- END RUNX AGENT INSTRUCTIONS -->'

export const installAgentSkill = async (scope: AgentScope, tool: AgentTool, cwd: string): Promise<string[]> => {
  const root = scope === 'global' ? homedir() : cwd
  const skill = await readBundledSkill()
  const targets = tool === 'all' ? ['agents', 'claude'] : [tool]
  const installed: string[] = []

  for (const target of targets) {
    const directory = target === 'agents'
      ? join(root, '.agents', 'skills', 'guiho-s-runx')
      : join(root, '.claude', 'skills', 'guiho-s-runx')
    const path = join(directory, 'SKILL.md')
    await mkdir(directory, { recursive: true })
    await writeFile(path, skill, 'utf8')
    installed.push(path)
  }

  return installed
}

export const installAgentInstructions = async (cwd: string): Promise<string> => {
  const path = join(cwd, 'AGENTS.md')
  let existing = ''
  try {
    existing = await readFile(path, 'utf8')
  } catch (error) {
    if (!(error instanceof Error && 'code' in error && error.code === 'ENOENT')) throw error
  }

  const section = `${managedStart}\n## RunX Command Catalog\n\nUse the bundled \`guiho-s-runx\` skill whenever working with \`runx.yaml\` manifests. Inspect with \`runx check --format json\` and \`runx list --format json\` before selecting a command. Use a UID for automation and require explicit developer approval before \`--yes\` on confirmation-gated commands.\n${managedEnd}`
  const pattern = new RegExp(`${escapeRegExp(managedStart)}[\\s\\S]*?${escapeRegExp(managedEnd)}`, 'g')
  const next = pattern.test(existing) ? existing.replace(pattern, section) : `${existing.trimEnd()}${existing.trim() ? '\n\n' : ''}${section}\n`

  await mkdir(dirname(path), { recursive: true })
  await writeFile(path, next, 'utf8')
  return path
}

const readBundledSkill = async (): Promise<string> => {
  if (globalThis.__RUNX_EMBEDDED_RESOURCES__?.skill) return globalThis.__RUNX_EMBEDDED_RESOURCES__.skill
  const path = new URL('../skills/guiho-s-runx/SKILL.md', import.meta.url)
  try {
    return await Bun.file(path).text()
  } catch {
    throw new RunXError('Bundled guiho-s-runx skill is unavailable. Reinstall RunX from an official release.')
  }
}

const escapeRegExp = (value: string): string => value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
