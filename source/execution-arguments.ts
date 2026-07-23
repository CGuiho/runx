/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'

import type { Static } from '@sinclair/typebox'

export {
  partitionRunInvocation,
  runInvocationSchema,
}
export type {
  RunInvocation,
}

const runInvocationSchema = Type.Object({
  routerArguments: Type.Array(Type.String()),
  forwardedArguments: Type.Array(Type.String()),
}, { additionalProperties: false })
type RunInvocation = Static<typeof runInvocationSchema>

const valueOptions = new Set(['--cwd', '--config', '--format'])

function partitionRunInvocation(rawArguments: readonly string[]): RunInvocation {
  if (rawArguments[0] !== 'run') {
    return Value.Decode(runInvocationSchema, {
      routerArguments: [...rawArguments],
      forwardedArguments: [],
    })
  }

  let selectorIndex = -1
  for (let index = 1; index < rawArguments.length; index += 1) {
    const token = rawArguments[index] ?? ''
    if (token === '--') break
    if (valueOptions.has(token)) {
      index += 1
      continue
    }
    if (token.startsWith('--cwd=') || token.startsWith('--config=') || token.startsWith('--format=')) continue
    if (!token.startsWith('-')) {
      selectorIndex = index
      break
    }
  }

  if (selectorIndex === -1) {
    return Value.Decode(runInvocationSchema, {
      routerArguments: [...rawArguments],
      forwardedArguments: [],
    })
  }

  const boundary = selectorIndex + 1
  const forwarded = rawArguments.slice(boundary)
  const forwardedArguments = forwarded[0] === '--' ? forwarded.slice(1) : forwarded
  return Value.Decode(runInvocationSchema, {
    routerArguments: rawArguments.slice(0, boundary),
    forwardedArguments,
  })
}
