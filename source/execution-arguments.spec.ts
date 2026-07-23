/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { partitionRunInvocation } from './execution-arguments.js'

describe('RunX run argument ownership', () => {
  test('leaves non-run commands unchanged', () => {
    expect(partitionRunInvocation(['list', '--format', 'json'])).toEqual({
      routerArguments: ['list', '--format', 'json'],
      forwardedArguments: [],
    })
  })

  test('keeps RunX options before the selector and forwards everything after it', () => {
    expect(partitionRunInvocation(['run', '--dry-run', '--cwd', 'project', 'cli-ts', '-v', 'build', '--watch'])).toEqual({
      routerArguments: ['run', '--dry-run', '--cwd', 'project', 'cli-ts'],
      forwardedArguments: ['-v', 'build', '--watch'],
    })
  })

  test('removes one immediate delimiter and preserves later delimiters', () => {
    expect(partitionRunInvocation(['run', 'cli-ts', '--', '--dry-run', '--', 'tail'])).toEqual({
      routerArguments: ['run', 'cli-ts'],
      forwardedArguments: ['--dry-run', '--', 'tail'],
    })
  })

  test('preserves empty, Unicode, and hostile-looking values without interpretation', () => {
    const values = ['', 'Olá 世界', '; exit 99', '&& echo injected', '$(whoami)', '%PATH%', '!ERRORLEVEL!']
    expect(partitionRunInvocation(['run', 'literal', ...values]).forwardedArguments).toEqual(values)
  })
})
