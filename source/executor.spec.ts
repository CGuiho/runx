/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { createShellExecution } from './executor.js'

import type { ResolvedCommand } from './types.js'

describe('RunX shell argument transport', () => {
  const base = (shell: ResolvedCommand['shell']): ResolvedCommand => ({
    uid: 'shell-test', id: 'shell-test', group: 'public', summary: 'Shell test.',
    description: 'Shell test.', command: 'child-command', shell,
    index: 0, selector: 'public/shell-test', manifestPath: 'runx.yaml', cwd: '.',
  })
  const values = ['-v', '', 'space value', 'Olá 世界', '; exit 99', '&& injected', '$(whoami)', '%PATH%', '!ERRORLEVEL!']

  test('uses positional parameters for POSIX shells', () => {
    for (const shell of ['bash', 'sh'] as const) {
      const execution = createShellExecution(base(shell), values)
      expect(execution.arguments.slice(-values.length)).toEqual(values)
      expect(execution.arguments[2]).toBe('child-command "$@"')
      for (const value of values.filter(Boolean)) expect(execution.arguments[2]).not.toContain(value)
    }
  })

  test('uses JSON-backed PowerShell splatting without interpolating values', () => {
    const execution = createShellExecution(base('powershell'), values)
    expect(JSON.parse(execution.environment.RUNX_FORWARDED_ARGUMENTS_JSON ?? '')).toEqual(values.map((value) => value === '' ? '""' : value))
    for (const value of values.filter(Boolean)) expect(execution.arguments.at(-1)).not.toContain(value)
  })

  test('uses delayed-expansion-disabled cmd environment references', () => {
    const execution = createShellExecution(base('cmd'), values)
    expect(execution.arguments).toContain('/v:off')
    expect(execution.script?.content).toContain('setlocal DisableDelayedExpansion')
    values.forEach((value, index) => {
      expect(execution.environment[`RUNX_FORWARDED_ARGUMENT_${index}`]).toBe(value)
      expect(execution.script?.content).toContain(`"%RUNX_FORWARDED_ARGUMENT_${index}%"`)
      if (value) expect(execution.script?.content).not.toContain(value)
    })
  })
})
