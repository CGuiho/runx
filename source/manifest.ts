import { dirname, isAbsolute, relative, resolve } from 'node:path'
import { Type } from '@sinclair/typebox'
import { TypeCompiler } from '@sinclair/typebox/compiler'
import { RunXError, invariant } from './errors.js'
import type { ResolvedCommand, RunXCommand, RunXManifest } from './types.js'

const identifier = '^[a-z][a-z0-9-]*$'

export const CommandSchema = Type.Object({
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

export const ManifestSchema = Type.Object({
  version: Type.Literal(1),
  project: Type.Optional(Type.Object({ name: Type.String({ minLength: 1 }) }, { additionalProperties: false })),
  groups: Type.Record(Type.String({ pattern: identifier }), Type.Object({ summary: Type.String({ minLength: 1 }) }, { additionalProperties: false })),
  commands: Type.Array(CommandSchema, { minItems: 1 }),
}, { additionalProperties: false })

const manifestCheck = TypeCompiler.Compile(ManifestSchema)

export const findManifest = async (cwd: string, explicitFile?: string): Promise<string> => {
  if (explicitFile) {
    const candidate = resolve(cwd, explicitFile)
    if (await Bun.file(candidate).exists()) return candidate
    throw new RunXError(`Manifest not found: ${candidate}`)
  }

  let current = resolve(cwd)
  while (true) {
    const candidate = resolve(current, 'runx.yaml')
    if (await Bun.file(candidate).exists()) return candidate
    const parent = dirname(current)
    if (parent === current) break
    current = parent
  }

  throw new RunXError('No runx.yaml found. Run this command inside a configured project or pass --file <path>.')
}

export const readManifest = async (cwd: string, explicitFile?: string): Promise<{ manifest: RunXManifest, path: string }> => {
  const path = await findManifest(cwd, explicitFile)
  let parsed: unknown

  try {
    parsed = Bun.YAML.parse(await Bun.file(path).text())
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown YAML parsing error.'
    throw new RunXError(`Invalid YAML in ${path}: ${message}`)
  }

  if (!manifestCheck.Check(parsed)) {
    const errors = [...manifestCheck.Errors(parsed)].map((error) => `${error.path || '/'}: ${error.message}`)
    throw new RunXError(`Invalid RunX manifest ${path}:\n${errors.join('\n')}`)
  }

  const manifest = parsed as RunXManifest
  validateManifestSemantics(manifest, path)
  return { manifest, path }
}

export const resolveCommand = (manifest: RunXManifest, manifestPath: string, selector: string): ResolvedCommand => {
  const commands = manifest.commands
  let command: RunXCommand | undefined = commands.find((entry) => entry.uid === selector)

  if (!command && selector.includes('/')) {
    command = commands.find((entry) => `${entry.group}/${entry.id}` === selector)
  }

  if (!command && /^\d+$/.test(selector)) {
    command = commands[Number.parseInt(selector, 10) - 1]
  }

  if (!command) {
    const matches = commands.filter((entry) => entry.id === selector)
    if (matches.length > 1) {
      throw new RunXError(`Ambiguous command ID "${selector}". Use one of: ${matches.map((entry) => `${entry.group}/${entry.id} (${entry.uid})`).join(', ')}`)
    }
    command = matches[0]
  }

  if (!command) throw new RunXError(`Unknown command selector: ${selector}`)

  const manifestDirectory = dirname(manifestPath)
  const commandCwd = resolve(manifestDirectory, command.cwd ?? '.')
  const relativeCwd = relative(manifestDirectory, commandCwd)
  invariant(!isAbsolute(relativeCwd) && !relativeCwd.startsWith('..'), `Command ${command.uid} has a cwd outside the manifest directory.`)

  return {
    ...command,
    index: commands.indexOf(command) + 1,
    selector: `${command.group}/${command.id}`,
    manifestPath,
    cwd: commandCwd,
  }
}

const validateManifestSemantics = (manifest: RunXManifest, path: string): void => {
  const uids = new Set<string>()
  const selectors = new Set<string>()

  for (const command of manifest.commands) {
    if (!manifest.groups[command.group]) throw new RunXError(`Command ${command.uid} references unknown group "${command.group}" in ${path}.`)
    if (uids.has(command.uid)) throw new RunXError(`Duplicate command UID "${command.uid}" in ${path}.`)
    const selector = `${command.group}/${command.id}`
    if (selectors.has(selector)) throw new RunXError(`Duplicate command selector "${selector}" in ${path}.`)
    uids.add(command.uid)
    selectors.add(selector)
  }
}
