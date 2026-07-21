/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { defineCommand, renderUsage, runCommand as runCittyCommand } from 'citty'
import { agentMaintenanceWorkerCwd, runAgentMaintenanceWorker, spawnAgentMaintenanceWorker } from './agent-maintenance.js'
import {
  applyAgentInstructions,
  installAgentSkill,
  listAgentPrompts,
  listAgentSkills,
  removeAgentInstructions,
  showAgentInstructions,
  showAgentPrompt,
  showAgentSkill,
  uninstallAgentSkill,
  updateAgentInstructions,
  updateAgentSkill,
} from './agents.js'
import { RunXError } from './errors.js'
import { runCommand } from './executor.js'
import { readVersion, renderHelpDocs, renderHelpTree } from './help.js'
import { initializeRunXManifest } from './init.js'
import { readManifest, resolveCommand } from './manifest.js'
import { resolvePath } from './path-utils.js'
import { renderDescription, renderExecutionPlan, renderJson, renderList } from './render.js'
import { fetchReleaseCatalog, paginateReleaseCatalog, resolveUpgradePlatform } from './release-catalog.js'
import { checkForLatestVersion, uninstallSelf, upgradeSelf } from './self-management.js'
import { readCachedUpdateNotice, runUpdateWorker, spawnUpdateWorker } from './update-cache.js'
import { renderReleaseCatalog, renderUpgradeEvent, renderUpgradeHeading, renderUpgradePlan, renderUpgradeResult } from './upgrade-reporting.js'

import type { CommandContext, CommandDef, SubCommandsDef } from 'citty'
import type { AgentScope, CliOptions, OutputFormat } from './types.js'
import type { UpgradeArch, UpgradeVariant } from './upgrade-types.js'

export {
  renderStartupBanner,
  runCli,
  runCliWithErrorHandling,
  runxCommand,
}

type GlobalArgs = {
  _: string[]
  cwd?: string
  config?: string
  format?: string
  verbose?: boolean
  dryRun?: boolean
  yes?: boolean
  local?: boolean
  filter?: string
  names?: boolean
  version?: string | boolean
  arch?: string
  variant?: string
  page?: string
  perPage?: string
  preReleases?: boolean
  helpTree?: boolean
  helpTreeDepth?: string
  helpDocs?: boolean
  help?: boolean
  h?: boolean
  v?: boolean
  [key: string]: string | number | boolean | string[] | undefined
}

type CliState = {
  args: GlobalArgs
  command: CommandDef<any>
  commands: Map<string, CommandDef<any>>
}

class CliHandled extends Error {}

class CliUsageError extends Error {
  readonly usage: string
  constructor(message: string, usage: string) {
    super(message)
    this.name = 'CliUsageError'
    this.usage = usage
  }
}

const formatSchema = Type.Union([Type.Literal('text'), Type.Literal('json')])
const archSchema = Type.Union([Type.Literal('x64'), Type.Literal('arm64')])
const variantSchema = Type.Union([Type.Literal('baseline'), Type.Literal('default'), Type.Literal('modern')])
const positiveIntegerSchema = Type.Integer({ minimum: 1 })
const platformLabels: Readonly<Record<string, string>> = {
  darwin: 'macOS',
  linux: 'Linux',
  win32: 'Windows',
}

function renderStartupBanner(platform = process.platform, version = readVersion()): string {
  return `Hello ${platformLabels[platform] ?? platform} - runx v${version}\n`
}

const helpArgs = {
  help: { type: 'boolean', alias: 'h', description: 'Show command help.' },
  'help-tree': { type: 'boolean', description: 'Show this command hierarchy.' },
  'help-tree-depth': { type: 'string', valueHint: 'positive-integer', description: 'Limit help-tree recursion depth.' },
  'help-docs': { type: 'boolean', description: 'Emit Markdown documentation for this command.' },
} as const

