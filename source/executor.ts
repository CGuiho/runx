/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { joinPath } from './path-utils.js'
import { removePath, writeTextFile } from './storage.js'

import type { ResolvedCommand } from './types.js'

export {
  createShellExecution,
  runCommand,
}

type ShellExecution = {
  readonly arguments: string[]
  readonly environment: Record<string, string>
  readonly script?: {
    readonly content: string
    readonly path: string
  }
}

async function runCommand(command: ResolvedCommand, forwardedArguments: readonly string[] = []): Promise<number> {
  const execution = createShellExecution(command, forwardedArguments)
  if (execution.script) await writeTextFile(execution.script.path, execution.script.content)
  try {
    const child = Bun.spawn(execution.arguments, {
      cwd: command.cwd,
      env: { ...process.env, ...execution.environment },
      stdin: 'inherit',
      stdout: 'inherit',
      stderr: 'inherit',
    })
    return await child.exited
  } finally {
    if (execution.script) await removePath(execution.script.path).catch(() => undefined)
  }
}

function createShellExecution(command: ResolvedCommand, forwardedArguments: readonly string[] = []): ShellExecution {
  switch (command.shell ?? 'auto') {
    case 'bash': return positionalShellExecution('bash', command.command, forwardedArguments)
    case 'sh': return positionalShellExecution('sh', command.command, forwardedArguments)
    case 'powershell': return powerShellExecution(command.command, forwardedArguments)
    case 'cmd': return cmdExecution(command.command, forwardedArguments)
    case 'auto': return process.platform === 'win32'
      ? cmdExecution(command.command, forwardedArguments)
      : positionalShellExecution('sh', command.command, forwardedArguments)
  }
}

function positionalShellExecution(shell: 'bash' | 'sh', command: string, forwardedArguments: readonly string[]): ShellExecution {
  return {
    arguments: [shell, '-lc', `${command} "$@"`, 'runx-child', ...forwardedArguments],
    environment: {},
  }
}

function powerShellExecution(command: string, forwardedArguments: readonly string[]): ShellExecution {
  const environmentKey = 'RUNX_FORWARDED_ARGUMENTS_JSON'
  const powerShellArguments = forwardedArguments.map((argument) => argument === '' ? '""' : argument)
  const script = [
    `$runxForwarded = @(ConvertFrom-Json -InputObject $env:${environmentKey})`,
    `& { ${command} @runxForwarded }`,
    'exit $LASTEXITCODE',
  ].join('; ')
  return {
    arguments: [process.platform === 'win32' ? 'powershell.exe' : 'pwsh', '-NoProfile', '-NonInteractive', '-Command', script],
    environment: { [environmentKey]: JSON.stringify(powerShellArguments) },
  }
}

function cmdExecution(command: string, forwardedArguments: readonly string[]): ShellExecution {
  const environment: Record<string, string> = {}
  const references = forwardedArguments.map((value, index) => {
    const key = `RUNX_FORWARDED_ARGUMENT_${index}`
    environment[key] = value
    return `"%${key}%"`
  })
  const scriptPath = joinPath(Bun.env.TEMP ?? Bun.env.TMP ?? '.', `runx-command-${process.pid}-${crypto.randomUUID()}.cmd`)
  return {
    arguments: ['cmd.exe', '/d', '/v:off', '/s', '/c', scriptPath],
    environment,
    script: {
      path: scriptPath,
      content: `@echo off\r\nsetlocal DisableDelayedExpansion\r\n${[command, ...references].join(' ')}\r\nexit /b %ERRORLEVEL%\r\n`,
    },
  }
}
