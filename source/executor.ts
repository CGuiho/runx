/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import type { ResolvedCommand } from './types.js'

export {
  runCommand,
}

async function runCommand(command: ResolvedCommand): Promise<number> {
  const child = Bun.spawn(shellArguments(command), {
    cwd: command.cwd,
    stdin: 'inherit',
    stdout: 'inherit',
    stderr: 'inherit',
  })
  return child.exited
}

function shellArguments(command: ResolvedCommand): string[] {
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