const catalogArgs = {
  cwd: { type: 'string', valueHint: 'path', description: 'Use this effective working directory.' },
  config: { type: 'string', valueHint: 'path', description: 'Use this runx.yaml configuration file.' },
  format: { type: 'string', valueHint: 'text|json', default: 'text', description: 'Select output format.' },
  verbose: { type: 'boolean', description: 'Enable diagnostics.' },
  ...helpArgs,
} as const

function createCommandTree(): { command: CommandDef<any>, state: CliState } {
  const state: CliState = { args: { _: [] }, command: {}, commands: new Map() }
  const leaf = (name: string, description: string, args: Record<string, any>, run: (context: CommandContext<any>) => Promise<void> | void): CommandDef<any> => {
    const command = defineCommand({ meta: { name, description }, args: { ...args, ...helpArgs }, setup: helpSetup(state), run })
    state.commands.set(name.replace(/^runx(?:\s+|$)/, ''), command)
    return command
  }

  const list = leaf('runx list', 'List commands in a RunX configuration.', catalogArgs, async () => listCommands(options(state.args)))
  const describe = leaf('runx describe', 'Describe one catalog command without execution.', {
    selector: { type: 'positional', description: 'UID, group/id, index, or unambiguous ID.' },
    ...catalogArgs,
  }, async ({ args }) => {
    if (!args.selector) throw await usageError(state, 'Missing required positional argument: SELECTOR')
    await describeCommand(String(args.selector), options(state.args))
  })
  const run = leaf('runx run', 'Execute one selected catalog command.', {
    selector: { type: 'positional', description: 'UID, group/id, index, or unambiguous ID.' },
    ...catalogArgs,
    'dry-run': { type: 'boolean', description: 'Print the execution plan without spawning.' },
    yes: { type: 'boolean', description: 'Approve a confirmation-gated command.' },
  }, async ({ args }) => {
    if (!args.selector) throw await usageError(state, 'Missing required positional argument: SELECTOR')
    await runSelected(String(args.selector), options(state.args), state.args)
  })
  const check = leaf('runx check', 'Validate a RunX configuration without execution.', catalogArgs, async () => checkManifest(options(state.args)))
  const init = leaf('runx init', 'Create a new YAML RunX configuration.', catalogArgs, async () => {
    const result = await initializeRunXManifest({ cwd: options(state.args).cwd, config: options(state.args).config })
    write(state.args.format === 'json' ? renderJson(result) : `created: ${result.path}\n`)
  })

  const skillInstall = leaf('runx agent skill install', 'Install the bundled skill into both global tool locations.', {
    cwd: catalogArgs.cwd, local: { type: 'boolean', description: 'Use project-local tool directories.' }, format: catalogArgs.format,
  }, async () => agentMutation('installed', installAgentSkill, state.args))
  const skillUninstall = leaf('runx agent skill uninstall', 'Remove the bundled skill from both tool locations.', {
    cwd: catalogArgs.cwd, local: { type: 'boolean', description: 'Use project-local tool directories.' }, format: catalogArgs.format,
  }, async () => agentMutation('removed', uninstallAgentSkill, state.args))
  const skillUpdate = leaf('runx agent skill update', 'Refresh the bundled skill in both tool locations.', {
    cwd: catalogArgs.cwd, local: { type: 'boolean', description: 'Use project-local tool directories.' }, format: catalogArgs.format,
  }, async () => agentMutation('updated', updateAgentSkill, state.args))
  const skillList = leaf('runx agent skill list', 'List bundled RunX skills.', {
    filter: { type: 'string', valueHint: 'keyword', description: 'Filter skill metadata.' }, format: catalogArgs.format,
  }, () => writeFormatted(listAgentSkills(state.args.filter), state.args))
  const skillShow = leaf('runx agent skill show', 'Show metadata for one bundled skill.', {
    id: { type: 'positional', description: 'Bundled skill ID.' }, format: catalogArgs.format,
  }, async ({ args }) => {
    if (!args.id) throw await usageError(state, 'Missing required positional argument: ID')
    writeFormatted(await showAgentSkill(String(args.id)), state.args)
  })
  const skill = group('runx agent skill', 'Manage the bundled RunX skill.', { install: skillInstall, uninstall: skillUninstall, update: skillUpdate, list: skillList, show: skillShow }, state)

  const instructionApply = leaf('runx agent instruction apply', 'Apply the managed instruction block.', { cwd: catalogArgs.cwd, format: catalogArgs.format }, async () => instructionMutation('updated', applyAgentInstructions, state.args))
  const instructionRemove = leaf('runx agent instruction remove', 'Remove the managed instruction block.', { cwd: catalogArgs.cwd, format: catalogArgs.format }, async () => instructionMutation('removed', removeAgentInstructions, state.args))
  const instructionUpdate = leaf('runx agent instruction update', 'Replace stale managed instruction content.', { cwd: catalogArgs.cwd, format: catalogArgs.format }, async () => instructionMutation('updated', updateAgentInstructions, state.args))
  const instructionShow = leaf('runx agent instruction show', 'Print the raw instruction template.', {}, async () => write(await showAgentInstructions()))
  const instruction = group('runx agent instruction', 'Manage RunX instruction blocks.', { apply: instructionApply, remove: instructionRemove, update: instructionUpdate, show: instructionShow }, state)

  const promptList = leaf('runx agent prompt list', 'List bundled RunX prompts.', {
    names: { type: 'boolean', description: 'Print prompt names only.' }, format: catalogArgs.format,
  }, () => {
    if (state.args.names) {
      const names = listAgentPrompts(true)
      write(options(state.args).format === 'json' ? renderJson(names) : `${names.join('\n')}\n`)
      return
    }
    writeFormatted(listAgentPrompts(false), state.args)
  })
  const promptShow = leaf('runx agent prompt show', 'Print one raw bundled prompt.', {
    id: { type: 'positional', description: 'Bundled prompt ID.' },
  }, async ({ args }) => {
    if (!args.id) throw await usageError(state, 'Missing required positional argument: ID')
    write(await showAgentPrompt(String(args.id)))
  })
  const prompt = group('runx agent prompt', 'Inspect bundled agent prompts.', { list: promptList, show: promptShow }, state)
  const agent = group('runx agent', 'Manage RunX agent integration.', { skill, instruction, prompt }, state)

  const upgradeArgs = {
    version: { type: 'string', valueHint: 'version', description: 'Select an exact release version.' },
    arch: { type: 'string', valueHint: 'x64|arm64', description: 'Select target architecture.' },
    variant: { type: 'string', valueHint: 'baseline|default|modern', description: 'Select x64 binary variant.' },
    'dry-run': { type: 'boolean', description: 'Plan without mutation.' },
    format: catalogArgs.format,
  } as const
  const upgradeCheck = leaf('runx upgrade check', 'Check whether a newer stable release exists.', {
    format: catalogArgs.format,
  }, async () => {
    const result = await checkForLatestVersion()
    writeFormatted(result, state.args)
  })
  const upgradeList = leaf('runx upgrade list', 'List RunX releases newest first.', {
    page: { type: 'string', valueHint: 'positive-integer', description: 'Select result page.' },
    'per-page': { type: 'string', valueHint: 'positive-integer', description: 'Select page size.' },
    'pre-releases': { type: 'boolean', description: 'Accepted explicitly; prereleases are always included.' },
    arch: { type: 'string', valueHint: 'x64|arm64', description: 'Select target architecture.' },
    variant: { type: 'string', valueHint: 'baseline|default|modern', description: 'Select x64 variant.' },
    format: catalogArgs.format,
  }, async () => runUpgradeList(state.args))
  const upgrade = defineCommand({
    meta: { name: 'runx upgrade', description: 'Inspect or upgrade a native RunX executable.' },
    args: { ...upgradeArgs, ...helpArgs },
    subCommands: { check: upgradeCheck, list: upgradeList },
    setup: helpSetup(state),
    run: async () => {
      if (state.command === upgradeCheck || state.command === upgradeList) return
      await runUpgrade(options(state.args), state.args)
    },
  })
  state.commands.set('upgrade', upgrade)
  const uninstall = leaf('runx uninstall', 'Uninstall the native RunX executable.', {
    'dry-run': { type: 'boolean', description: 'Print the target without deleting it.' }, format: catalogArgs.format,
  }, async () => {
    const result = await uninstallSelf(Boolean(state.args.dryRun))
    writeFormatted(result, state.args)
  })

  const root = defineCommand({
    meta: { name: 'runx', version: readVersion(), description: 'A documented local command catalog for runx.yaml.' },
    args: { version: { type: 'boolean', alias: 'v', description: 'Show the RunX version.' }, ...helpArgs },
    subCommands: { list, describe, run, check, init, agent, upgrade, uninstall },
    setup: async (context) => {
      state.args = context.args as GlobalArgs
      state.command = resolveHelpCommand(state)
      await validateInvocation(state)
      if (context.rawArgs.length === 1 && ['-v', '--version'].includes(context.rawArgs[0] ?? '')) {
        return handled(`${readVersion()}\n`)
      }
      await handleHelp(state)
    },
    run: ({ args }) => {
      if ((args._ as string[]).length === 0) write(renderStartupBanner())
    },
  })
  state.commands.set('', root)
  state.command = root
  return { command: root, state }
}

