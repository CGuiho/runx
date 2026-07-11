import { existsSync } from 'node:fs'
import type { ResolvedCommand } from './types.js'
import { RunXError } from './errors.js'

export const runCommand = async (command: ResolvedCommand): Promise<number> => {
  if (!existsSync(command.cwd)) throw new RunXError(`Command working directory does not exist: ${command.cwd}`)

  const process = Bun.spawn(shellArguments(command), {
    cwd: command.cwd,
    stdin: 'inherit',
    stdout: 'inherit',
    stderr: 'inherit',
  })

  return process.exited
}

const shellArguments = (command: ResolvedCommand): string[] => {
  switch (command.shell ?? 'auto') {
    case 'bash': return ['bash', '-lc', command.command]
    case 'sh': return ['sh', '-lc', command.command]
    case 'powershell': return [process.platform === 'win32' ? 'powershell.exe' : 'pwsh', '-NoProfile', '-NonInteractive', '-Command', command.command]
    case 'cmd': return ['cmd.exe', '/d', '/s', '/c', command.command]
    case 'auto': return process.platform === 'win32'
      ? ['cmd.exe', '/d', '/s', '/c', command.command]
      : ['sh', '-lc', command.command]
  }
}
