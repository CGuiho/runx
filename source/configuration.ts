/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { RunXError } from './errors.js'
import { directoryName, homeDirectory, isAbsolutePath, joinPath, relativePath, resolvePath } from './path-utils.js'

import type { Static } from '@sinclair/typebox'

export {
  commandSchema,
  findConfiguration,
  manifestSchema,
  readManifest,
  resolveCommand,
}
export type {
  RunXCommand,
  RunXManifest,
}

const identifier = '^[a-z][a-z0-9-]*$'
const semver = '^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$'

const commandSchema = Type.Object({
  uid: Type.String({ pattern: identifier }),
  id: Type.String({ pattern: identifier }),
  group: Type.String({ pattern: identifier }),
  summary: Type.String({ minLength: 1 }),
  description: Type.String({ minLength: 1 }),
  command: Type.String({ minLength: 1 }),
  cwd: Type.Optional(Type.String({ minLength: 1 })),
  shell: Type.Optional(Type.Union([
    Type.Literal('auto'),
    Type.Literal('bash'),
    Type.Literal('sh'),
    Type.Literal('powershell'),
    Type.Literal('cmd'),
  ])),
  tags: Type.Optional(Type.Array(Type.String({ minLength: 1 }))),
  confirm: Type.Optional(Type.Union([Type.Literal('never'), Type.Literal('always')])),
}, { additionalProperties: false })

const manifestSchema = Type.Object({
  version: Type.String({ pattern: semver }),
  project: Type.Optional(Type.Object({ name: Type.String({ minLength: 1 }) }, { additionalProperties: false })),
  scripts: Type.Object({ directory: Type.String({ minLength: 1 }) }, { additionalProperties: false }),
  groups: Type.Record(Type.String({ pattern: identifier }), Type.Object({ summary: Type.String({ minLength: 1 }) }, { additionalProperties: false })),
  commands: Type.Array(commandSchema),
}, { additionalProperties: false })

type RunXCommand = Static<typeof commandSchema>
type RunXManifest = Static<typeof manifestSchema>

type ResolvedCommand = RunXCommand & {
  index: number
  selector: string
  manifestPath: string
  cwd: string
}

async function findConfiguration(cwd: string, explicitConfig?: string): Promise<string> {
  const candidates = explicitConfig
    ? [resolvePath(cwd, explicitConfig)]
    : [joinPath(resolvePath(cwd), 'runx.yaml'), joinPath(homeDirectory(), '.guiho', 'runx', 'runx.yaml')]
  for (const candidate of candidates) if (await Bun.file(candidate).exists()) return candidate
  throw new RunXError(
    explicitConfig
      ? `Configuration file not found: ${candidates[0]}`
      : `No runx.yaml found in ${resolvePath(cwd)} or ${joinPath(homeDirectory(), '.guiho', 'runx', 'runx.yaml')}.`,
    3,
  )
}

async function readManifest(cwd: string, explicitConfig?: string): Promise<{ manifest: RunXManifest, path: string }> {
  const path = await findConfiguration(cwd, explicitConfig)
  let input: unknown
  try {
    input = Bun.YAML.parse(await Bun.file(path).text())
  } catch (error) {
    throw new RunXError(`Invalid YAML in ${path}: ${error instanceof Error ? error.message : 'Unknown YAML parsing error.'}`, 3)
  }
  let parsed: RunXManifest
  try {
    parsed = Value.Decode(manifestSchema, input)
  } catch (error) {
    throw new RunXError(`Invalid RunX configuration ${path}: ${error instanceof Error ? error.message : 'schema validation failed'}`, 3)
  }
  validateManifestSemantics(parsed, path)
  return { manifest: parsed, path }
}

function resolveCommand(manifest: RunXManifest, manifestPath: string, selector: string): ResolvedCommand {
  let command = manifest.commands.find((entry) => entry.uid === selector)
  if (!command && selector.includes('/')) command = manifest.commands.find((entry) => `${entry.group}/${entry.id}` === selector)
  if (!command && /^\d+$/.test(selector)) command = manifest.commands[Number.parseInt(selector, 10) - 1]
  if (!command) {
    const matches = manifest.commands.filter((entry) => entry.id === selector)
    if (matches.length > 1) throw new RunXError(`Ambiguous command ID "${selector}". Use a UID or group/id selector.`, 2)
    command = matches[0]
  }
  if (!command) throw new RunXError(`Unknown command selector: ${selector}`, 2)
  const root = directoryName(manifestPath)
  const commandCwd = resolvePath(root, command.cwd ?? '.')
  const relative = relativePath(root, commandCwd)
  if (isAbsolutePath(relative) || relative.startsWith('..')) throw new RunXError(`Command ${command.uid} has a cwd outside the configuration directory.`, 3)
  return { ...command, index: manifest.commands.indexOf(command) + 1, selector: `${command.group}/${command.id}`, manifestPath, cwd: commandCwd }
}

function validateManifestSemantics(manifest: RunXManifest, path: string): void {
  if (manifest.version.split('.', 1)[0] !== '1') throw new RunXError(`Unsupported RunX manifest version "${manifest.version}" in ${path}.`, 3)
  if (!manifest.groups.public) throw new RunXError(`RunX configuration ${path} must define the "public" group.`, 3)
  const root = directoryName(path)
  const scripts = resolvePath(root, manifest.scripts.directory)
  const relative = relativePath(root, scripts)
  if (relative === '.' || isAbsolutePath(relative) || relative.startsWith('..')) {
    throw new RunXError(`Scripts directory must be a relative subdirectory in ${path}.`, 3)
  }
  const uids = new Set<string>()
  const selectors = new Set<string>()
  for (const command of manifest.commands) {
    if (!manifest.groups[command.group]) throw new RunXError(`Command ${command.uid} references unknown group "${command.group}".`, 3)
    if (uids.has(command.uid)) throw new RunXError(`Duplicate command UID "${command.uid}".`, 3)
    const selector = `${command.group}/${command.id}`
    if (selectors.has(selector)) throw new RunXError(`Duplicate command selector "${selector}".`, 3)
    uids.add(command.uid)
    selectors.add(selector)
  }
}