async function validateInvocation(state: CliState): Promise<void> {
  const first = state.args._[0]
  const topLevel = new Set(['list', 'describe', 'run', 'check', 'init', 'agent', 'upgrade', 'uninstall'])
  if (first && !topLevel.has(first)) throw await usageError(state, `Unknown command: ${first}`)
  const allowed = new Set(['_', 'h', 'v'])
  for (const command of state.commands.values()) {
    for (const key of Object.keys(command.args ?? {})) {
      allowed.add(key)
      allowed.add(key.replace(/-([a-z])/g, (_, character: string) => character.toUpperCase()))
    }
  }
  const unknown = Object.keys(state.args).find((key) => !allowed.has(key))
  if (unknown) throw await usageError(state, `Unknown option: --${unknown.replace(/[A-Z]/g, (character) => `-${character.toLowerCase()}`)}`)
}

function group(name: string, description: string, subCommands: SubCommandsDef, state: CliState, defaultCommand?: string): CommandDef<any> {
  const command = defineCommand({ meta: { name, description }, args: helpArgs, subCommands, ...(defaultCommand ? { default: defaultCommand } : {}), setup: helpSetup(state) })
  state.commands.set(name.replace(/^runx(?:\s+|$)/, ''), command)
  return command
}

function helpSetup(state: CliState): (context: CommandContext<any>) => Promise<void> {
  return async (context) => {
    state.args = context.args as GlobalArgs
    state.command = context.cmd
    await handleHelp(state)
  }
}

