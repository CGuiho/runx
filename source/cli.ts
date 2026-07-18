/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { resolve } from 'node:path'
import { defineCommand, renderUsage, runCommand as runCittyCommand } from 'citty'
import { installAgentInstructions, installAgentSkill } from './agents.js'
import { RunXError } from './errors.js'
import { runCommand } from './executor.js'
import { readVersion, showHelpDocs, showHelpTree, showHome } from './help.js'
import { initializeRunXManifest } from './init.js'
import { readManifest, resolveCommand } from './manifest.js'
import { renderDescription, renderExecutionPlan, renderJson, renderList } from './render.js'
import { checkForLatestVersion, listAvailableVersions, uninstallSelf, upgradeSelf } from './self-management.js'
import { renderReleaseCatalog, renderUpgradeEvent, renderUpgradeHeading, renderUpgradePlan, renderUpgradeResult } from './upgrade-reporting.js'

import type { CommandContext, CommandDef, SubCommandsDef } from 'citty'
import type { AgentScope, AgentTool, CliOptions, OutputFormat } from './types.js'

export {
  runCli,
  runCliWithErrorHandling,
  runxCommand,
}

type GlobalArgs = {
  _: string[]
  cwd?: string
  file?: string
  format?: string
  verbose?: boolean
  dryRun?: boolean
  yes?: boolean
  tool?: string
  helpTree?: boolean
  helpDocs?: boolean
  help?: boolean
  h?: boolean
  version?: boolean
  v?: boolean
  [key: string]: string | number | boolean | string[] | undefined
}

