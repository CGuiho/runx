/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { RunXError } from './errors.js'
import { directoryName, homeDirectory, joinPath, resolvePath } from './path-utils.js'
import { ensureDirectory, pathExists, readTextIfExists, removePath, writeTextFile, writeTextFileAtomic } from './storage.js'

import type { AgentScope } from './types.js'

export {
  applyAgentInstructions,
  installAgentSkill,
  listAgentPrompts,
  listAgentSkills,
  maintainAgentIntegration,
  removeAgentInstructions,
  showAgentInstructions,
  showAgentPrompt,
  showAgentSkill,
  uninstallAgentSkill,
  updateAgentInstructions,
  updateAgentSkill,
}

type EmbeddedResources = { skill: string, prompt: string } | undefined
type AgentMaintenanceResult = {
  readonly skills: string[]
  readonly instructions: string[]
}

declare global {
  var __RUNX_EMBEDDED_RESOURCES__: EmbeddedResources
}

const skillId = 'guiho-s-runx'
const promptId = 'guiho-i-runx'
const managedStart = '<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->'
const managedEnd = '<!-- END RUNX -->'
const mojibakeManagedStart = '<!-- BEGIN RUNX \u00e2\u20ac\u201d DO NOT EDIT THIS SECTION -->'
const legacyManagedStart = '<!-- BEGIN RUNX AGENT INSTRUCTIONS -->'
const legacyManagedEnd = '<!-- END RUNX AGENT INSTRUCTIONS -->'

async function installAgentSkill(scope: AgentScope, cwd: string): Promise<string[]> {
  return writeSkillTargets(scope, cwd)
}

async function updateAgentSkill(scope: AgentScope, cwd: string): Promise<string[]> {
  return writeSkillTargets(scope, cwd)
}

async function uninstallAgentSkill(scope: AgentScope, cwd: string): Promise<string[]> {
  const targets = skillDirectories(scope, cwd)
  for (const target of targets) await removePath(target)
  return targets
}

function listAgentSkills(filter?: string): Array<{ id: string, description: string }> {
  const skills = [{ id: skillId, description: 'Inspect, validate, describe, and safely execute RunX command catalogs.' }]
  const keyword = filter?.trim().toLowerCase()
  return keyword ? skills.filter((skill) => `${skill.id} ${skill.description}`.toLowerCase().includes(keyword)) : skills
}

async function showAgentSkill(id: string): Promise<{ id: string, path: string, description: string, metadata: Record<string, string> }> {
  if (id !== skillId) throw new RunXError(`Unknown RunX skill: ${id}`, 2)
  return {
    id,
    path: `skills/${skillId}/SKILL.md`,
    description: listAgentSkills()[0]!.description,
    metadata: { version: '0.1.0', package: '@guiho/runx' },
  }
}

async function applyAgentInstructions(cwd: string): Promise<string[]> {
  return reconcileInstructions(cwd, 'apply')
}

async function updateAgentInstructions(cwd: string): Promise<string[]> {
  return reconcileInstructions(cwd, 'update')
}

async function removeAgentInstructions(cwd: string): Promise<string[]> {
  const changed: string[] = []
  for (const path of await instructionTargets(cwd)) {
    const existing = await readTextIfExists(path)
    if (existing === null) continue
    const next = removeManagedBlock(existing)
    if (next !== existing) {
      await writeTextFile(path, next)
      changed.push(path)
    }
  }
  return changed
}

async function showAgentInstructions(): Promise<string> {
  return instructionBlock()
}

function listAgentPrompts(namesOnly: true): string[]
function listAgentPrompts(namesOnly?: false): Array<{ id: string, description: string }>
function listAgentPrompts(namesOnly = false): string[] | Array<{ id: string, description: string }> {
  return namesOnly
    ? [promptId]
    : [{ id: promptId, description: 'Guide an agent through safe RunX catalog inspection and execution.' }]
}

async function showAgentPrompt(id: string): Promise<string> {
  if (id !== promptId) throw new RunXError(`Unknown RunX prompt: ${id}`, 2)
  return readBundledPrompt()
}

async function maintainAgentIntegration(cwd: string): Promise<AgentMaintenanceResult> {
  const skill = await readBundledSkill()
  const skills: string[] = []
  for (const directory of skillDirectories('global', cwd)) {
    const path = joinPath(directory, 'SKILL.md')
    if (await readTextIfExists(path) === skill) continue
    await writeTextFileAtomic(path, skill)
    skills.push(path)
  }

  const instructionPath = await nearestAgentsPath(cwd)
  const existing = await readTextIfExists(instructionPath) ?? ''
  const next = replaceManagedBlock(existing, instructionBlock())
  const instructions: string[] = []
  if (next !== existing) {
    await writeTextFileAtomic(instructionPath, next)
    instructions.push(instructionPath)
  }

  return { skills, instructions }
}