async function handleHelp(state: CliState): Promise<void> {
  if (state.args.helpDocs) return handled(renderHelpDocs(state.command))
  if (state.args.helpTree || state.args.helpTreeDepth) {
    const depth = state.args.helpTreeDepth ? positiveInteger(state.args.helpTreeDepth, '--help-tree-depth') : Number.POSITIVE_INFINITY
    return handled(renderHelpTree(state.command, depth))
  }
  if (state.args.help || state.args.h) return handled(`${await renderUsage(state.command)}\n`)
}

function resolveHelpCommand(state: CliState): CommandDef<any> {
  const tokens = state.args._.filter((token) => !token.startsWith('-'))
  for (let length = tokens.length; length > 0; length -= 1) {
    const command = state.commands.get(tokens.slice(0, length).join(' '))
    if (command) return command
  }
  return state.commands.get('') ?? state.command
}

const { command: runxCommand } = createCommandTree()

async function runCli(rawArgs: string[] = process.argv.slice(2)): Promise<void> {
  if (rawArgs.length === 1 && rawArgs[0] === '--check-updates-worker') {
    await runUpdateWorker({ leaseToken: Bun.env.RUNX_UPDATE_WORKER_LEASE_TOKEN })
    return
  }
  const maintenanceWorkerCwd = agentMaintenanceWorkerCwd(rawArgs)
  if (maintenanceWorkerCwd !== null) {
    await runAgentMaintenanceWorker(maintenanceWorkerCwd)
    return
  }
  const { command, state } = createCommandTree()
  const cleanOutput = rawArgs.some((arg) => ['-h', '--help', '-v', '--version', '--help-tree', '--help-docs'].includes(arg) || arg.startsWith('--help-tree-depth'))
  if (!cleanOutput) {
    const notice = await readCachedUpdateNotice(rawArgs.includes('--verbose'))
    if (notice) {
      if (rawArgs.length === 0) write(`${notice}\n`)
      else process.stderr.write(`${notice}\n`)
    }
  }
  await spawnUpdateWorker()
  try {
    await runCittyCommand(command, { rawArgs })
  } catch (error) {
    if (error instanceof CliHandled) return
    if (error instanceof CliUsageError) throw error
    if (error instanceof Error && error.name === 'CLIError') throw new CliUsageError(error.message, await renderUsage(state.command))
    throw error
  } finally {
    if (shouldScheduleAgentMaintenance(rawArgs)) {
      spawnAgentMaintenanceWorker(resolvePath(state.args.cwd ?? process.cwd()))
    }
  }
}

