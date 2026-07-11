import { resolve } from 'node:path'
import { installAgentInstructions, installAgentSkill } from './agents.js'
import { RunXError } from './errors.js'
import { runCommand } from './executor.js'
import { booleanFlag, parseArgs, stringFlag } from './flags.js'
import { readVersion, showCommandHelp, showHelpDocs, showHelpTree, showHome } from './help.js'
import { readManifest, resolveCommand } from './manifest.js'
import { renderDescription, renderExecutionPlan, renderJson, renderList } from './render.js'
import { checkForLatestVersion, listAvailableVersions, uninstallSelf, upgradeSelf } from './self-management.js'
import type { AgentScope, AgentTool, CliOptions, OutputFormat } from './types.js'

const builtins = new Set(['list', 'describe', 'run', 'r', 'check', 'agents', 'upgrade', 'uninstall'])

export const runCli = async (rawArgs: string[] = process.argv.slice(2)): Promise<void> => {
  const parsed = parseArgs(rawArgs)
  const options = resolveOptions(parsed.flags)

  if (booleanFlag(parsed.flags, 'version')) return write(`${readVersion()}\n`)
  if (booleanFlag(parsed.flags, 'helpTree')) return write(`${showHelpTree()}\n`)
  if (booleanFlag(parsed.flags, 'helpDocs')) return write(showHelpDocs())
  if (booleanFlag(parsed.flags, 'help')) return write(parsed.command && builtins.has(parsed.command) ? showCommandHelp(parsed.command) : showHome())
  if (!parsed.command) return write(showHome())

  if (!builtins.has(parsed.command)) {
    await runSelectedCommand(parsed.command, parsed.positionals, options, parsed.flags)
    return
  }

  switch (parsed.command) {
    case 'list': await listCommands(options); return
    case 'describe': await describeCommand(parsed.positionals[0], options); return
    case 'run':
    case 'r': await runSelectedCommand(parsed.positionals[0], parsed.positionals.slice(1), options, parsed.flags); return
    case 'check': await checkManifest(options); return
    case 'agents': await runAgents(parsed.positionals, options, parsed.flags); return
    case 'upgrade': await runUpgrade(parsed.positionals, options, parsed.flags); return
    case 'uninstall': await runUninstall(options, parsed.flags); return
  }
}

export const runCliWithErrorHandling = async (rawArgs?: string[]): Promise<void> => {
  try {
    await runCli(rawArgs)
  } catch (error) {
    if (error instanceof RunXError) {
      process.stderr.write(`error: ${error.message}\n`)
      process.exitCode = error.exitCode
      return
    }
    process.stderr.write(`error: ${error instanceof Error ? error.message : 'Unexpected failure.'}\n`)
    process.exitCode = 1
  }
}

const listCommands = async (options: CliOptions): Promise<void> => {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  write(options.format === 'json' ? renderJson({ manifestPath: path, manifest }) : renderList(manifest, path))
}

const describeCommand = async (selector: string | undefined, options: CliOptions): Promise<void> => {
  if (!selector) throw new RunXError('Missing selector. Usage: runx describe <selector>')
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const command = resolveCommand(manifest, path, selector)
  write(options.format === 'json' ? renderJson(command) : renderDescription(command))
}

const checkManifest = async (options: CliOptions): Promise<void> => {
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const result = { valid: true, manifestPath: path, commandCount: manifest.commands.length, groups: Object.keys(manifest.groups) }
  write(options.format === 'json' ? renderJson(result) : `valid: true\nmanifest: ${path}\ncommands: ${result.commandCount}\n`)
}

const runSelectedCommand = async (selector: string | undefined, _: string[], options: CliOptions, flags: Record<string, boolean | string | string[]>): Promise<void> => {
  if (!selector) throw new RunXError('Missing selector. Usage: runx run <selector>')
  const { manifest, path } = await readManifest(options.cwd, options.file)
  const command = resolveCommand(manifest, path, selector)

  if (booleanFlag(flags, 'dryRun')) {
    write(options.format === 'json' ? renderJson({ dryRun: true, command }) : renderExecutionPlan(command))
    return
  }

  if (command.confirm === 'always' && !booleanFlag(flags, 'yes')) {
    throw new RunXError(`Command ${command.uid} requires confirmation. Review it with \`runx describe ${command.uid}\` and rerun with --yes when authorized.`)
  }

  if (options.format === 'text') write(`Running ${command.uid} (${command.selector})\n`)
  const exitCode = await runCommand(command)
  process.exitCode = exitCode
}

const runAgents = async (positionals: string[], options: CliOptions, flags: Record<string, boolean | string | string[]>): Promise<void> => {
  const action = positionals[0]
  if (action === 'install') {
    const scope = (positionals[1] ?? 'local') as AgentScope
    if (scope !== 'local' && scope !== 'global') throw new RunXError('Agent scope must be local or global.')
    const tool = (stringFlag(flags, 'tool') ?? 'agents') as AgentTool
    if (!['agents', 'claude', 'all'].includes(tool)) throw new RunXError('Agent tool must be agents, claude, or all.')
    const installed = await installAgentSkill(scope, tool, options.cwd)
    write(options.format === 'json' ? renderJson({ installed }) : `${installed.map((path) => `installed: ${path}`).join('\n')}\n`)
    return
  }
  if (action === 'instructions') {
    const path = await installAgentInstructions(options.cwd)
    write(options.format === 'json' ? renderJson({ instructions: path }) : `updated: ${path}\n`)
    return
  }
  throw new RunXError('Usage: runx agents install <local|global> [--tool agents|claude|all] | runx agents instructions')
}

const runUpgrade = async (positionals: string[], options: CliOptions, flags: Record<string, boolean | string | string[]>): Promise<void> => {
  if (positionals[0] === 'check') {
    const result = await checkForLatestVersion()
    write(options.format === 'json' ? renderJson(result) : `current: ${result.currentVersion}\nlatest: ${result.latestVersion}\nupdate_available: ${result.updateAvailable}\n`)
    return
  }
  if (positionals[0] === 'list') {
    const versions = await listAvailableVersions()
    write(options.format === 'json' ? renderJson({ versions }) : `${versions.join('\n')}\n`)
    return
  }
  const result = await upgradeSelf(booleanFlag(flags, 'dryRun'))
  write(options.format === 'json' ? renderJson(result) : `current: ${result.currentVersion}\ntarget: ${result.latestVersion}\npath: ${result.executablePath}\n${result.scheduled ? 'scheduled: true\n' : ''}`)
}

const runUninstall = async (options: CliOptions, flags: Record<string, boolean | string | string[]>): Promise<void> => {
  const result = await uninstallSelf(booleanFlag(flags, 'dryRun'))
  write(options.format === 'json' ? renderJson(result) : `path: ${result.executablePath}\n${result.dryRun ? 'dry_run: true\n' : result.scheduled ? 'scheduled: true\n' : 'uninstalled: true\n'}`)
}

const resolveOptions = (flags: Record<string, boolean | string | string[]>): CliOptions => {
  const format = stringFlag(flags, 'format') ?? 'text'
  if (format !== 'text' && format !== 'json') throw new RunXError('Invalid --format value. Expected text or json.')
  return { cwd: resolve(stringFlag(flags, 'cwd') ?? process.cwd()), file: stringFlag(flags, 'file'), format: format as OutputFormat, verbose: booleanFlag(flags, 'verbose') }
}

const write = (value: string): void => {
  process.stdout.write(value)
}
