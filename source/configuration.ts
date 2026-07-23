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
  CatalogChild,
  CatalogSource,
  RunXCommand,
  RunXManifest,
}

const identifier = '^[a-z][a-z0-9-]*$'
const semver = '^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)(?:-((?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\\.(?:0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\\+([0-9a-zA-Z-]+(?:\\.[0-9a-zA-Z-]+)*))?$'
const maximumCatalogBytes = 1024 * 1024
const maximumCatalogDepth = 32

const commandSchema = Type.Object({
  uid: Type.String({ pattern: identifier }),
  id: Type.String({ pattern: identifier }),
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

const groupSchema = Type.Recursive(This => Type.Object({
  group: Type.String({ pattern: identifier }),
  summary: Type.String({ minLength: 1 }),
  commands: Type.Optional(Type.Array(Type.Union([commandSchema, This]))),
  runx: Type.Optional(Type.String({ minLength: 1 })),
}, { additionalProperties: false }))

const manifestSchema = Type.Object({
  version: Type.String({ pattern: semver }),
  namespace: Type.String({ pattern: identifier }),
  scripts: Type.Object({ directory: Type.String({ minLength: 1 }) }, { additionalProperties: false }),
  parent: Type.Optional(Type.String({ minLength: 1 })),
  commands: Type.Array(Type.Union([commandSchema, groupSchema])),
}, { additionalProperties: false })

type SourceCommand = Static<typeof commandSchema>
type SourceGroup = Static<typeof groupSchema>
type SourceManifest = Static<typeof manifestSchema>
type CatalogSource = 'local' | 'foreign'

type RunXCommand = SourceCommand & {
  group: string
  selector?: string
  catalogPath?: string
  catalogSource?: CatalogSource
  basePath?: string
}

type CatalogChild = {
  namespace: string
  declaredNamespace: string
  path: string
  source: CatalogSource
  parent: string
}

type RunXManifest = {
  version: string
  namespace: string
  scripts: { directory: string }
  parent?: { path: string, source: CatalogSource }
  commands: RunXCommand[]
  groups: Record<string, { summary: string, catalogPath: string, catalogSource: CatalogSource }>
  children: CatalogChild[]
}

type ResolvedCommand = RunXCommand & {
  index: number
  selector: string
  manifestPath: string
  cwd: string
}

type LoadedCatalog = {
  manifest: SourceManifest
  path: string
  source: CatalogSource
  basePath: string
}

type LoadState = {
  active: Set<string>
  commands: RunXCommand[]
  groups: RunXManifest['groups']
  children: CatalogChild[]
  uids: Set<string>
  selectors: Set<string>
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
  const root = await loadLocalCatalog(path)
  const state: LoadState = {
    active: new Set(), commands: [], groups: {}, children: [], uids: new Set(), selectors: new Set(),
  }
  validateSourceManifest(root)
  if (root.manifest.parent) await validateDeclaredParent(root)
  await expandCatalog(root, [], state, 0)
  const manifest: RunXManifest = {
    version: root.manifest.version,
    namespace: root.manifest.namespace,
    scripts: root.manifest.scripts,
    ...(root.manifest.parent ? { parent: referenceMetadata(root, root.manifest.parent) } : {}),
    commands: state.commands,
    groups: state.groups,
    children: state.children,
  }
  return { manifest, path }
}

function resolveCommand(manifest: RunXManifest, manifestPath: string, selector: string): ResolvedCommand {
  let command = manifest.commands.find(entry => entry.uid === selector)
  if (!command && selector.includes('/')) command = manifest.commands.find(entry => commandSelector(entry) === selector)
  if (!command && /^\\d+$/.test(selector)) command = manifest.commands[Number.parseInt(selector, 10) - 1]
  if (!command) {
    const matches = manifest.commands.filter(entry => entry.id === selector)
    if (matches.length > 1) throw new RunXError(`Ambiguous command ID "${selector}". Use a UID or full selector.`, 2)
    command = matches[0]
  }
  if (!command) throw new RunXError(`Unknown command selector: ${selector}`, 2)
  const root = command.basePath ?? directoryName(manifestPath)
  const commandCwd = resolvePath(root, command.cwd ?? '.')
  const relative = relativePath(root, commandCwd)
  if (isAbsolutePath(relative) || relative.startsWith('..')) throw new RunXError(`Command ${command.uid} has a cwd outside its catalog directory.`, 3)
  return {
    ...command,
    index: manifest.commands.indexOf(command) + 1,
    selector: commandSelector(command),
    manifestPath: command.catalogPath ?? manifestPath,
    cwd: commandCwd,
  }
}

async function expandCatalog(catalog: LoadedCatalog, prefix: string[], state: LoadState, depth: number): Promise<void> {
  if (depth > maximumCatalogDepth) throw new RunXError(`RunX catalog depth exceeds ${maximumCatalogDepth} at ${catalog.path}.`, 3)
  if (state.active.has(catalog.path)) throw new RunXError(`RunX catalog cycle detected at ${catalog.path}.`, 3)
  state.active.add(catalog.path)
  try {
    await expandEntries(catalog, catalog.manifest.commands, prefix, state, depth)
  } finally {
    state.active.delete(catalog.path)
  }
}

async function expandEntries(catalog: LoadedCatalog, entries: Array<SourceCommand | SourceGroup>, prefix: string[], state: LoadState, depth: number): Promise<void> {
  const siblingNames = new Set<string>()
  for (const entry of entries) {
    const name = 'group' in entry ? entry.group : entry.id
    if (siblingNames.has(name)) throw new RunXError(`Duplicate command or group name "${name}" at ${catalog.path}:${prefix.join('/') || catalog.manifest.namespace}.`, 3)
    siblingNames.add(name)
    if (prefix.length === 0 && name === catalog.manifest.namespace) {
      throw new RunXError(`Namespace "${catalog.manifest.namespace}" conflicts with a top-level command or group in ${catalog.path}.`, 3)
    }
  }

  for (const entry of entries) {
    if (!('group' in entry)) {
      const selector = [...prefix, entry.id].join('/')
      if (state.uids.has(entry.uid)) throw new RunXError(`Duplicate command UID "${entry.uid}" in composed catalog.`, 3)
      if (state.selectors.has(selector)) throw new RunXError(`Duplicate command selector "${selector}" in composed catalog.`, 3)
      state.uids.add(entry.uid)
      state.selectors.add(selector)
      state.commands.push({
        ...entry,
        group: prefix.join('/'),
        selector,
        catalogPath: catalog.path,
        catalogSource: catalog.source,
        basePath: catalog.basePath,
      })
      continue
    }

    const groupPath = [...prefix, entry.group]
    const groupSelector = groupPath.join('/')
    state.groups[groupSelector] = { summary: entry.summary, catalogPath: catalog.path, catalogSource: catalog.source }
    const hasCommands = entry.commands !== undefined
    const hasRunX = entry.runx !== undefined
    if (hasCommands === hasRunX) throw new RunXError(`Group "${groupSelector}" must define exactly one of commands or runx.`, 3)
    if (entry.commands) {
      await expandEntries(catalog, entry.commands, groupPath, state, depth)
      continue
    }

    const child = await loadReferencedCatalog(catalog, entry.runx!)
    validateSourceManifest(child)
    if (!child.manifest.parent) throw new RunXError(`Child catalog ${child.path} must declare its parent.`, 3)
    const declaredParent = resolveReference(child, child.manifest.parent)
    if (declaredParent.path !== catalog.path) {
      if (child.source === 'foreign' && catalog.source === 'local') {
        await validateDeclaredParent(child)
      } else {
        throw new RunXError(`Child catalog ${child.path} declares parent ${declaredParent.path}; expected ${catalog.path}.`, 3)
      }
    }
    state.children.push({
      namespace: entry.group,
      declaredNamespace: child.manifest.namespace,
      path: child.path,
      source: child.source,
      parent: catalog.path,
    })
    await expandCatalog(child, groupPath, state, depth + 1)
  }
}

function validateSourceManifest(catalog: LoadedCatalog): void {
  if (catalog.manifest.version.split('.', 1)[0] !== '2') {
    throw new RunXError(`Unsupported RunX manifest version "${catalog.manifest.version}" in ${catalog.path}; version 2 is required.`, 3)
  }
  if (catalog.source === 'local') validateScriptsDirectory(catalog)
}

function validateScriptsDirectory(catalog: LoadedCatalog): void {
  const scripts = resolvePath(catalog.basePath, catalog.manifest.scripts.directory)
  const relative = relativePath(catalog.basePath, scripts)
  if (relative === '.' || isAbsolutePath(relative) || relative.startsWith('..')) {
    throw new RunXError(`Scripts directory must be a relative subdirectory in ${catalog.path}.`, 3)
  }
}

async function validateDeclaredParent(catalog: LoadedCatalog): Promise<void> {
  const reference = resolveReference(catalog, catalog.manifest.parent!)
  const parent = await loadResolvedCatalog(reference, catalog.basePath)
  validateSourceManifest(parent)
  const references = collectRunXReferences(parent.manifest.commands)
    .map(value => resolveReference(parent, value).path)
  if (!references.includes(catalog.path)) {
    throw new RunXError(`Parent catalog ${parent.path} does not declare child ${catalog.path}.`, 3)
  }
}

function collectRunXReferences(entries: Array<SourceCommand | SourceGroup>): string[] {
  const references: string[] = []
  for (const entry of entries) {
    if (!('group' in entry)) continue
    if (entry.runx) references.push(entry.runx)
    if (entry.commands) references.push(...collectRunXReferences(entry.commands))
  }
  return references
}

async function loadLocalCatalog(path: string): Promise<LoadedCatalog> {
  return decodeCatalog(await Bun.file(path).text(), path, 'local', directoryName(path))
}

async function loadReferencedCatalog(owner: LoadedCatalog, reference: string): Promise<LoadedCatalog> {
  const resolved = resolveReference(owner, reference)
  return loadResolvedCatalog(resolved, owner.basePath)
}

async function loadResolvedCatalog(reference: { path: string, source: CatalogSource }, localBase: string): Promise<LoadedCatalog> {
  if (reference.source === 'local') {
    if (!await Bun.file(reference.path).exists()) throw new RunXError(`Referenced RunX catalog not found: ${reference.path}`, 3)
    return loadLocalCatalog(reference.path)
  }
  let response: Response
  try {
    response = await fetch(reference.path, { signal: AbortSignal.timeout(10_000), headers: { Accept: 'text/yaml, text/plain' } })
  } catch (error) {
    throw new RunXError(`Could not load foreign RunX catalog ${reference.path}: ${error instanceof Error ? error.message : String(error)}.`, 3)
  }
  if (response.url) normalizeGitHubUrl(response.url)
  if (!response.ok) throw new RunXError(`Could not load foreign RunX catalog ${reference.path}: HTTP ${response.status}.`, 3)
  const declaredLength = Number(response.headers.get('content-length') ?? '0')
  if (declaredLength > maximumCatalogBytes) throw new RunXError(`Foreign RunX catalog exceeds ${maximumCatalogBytes} bytes: ${reference.path}`, 3)
  const text = await response.text()
  if (new TextEncoder().encode(text).byteLength > maximumCatalogBytes) throw new RunXError(`Foreign RunX catalog exceeds ${maximumCatalogBytes} bytes: ${reference.path}`, 3)
  return decodeCatalog(text, reference.path, 'foreign', localBase)
}

function decodeCatalog(text: string, path: string, source: CatalogSource, basePath: string): LoadedCatalog {
  let input: unknown
  try {
    input = Bun.YAML.parse(text)
  } catch (error) {
    throw new RunXError(`Invalid YAML in ${path}: ${error instanceof Error ? error.message : 'Unknown YAML parsing error.'}`, 3)
  }
  try {
    return { manifest: Value.Decode(manifestSchema, input), path, source, basePath }
  } catch (error) {
    throw new RunXError(`Invalid RunX configuration ${path}: ${error instanceof Error ? error.message : 'schema validation failed'}`, 3)
  }
}

function resolveReference(owner: LoadedCatalog, reference: string): { path: string, source: CatalogSource } {
  if (/^https:\/\//i.test(reference)) return { path: normalizeGitHubUrl(reference), source: 'foreign' }
  if (isAbsolutePath(reference)) throw new RunXError(`RunX catalog references must be relative paths or full GitHub URLs: ${reference}`, 3)
  if (owner.source === 'foreign') {
    return { path: normalizeGitHubUrl(new URL(reference, owner.path).toString()), source: 'foreign' }
  }
  return { path: resolvePath(directoryName(owner.path), reference), source: 'local' }
}

function referenceMetadata(owner: LoadedCatalog, reference: string): { path: string, source: CatalogSource } {
  return resolveReference(owner, reference)
}

function normalizeGitHubUrl(value: string): string {
  const url = new URL(value)
  if (url.protocol !== 'https:') throw new RunXError(`Foreign RunX catalog URL must use HTTPS: ${value}`, 3)
  if (url.hostname === 'github.com') {
    const match = url.pathname.match(/^\/([^/]+)\/([^/]+)\/blob\/([^/]+)\/(.+)$/)
    if (!match) throw new RunXError(`GitHub RunX catalog URL must use /owner/repository/blob/ref/path: ${value}`, 3)
    return `https://raw.githubusercontent.com/${match[1]}/${match[2]}/${match[3]}/${match[4]}`
  }
  if (url.hostname !== 'raw.githubusercontent.com') throw new RunXError(`Foreign RunX catalogs must use a full GitHub URL: ${value}`, 3)
  if (!/^\/[^/]+\/[^/]+\/[^/]+\/.+/.test(url.pathname)) throw new RunXError(`Invalid raw GitHub RunX catalog URL: ${value}`, 3)
  url.search = ''
  url.hash = ''
  return url.toString()
}

function commandSelector(command: RunXCommand): string {
  return command.selector ?? [command.group, command.id].filter(Boolean).join('/')
}