function shouldScheduleAgentMaintenance(rawArgs: string[]): boolean {
  return !['agent', 'uninstall'].includes(rawArgs[0] ?? '')
}

async function runCliWithErrorHandling(rawArgs?: string[]): Promise<void> {
  try {
    await runCli(rawArgs)
  } catch (error) {
    if (error instanceof CliUsageError) {
      process.stderr.write(`${error.usage}\n\nerror: ${error.message}\n`)
      process.exitCode = 2
    } else if (error instanceof RunXError) {
      process.stderr.write(`error: ${error.message}\n`)
      process.exitCode = error.exitCode
    } else {
      process.stderr.write(`error: ${error instanceof Error ? error.message : 'Unexpected failure.'}\n`)
      process.exitCode = 1
    }
  }
}

async function load(optionsValue: CliOptions): Promise<Awaited<ReturnType<typeof readManifest>>> {
  const loaded = await readManifest(optionsValue.cwd, optionsValue.config)
  const message = `configuration file loaded: ${loaded.path}\n`
  if (optionsValue.format === 'json') process.stderr.write(message)
  else write(message)
  return loaded
}

async function listCommands(optionsValue: CliOptions): Promise<void> {
  const { manifest, path } = await load(optionsValue)
  write(optionsValue.format === 'json' ? renderJson({ manifestPath: path, manifest }) : renderList(manifest, path))
}

async function describeCommand(selector: string, optionsValue: CliOptions): Promise<void> {
  const { manifest, path } = await load(optionsValue)
  const command = resolveCommand(manifest, path, selector)
  write(optionsValue.format === 'json' ? renderJson(command) : renderDescription(command))
}

async function checkManifest(optionsValue: CliOptions): Promise<void> {
  const { manifest, path } = await load(optionsValue)
  const result = { valid: true, manifestPath: path, commandCount: manifest.commands.length, groups: Object.keys(manifest.groups) }
  write(optionsValue.format === 'json' ? renderJson(result) : `valid: true\nmanifest: ${path}\ncommands: ${result.commandCount}\n`)
}

async function runSelected(selector: string, optionsValue: CliOptions, args: GlobalArgs): Promise<void> {
  const { manifest, path } = await load(optionsValue)
  const command = resolveCommand(manifest, path, selector)
  if (args.dryRun) {
    write(optionsValue.format === 'json' ? renderJson({ dryRun: true, command }) : renderExecutionPlan(command))
    return
  }
  if (command.confirm === 'always' && !args.yes) throw new RunXError(`Command ${command.uid} requires confirmation. Rerun with --yes only after authorization.`, 2)
  if (optionsValue.format === 'text') write(`Running ${command.uid} (${command.selector})\n`)
  process.exitCode = await runCommand(command)
}