async function writeSkillTargets(scope: AgentScope, cwd: string): Promise<string[]> {
  const skill = await readBundledSkill()
  const installed: string[] = []
  for (const directory of skillDirectories(scope, cwd)) {
    const path = joinPath(directory, 'SKILL.md')
    await ensureDirectory(directory)
    await writeTextFile(path, skill)
    installed.push(path)
  }
  return installed
}

function skillDirectories(scope: AgentScope, cwd: string): string[] {
  const root = scope === 'global' ? homeDirectory() : resolvePath(cwd)
  return [
    joinPath(root, '.agents', 'skills', skillId),
    joinPath(root, '.claude', 'skills', skillId),
  ]
}

async function instructionTargets(cwd: string): Promise<string[]> {
  const root = resolvePath(cwd)
  const agents = joinPath(root, 'AGENTS.md')
  const claude = joinPath(root, 'CLAUDE.md')
  const agentsExists = await pathExists(agents)
  const claudeExists = await pathExists(claude)
  if (agentsExists && claudeExists) return [agents, claude]
  if (claudeExists) return [claude]
  return [agents]
}

async function reconcileInstructions(cwd: string, _action: 'apply' | 'update'): Promise<string[]> {
  const changed: string[] = []
  for (const path of await instructionTargets(cwd)) {
    const existing = await readTextIfExists(path) ?? ''
    const next = replaceManagedBlock(existing, instructionBlock())
    if (next !== existing) {
      await writeTextFile(path, next)
      changed.push(path)
    }
  }
  return changed
}

function instructionBlock(): string {
  return `${managedStart}
## RunX Command Catalog

Load the \`guiho-s-runx\` skill whenever discovering commands, creating or
updating catalog entries, validating \`runx.yaml\`, inspecting command details,
or executing RunX commands.
Start with \`runx check --format json\` and \`runx list --format json\`, select
stable UIDs, use \`runx describe <uid>\`, and run
\`runx run <uid> --dry-run\` before unfamiliar or side-effecting work.
${managedEnd}
`
}

function replaceManagedBlock(existing: string, block: string): string {
  const stripped = removeKnownManagedBlocks(existing).trimEnd()
  return `${stripped}${stripped ? '\n\n' : ''}${block}`
}

function removeManagedBlock(existing: string): string {
  const next = removeKnownManagedBlocks(existing, true).trimEnd()
  return next ? `${next}\n` : ''
}

function removeKnownManagedBlocks(existing: string, includeLeadingWhitespace = false): string {
  let output = existing
  for (const [start, end] of [
    [managedStart, managedEnd],
    [mojibakeManagedStart, managedEnd],
    [legacyManagedStart, legacyManagedEnd],
  ]) {
    const prefix = includeLeadingWhitespace ? '\\s*' : ''
    const pattern = new RegExp(`${prefix}${escapeRegExp(start)}[\\s\\S]*?${escapeRegExp(end)}\\s*`, 'g')
    output = output.replace(pattern, includeLeadingWhitespace ? '\n' : '')
  }
  return output
}

async function nearestAgentsPath(cwd: string): Promise<string> {
  const effectiveCwd = resolvePath(cwd)
  let current = effectiveCwd
  while (true) {
    const candidate = joinPath(current, 'AGENTS.md')
    if (await pathExists(candidate)) return candidate
    const parent = directoryName(current)
    if (parent === current) return joinPath(effectiveCwd, 'AGENTS.md')
    current = parent
  }
}

async function readBundledSkill(): Promise<string> {
  if (globalThis.__RUNX_EMBEDDED_RESOURCES__?.skill) return globalThis.__RUNX_EMBEDDED_RESOURCES__.skill
  const path = new URL('../skills/guiho-s-runx/SKILL.md', import.meta.url)
  if (!(await Bun.file(path).exists())) throw new RunXError('Bundled guiho-s-runx skill is unavailable.', 5)
  return Bun.file(path).text()
}

async function readBundledPrompt(): Promise<string> {
  if (globalThis.__RUNX_EMBEDDED_RESOURCES__?.prompt) return globalThis.__RUNX_EMBEDDED_RESOURCES__.prompt
  const path = new URL('../prompts/guiho-i-runx.md', import.meta.url)
  if (!(await Bun.file(path).exists())) throw new RunXError('Bundled guiho-i-runx prompt is unavailable.', 5)
  return Bun.file(path).text()
}

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}