type CliState = {
  globalArgs: GlobalArgs
  usageCommand: CommandDef<any>
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

const commonArgs = {
  cwd: { type: 'string', description: 'Resolve the manifest from this directory.', valueHint: 'path' },
  file: { type: 'string', description: 'Use this runx.yaml file.', valueHint: 'path' },
  format: { type: 'string', description: 'Select text or JSON output.', valueHint: 'text|json', default: 'text' },
  verbose: { type: 'boolean', description: 'Enable additional diagnostics.' },
  'help-tree': { type: 'boolean', description: 'Show the complete command hierarchy.' },
  'help-docs': { type: 'boolean', description: 'Show manifest and agent guidance.' },
  help: { type: 'boolean', alias: 'h', description: 'Show command help.' },
  version: { type: 'boolean', alias: 'v', description: 'Show the RunX version.' },
} as const

const runArgs = {
  ...commonArgs,
  'dry-run': { type: 'boolean', description: 'Show the action without executing it.' },
  yes: { type: 'boolean', description: 'Approve a confirmation-gated command.' },
} as const

const agentInstallArgs = {
  ...commonArgs,
  tool: { type: 'string', description: 'Select agents, claude, or all.', valueHint: 'agents|claude|all' },
} as const

const selfManagementArgs = {
  ...commonArgs,
  'dry-run': { type: 'boolean', description: 'Show the action without changing the executable.' },
} as const

const rootArgs = {
  ...runArgs,
  tool: agentInstallArgs.tool,
} as const

const knownArgumentKeys = new Set([
  '_',
  'cwd',
  'file',
  'format',
  'verbose',
  'dryRun',
  'dry-run',
  'yes',
  'tool',
  'helpTree',
  'help-tree',
  'helpDocs',
  'help-docs',
  'help',
  'h',
  'version',
  'v',
])

const createCommandTree = (): { command: CommandDef<any>, state: CliState } => {
  const state: CliState = { globalArgs: { _: [] }, usageCommand: {} }

  const listCommand = defineCommand({
    meta: { name: 'runx list', description: 'List commands in a RunX manifest.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => listCommands(resolveOptions(state.globalArgs)),
  })

  const describeCommandDefinition = defineCommand({
    meta: { name: 'runx describe', description: 'Describe one manifest command.' },
    args: {
      selector: { type: 'positional', description: 'Command UID, group/id, index, or unambiguous ID.' },
      ...commonArgs,
    },
    setup: withCommandHelp(state),
    run: async ({ args }) => {
      if (!args.selector) throw await missingSelector(state)
      await describeCommand(args.selector, resolveOptions(state.globalArgs))
    },
  })

  const runCommandDefinition = createRunCommand(state)

  const checkCommand = defineCommand({
    meta: { name: 'runx check', description: 'Validate a RunX manifest without executing it.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => checkManifest(resolveOptions(state.globalArgs)),
  })

  const initCommand = defineCommand({
    meta: { name: 'runx init', description: 'Interactively create an empty RunX manifest.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => initializeProject(state.globalArgs),
  })

  const agentsInstallCommand = defineCommand({
    meta: { name: 'runx agents install', description: 'Install the bundled RunX skill.' },
    args: {
      scope: { type: 'positional', description: 'Install locally or globally.', default: 'local', valueHint: 'local|global' },
      ...agentInstallArgs,
    },
    setup: withCommandHelp(state),
    run: async ({ args }) => runAgentsInstall(args.scope, resolveOptions(state.globalArgs), state.globalArgs),
  })

  const agentsInstructionsCommand = defineCommand({
    meta: { name: 'runx agents instructions', description: 'Install or refresh managed RunX agent instructions.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => runAgentsInstructions(resolveOptions(state.globalArgs)),
  })

  const agentsCommand = defineCommand({
    meta: { name: 'runx agents', description: 'Manage the bundled RunX agent integration.' },
    args: agentInstallArgs,
    subCommands: {
      install: agentsInstallCommand,
      instructions: agentsInstructionsCommand,
    },
    setup: withCommandHelp(state),
  })

  const upgradeApplyCommand = defineCommand({
    meta: { name: 'runx upgrade', description: 'Upgrade a native RunX executable.', hidden: true },
    args: selfManagementArgs,
    setup: withCommandHelp(state),
    run: async () => runUpgrade(resolveOptions(state.globalArgs), state.globalArgs),
  })

  const upgradeCheckCommand = defineCommand({
    meta: { name: 'runx upgrade check', description: 'Check whether a newer RunX release is available.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => runUpgradeCheck(resolveOptions(state.globalArgs)),
  })

  const upgradeListCommand = defineCommand({
    meta: { name: 'runx upgrade list', description: 'List every published RunX release, newest first.' },
    args: commonArgs,
    setup: withCommandHelp(state),
    run: async () => runUpgradeList(resolveOptions(state.globalArgs)),
  })

  const upgradeCommand = defineCommand({
    meta: { name: 'runx upgrade', description: 'Inspect or upgrade a native RunX executable.' },
    args: selfManagementArgs,
    default: '_apply',
    subCommands: {
      _apply: upgradeApplyCommand,
      check: upgradeCheckCommand,
      list: upgradeListCommand,
    },
    setup: withCommandHelp(state),
  })

  const uninstallCommand = defineCommand({
    meta: { name: 'runx uninstall', description: 'Uninstall a native RunX executable.' },
    args: selfManagementArgs,
    setup: withCommandHelp(state),
    run: async () => runUninstall(resolveOptions(state.globalArgs), state.globalArgs),
  })

  const homeCommand = defineCommand({
    meta: { name: 'runx', description: 'Show the RunX home page.', hidden: true },
    args: rootArgs,
    setup: withCommandHelp(state),
    run: () => write(showHome()),
  })

  const staticSubCommands = Object.assign(Object.create(null) as SubCommandsDef, {
    _home: homeCommand,
    list: listCommand,
    describe: describeCommandDefinition,
    run: runCommandDefinition,
    check: checkCommand,
    init: initCommand,
    agents: agentsCommand,
    upgrade: upgradeCommand,
    uninstall: uninstallCommand,
  })
  const usageCommands = new Map<string, CommandDef<any>>([
    ['list', listCommand],
    ['describe', describeCommandDefinition],
    ['run', runCommandDefinition],
    ['r', runCommandDefinition],
    ['check', checkCommand],
    ['init', initCommand],
    ['agents', agentsCommand],
    ['agents install', agentsInstallCommand],
    ['agents instructions', agentsInstructionsCommand],
    ['upgrade', upgradeCommand],
    ['upgrade check', upgradeCheckCommand],
    ['upgrade list', upgradeListCommand],
    ['uninstall', uninstallCommand],
  ])
  const subCommands = createSelectorCompatibleCommands(staticSubCommands, state)

  const command = defineCommand({
    meta: {
      name: 'runx',
      version: readVersion(),
      description: 'A documented, local command catalog for runx.yaml manifests.',
    },
    args: rootArgs,
    default: '_home',
    subCommands,
    setup: async (context) => {
      state.globalArgs = context.args as GlobalArgs
      state.usageCommand = resolveUsageCommand(context.args._, usageCommands, state)
      await validateArguments(context, state)
      if (state.globalArgs.version || state.globalArgs.v) return handleOutput(`${readVersion()}\n`)
      if (state.globalArgs.helpTree) return handleOutput(`${showHelpTree()}\n`)
      if (state.globalArgs.helpDocs) return handleOutput(showHelpDocs())
      if (state.globalArgs.help || state.globalArgs.h) return handleUsage(state.usageCommand)
    },
  })

  state.usageCommand = command
  return { command, state }
}

const { command: runxCommand } = createCommandTree()

async function runCli(rawArgs: string[] = process.argv.slice(2)): Promise<void> {
  const { command, state } = createCommandTree()
  try {
    await runCittyCommand(command, { rawArgs })
  } catch (error) {
    if (error instanceof CliHandled) return
    if (error instanceof CliUsageError) throw error
    if (isCittyUsageError(error)) {
      throw new CliUsageError(error.message, await renderUsage(state.usageCommand))
    }
    throw error
  }
}

async function runCliWithErrorHandling(rawArgs?: string[]): Promise<void> {
  try {
    await runCli(rawArgs)
  } catch (error) {
    if (error instanceof CliUsageError) {
      process.stderr.write(`${error.usage}\n\nerror: ${error.message}\n`)
      process.exitCode = 1
      return
    }
    if (error instanceof RunXError) {
      process.stderr.write(`error: ${error.message}\n`)
      process.exitCode = error.exitCode
      return
    }
    process.stderr.write(`error: ${error instanceof Error ? error.message : 'Unexpected failure.'}\n`)
    process.exitCode = 1
  }
}

function createRunCommand(state: CliState, selectorDefault?: string): CommandDef<any> {
  return defineCommand({
    meta: {
      name: selectorDefault ? 'runx <selector>' : 'runx run',
      alias: selectorDefault ? undefined : 'r',
      description: selectorDefault ? 'Run a selector using the root shorthand.' : 'Run one manifest command.',
    },
    args: selectorDefault
      ? runArgs
      : {
          selector: { type: 'positional', description: 'Command UID, group/id, index, or unambiguous ID.' },
          ...runArgs,
        },
    setup: withCommandHelp(state),
    run: async ({ args }) => {
      const selector = selectorDefault ?? (typeof args.selector === 'string' ? args.selector : undefined)
      if (!selector) throw await missingSelector(state)
      await runSelectedCommand(selector, resolveOptions(state.globalArgs), state.globalArgs)
    },
  })
}

function createSelectorCompatibleCommands(commands: SubCommandsDef, state: CliState): SubCommandsDef {
  const aliases = new Set(['r'])
  return new Proxy(commands, {
    has: (target, property) => {
      if (typeof property !== 'string') return Reflect.has(target, property)
      if (aliases.has(property)) return false
      return true
    },
    get: (target, property, receiver) => {
      if (typeof property !== 'string' || Reflect.has(target, property)) return Reflect.get(target, property, receiver)
      return createRunCommand(state, property)
    },
  })
}

function resolveUsageCommand(positionals: string[], commands: Map<string, CommandDef<any>>, state: CliState): CommandDef<any> {
  const name = positionals[0]
  if (!name) return state.usageCommand
  if (!commands.has(name)) return createRunCommand(state, name)
  for (let length = positionals.length; length > 0; length -= 1) {
    const command = commands.get(positionals.slice(0, length).join(' '))
    if (command) return command
  }
  return commands.get(name)!
}

function withCommandHelp(state: CliState): (context: CommandContext<any>) => void {
  return (context) => {
    state.usageCommand = context.cmd
  }
}

async function missingSelector(state: CliState): Promise<CliUsageError> {
  return new CliUsageError('Missing required positional argument: SELECTOR', await renderUsage(state.usageCommand))
}

async function validateArguments(context: CommandContext<any>, state: CliState): Promise<void> {
  const unknown = Object.keys(context.args).find((key) => !knownArgumentKeys.has(key))
  if (!unknown) return
  throw new CliUsageError(`Unknown option --${unknown}`, await renderUsage(state.usageCommand))
}

async function handleUsage(command: CommandDef<any>): Promise<never> {
  return handleOutput(`${await renderUsage(command)}\n`)
}

function handleOutput(value: string): never {
  write(value)
  throw new CliHandled()
}

async function listCommands(options: CliOptions): Promise<void> {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  write(options.format === 'json' ? renderJson({ manifestPath: path, manifest }) : renderList(manifest, path))
}

async function describeCommand(selector: string, options: CliOptions): Promise<void> {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const command = resolveCommand(manifest, path, selector)
  write(options.format === 'json' ? renderJson(command) : renderDescription(command))
}

async function checkManifest(options: CliOptions): Promise<void> {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const result = { valid: true, manifestPath: path, commandCount: manifest.commands.length, groups: Object.keys(manifest.groups) }
  write(options.format === 'json' ? renderJson(result) : `valid: true\nmanifest: ${path}\ncommands: ${result.commandCount}\n`)
}

async function initializeProject(args: GlobalArgs): Promise<void> {
  if (args.file) {
    throw new RunXError('runx init does not support --file. It always creates runx.yaml in --cwd or the current directory.')
  }
  if (args.format === 'json') {
    throw new RunXError('runx init does not support --format json because initialization is an interactive terminal workflow.')
  }
  await initializeRunXManifest({ cwd: resolve(args.cwd ?? process.cwd()) })
}

async function runSelectedCommand(selector: string, options: CliOptions, args: GlobalArgs): Promise<void> {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const command = resolveCommand(manifest, path, selector)

  if (args.dryRun) {
    write(options.format === 'json' ? renderJson({ dryRun: true, command }) : renderExecutionPlan(command))
    return
  }

  if (command.confirm === 'always' && !args.yes) {
    throw new RunXError(`Command ${command.uid} requires confirmation. Review it with \`runx describe ${command.uid}\` and rerun with --yes when authorized.`)
  }

  if (options.format === 'text') write(`Running ${command.uid} (${command.selector})\n`)
  process.exitCode = await runCommand(command)
}

async function runAgentsInstall(scopeValue: string, options: CliOptions, args: GlobalArgs): Promise<void> {
  const scope = scopeValue as AgentScope
  if (scope !== 'local' && scope !== 'global') throw new RunXError('Agent scope must be local or global.')
  const tool = (args.tool ?? 'agents') as AgentTool
  if (!['agents', 'claude', 'all'].includes(tool)) throw new RunXError('Agent tool must be agents, claude, or all.')
  const installed = await installAgentSkill(scope, tool, options.cwd)
  write(options.format === 'json' ? renderJson({ installed }) : `${installed.map((path) => `installed: ${path}`).join('\n')}\n`)
}

async function runAgentsInstructions(options: CliOptions): Promise<void> {
  const path = await installAgentInstructions(options.cwd)
  write(options.format === 'json' ? renderJson({ instructions: path }) : `updated: ${path}\n`)
}

async function runUpgradeCheck(options: CliOptions): Promise<void> {
  const result = await checkForLatestVersion()
  write(options.format === 'json' ? renderJson(result) : `current: ${result.currentVersion}\nlatest: ${result.latestVersion}\nupdate_available: ${result.updateAvailable}\n`)
}

async function runUpgradeList(options: CliOptions): Promise<void> {
  const catalog = await listAvailableVersions()
  write(options.format === 'json' ? renderJson(catalog) : renderReleaseCatalog(catalog))
}

async function runUpgrade(options: CliOptions, args: GlobalArgs): Promise<void> {
  const text = options.format === 'text'
  if (text) write(renderUpgradeHeading())
  const result = await upgradeSelf({
    dryRun: Boolean(args.dryRun),
    onPlan: text ? (plan) => write(renderUpgradePlan(plan)) : undefined,
    onEvent: text ? (event) => write(renderUpgradeEvent(event)) : undefined,
  })
  if (!text) write(renderJson(result))
  else write(renderUpgradeResult(result))
  if (result.outcome === 'failed' || result.outcome === 'rolled-back') {
    process.stderr.write(`error: ${result.error?.message ?? 'RunX upgrade failed.'}\n`)
    process.exitCode = 1
  }
}

async function runUninstall(options: CliOptions, args: GlobalArgs): Promise<void> {
  const result = await uninstallSelf(Boolean(args.dryRun))
  write(options.format === 'json' ? renderJson(result) : `path: ${result.executablePath}\n${result.dryRun ? 'dry_run: true\n' : result.scheduled ? 'scheduled: true\n' : 'uninstalled: true\n'}`)
}

function resolveOptions(args: GlobalArgs): CliOptions {
  const format = args.format ?? 'text'
  if (format !== 'text' && format !== 'json') throw new RunXError('Invalid --format value. Expected text or json.')
  return {
    cwd: resolve(args.cwd ?? process.cwd()),
    file: args.file,
    format: format as OutputFormat,
    verbose: Boolean(args.verbose),
  }
}

function isCittyUsageError(error: unknown): error is Error {
  return error instanceof Error && error.name === 'CLIError'
}

function write(value: string): void {
  process.stdout.write(value)
}