async function agentMutation(label: string, action: (scope: AgentScope, cwd: string) => Promise<string[]>, args: GlobalArgs): Promise<void> {
  try {
    const paths = await action(args.local ? 'local' : 'global', resolvePath(args.cwd ?? process.cwd()))
    writeFormatted({ [label]: paths }, args)
  } catch (error) {
    if (error instanceof RunXError) throw error
    throw new RunXError(`Agent skill mutation failed: ${error instanceof Error ? error.message : String(error)}`, 5)
  }
}

async function instructionMutation(label: string, action: (cwd: string) => Promise<string[]>, args: GlobalArgs): Promise<void> {
  try {
    const paths = await action(resolvePath(args.cwd ?? process.cwd()))
    writeFormatted({ [label]: paths }, args)
  } catch (error) {
    if (error instanceof RunXError) throw error
    throw new RunXError(`Agent instruction mutation failed: ${error instanceof Error ? error.message : String(error)}`, 5)
  }
}

async function runUpgradeList(args: GlobalArgs): Promise<void> {
  const requestedArch = args.arch ? decode<UpgradeArch>(archSchema, args.arch, '--arch') : process.arch
  const platform = resolveUpgradePlatform(process.platform, requestedArch)
  const variant = decode<UpgradeVariant>(variantSchema, args.variant ?? 'baseline', '--variant')
  const page = args.page ? positiveInteger(args.page, '--page') : undefined
  const perPage = args.perPage ? positiveInteger(args.perPage, '--per-page') : undefined
  const catalog = await fetchReleaseCatalog({ ...platform, variant, currentVersion: readVersion() })
  const paged = paginateReleaseCatalog(catalog, page, perPage)
  write(options(args).format === 'json' ? renderJson(paged) : renderReleaseCatalog(paged))
}

async function runUpgrade(optionsValue: CliOptions, args: GlobalArgs): Promise<void> {
  if (args.arch) decode(archSchema, args.arch, '--arch')
  if (args.variant) decode(variantSchema, args.variant, '--variant')
  const text = optionsValue.format === 'text'
  if (text) write(renderUpgradeHeading())
  const result = await upgradeSelf({
    dryRun: Boolean(args.dryRun),
    requestedVersion: typeof args.version === 'string' ? args.version : undefined,
    arch: args.arch as UpgradeArch | undefined,
    variant: args.variant as UpgradeVariant | undefined,
    onPlan: text ? (plan) => write(renderUpgradePlan(plan)) : undefined,
    onEvent: text ? (event) => write(renderUpgradeEvent(event)) : undefined,
  })
  write(text ? renderUpgradeResult(result) : renderJson(result))
  if (result.outcome === 'failed' || result.outcome === 'rolled-back') process.exitCode = result.error?.code.includes('download') ? 4 : 5
}

function options(args: GlobalArgs): CliOptions {
  const format = decode(formatSchema, args.format ?? 'text', '--format')
  return { cwd: resolvePath(args.cwd ?? process.cwd()), config: args.config, format: format as OutputFormat, verbose: Boolean(args.verbose) }
}

function decode<T>(schema: Parameters<typeof Value.Decode>[0], value: unknown, flag: string): T {
  try {
    return Value.Decode(schema, value) as T
  } catch {
    throw new RunXError(`Invalid ${flag} value: ${String(value)}`, 2)
  }
}

function positiveInteger(value: string, flag: string): number {
  const number = Number(value)
  try {
    return Value.Decode(positiveIntegerSchema, number)
  } catch {
    throw new RunXError(`${flag} must be a positive integer.`, 2)
  }
}

async function usageError(state: CliState, message: string): Promise<CliUsageError> {
  return new CliUsageError(message, await renderUsage(state.command))
}

function writeFormatted(value: unknown, args: GlobalArgs): void {
  write(options(args).format === 'json' ? renderJson(value) : `${typeof value === 'string' ? value : JSON.stringify(value, null, 2)}\n`)
}

function handled(value: string): never {
  write(value.endsWith('\n') ? value : `${value}\n`)
  throw new CliHandled()
}

function write(value: string): void {
  process.stdout.write(value)
}
