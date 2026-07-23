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
  validateManifestText,
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
  declaredParent: string
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
  identities: Map<string, string>
  idOwners: Map<string, Set<string>>
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
  const state = createLoadState()
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
  }

  for (const entry of entries) {
    if (!('group' in entry)) {
      const selector = [...prefix, entry.id].join('/')
      const owner = `${catalog.path}#${selector}`
      registerPrimaryIdentity(state, entry.uid, owner, 'UID')
      registerPrimaryIdentity(state, selector, owner, 'selector')
      registerIdIdentity(state, entry.id, owner)
      validateCommandCwd(entry, catalog)
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
    if (groupPath.length > maximumCatalogDepth) throw new RunXError(`RunX catalog depth exceeds ${maximumCatalogDepth} at ${catalog.path}:${groupPath.join('/')}.`, 3)
    const groupSelector = groupPath.join('/')
    state.groups[groupSelector] = { summary: entry.summary, catalogPath: catalog.path, catalogSource: catalog.source }
    const hasCommands = entry.commands !== undefined
    const hasRunX = entry.runx !== undefined
    if (hasCommands === hasRunX) throw new RunXError(`Group "${groupSelector}" must define exactly one of commands or runx.`, 3)
    if (entry.commands) {
      await expandEntries(catalog, entry.commands, groupPath, state, depth + 1)
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
      declaredParent: declaredParent.path,
    })
    await expandCatalog(child, groupPath, state, depth + 1)
  }
}

function validateSourceManifest(catalog: LoadedCatalog): void {
  if (catalog.manifest.version.split('.', 1)[0] !== '2') {
    throw new RunXError(`Unsupported RunX manifest version "${catalog.manifest.version}" in ${catalog.path}; version 2 is required.`, 3)
  }
  validateScriptsDirectory(catalog)
  for (const entry of catalog.manifest.commands) {
    const name = 'group' in entry ? entry.group : entry.id
    if (name === catalog.manifest.namespace) {
      throw new RunXError(`Namespace "${catalog.manifest.namespace}" conflicts with a top-level command or group in ${catalog.path}.`, 3)
    }
  }
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
  await expandCatalog(parent, [], createLoadState(), 0)
  const references = collectRunXReferences(parent.manifest.commands)
    .map(value => resolveReference(parent, value).path)
  if (!references.includes(catalog.path)) {
    throw new RunXError(`Parent catalog ${parent.path} does not declare child ${catalog.path}.`, 3)
  }
}

function createLoadState(): LoadState {
  return {
    active: new Set(), commands: [], groups: {}, children: [], identities: new Map(), idOwners: new Map(),
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
  if (!response.body) throw new RunXError(`Foreign RunX catalog has no response body: ${reference.path}`, 3)
  const reader = response.body.getReader()
  const decoder = new TextDecoder('utf-8', { fatal: true })
  let received = 0
  let text = ''
  try {
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      received += value.byteLength
      if (received > maximumCatalogBytes) {
        await reader.cancel('RunX catalog size limit exceeded')
        throw new RunXError(`Foreign RunX catalog exceeds ${maximumCatalogBytes} bytes: ${reference.path}`, 3)
      }
      text += decoder.decode(value, { stream: true })
    }
    text += decoder.decode()
  } catch (error) {
    if (error instanceof RunXError) throw error
    throw new RunXError(`Could not read foreign RunX catalog ${reference.path}: ${error instanceof Error ? error.message : String(error)}.`, 3)
  }
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

function validateManifestText(text: string, path: string): void {
  validateSourceManifest(decodeCatalog(text, path, 'local', directoryName(path)))
}

function resolveReference(owner: LoadedCatalog, reference: string): { path: string, source: CatalogSource } {
  if (/^https:\/\//i.test(reference)) return { path: normalizeGitHubUrl(reference), source: 'foreign' }
  if (isAbsolutePath(reference)) throw new RunXError(`RunX catalog references must be relative paths or full GitHub URLs: ${reference}`, 3)
  if (owner.source === 'foreign') {
    const resolved = normalizeGitHubUrl(new URL(reference, owner.path).toString())
    if (githubSourceRoot(resolved) !== githubSourceRoot(owner.path)) {
      throw new RunXError(`Relative foreign RunX catalog reference escapes its GitHub owner/repository/ref root: ${reference}`, 3)
    }
    return { path: resolved, source: 'foreign' }
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

function githubSourceRoot(value: string): string {
  const url = new URL(normalizeGitHubUrl(value))
  const segments = url.pathname.split('/').filter(Boolean)
  return `${url.origin}/${segments.slice(0, 3).join('/')}`
}

function validateCommandCwd(command: SourceCommand, catalog: LoadedCatalog): void {
  const commandCwd = resolvePath(catalog.basePath, command.cwd ?? '.')
  const relative = relativePath(catalog.basePath, commandCwd)
  if (isAbsolutePath(relative) || relative.startsWith('..')) {
    throw new RunXError(`Command ${command.uid} has a cwd outside its catalog directory.`, 3)
  }
}

function registerPrimaryIdentity(state: LoadState, identity: string, owner: string, kind: 'UID' | 'selector'): void {
  const primaryOwner = state.identities.get(identity)
  const shorthandOwners = state.idOwners.get(identity)
  if ((primaryOwner && primaryOwner !== owner) || (shorthandOwners && [...shorthandOwners].some(value => value !== owner))) {
    throw new RunXError(`Command ${kind} "${identity}" conflicts with another command UID, selector, or ID shorthand.`, 3)
  }
  state.identities.set(identity, owner)
}

function registerIdIdentity(state: LoadState, identity: string, owner: string): void {
  const primaryOwner = state.identities.get(identity)
  if (primaryOwner && primaryOwner !== owner) {
    throw new RunXError(`Command ID shorthand "${identity}" conflicts with another command UID or selector.`, 3)
  }
  const owners = state.idOwners.get(identity) ?? new Set<string>()
  owners.add(owner)
  state.idOwners.set(identity, owners)
}
