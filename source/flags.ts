import type { ParsedArgs } from './types.js'

const booleanFlags = new Set(['help', 'helpTree', 'helpDocs', 'version', 'verbose', 'dryRun', 'yes'])

export const parseArgs = (args: string[]): ParsedArgs => {
  const positionals: string[] = []
  const flags: ParsedArgs['flags'] = {}

  for (let index = 0; index < args.length; index += 1) {
    const value = args[index]!

    if (!value.startsWith('--')) {
      positionals.push(value)
      continue
    }

    const [rawName, inlineValue] = value.slice(2).split('=', 2)
    const name = toCamelCase(rawName)
    if (inlineValue !== undefined) {
      flags[name] = inlineValue
      continue
    }

    const next = args[index + 1]
    if (!booleanFlags.has(name) && next && !next.startsWith('--')) {
      flags[name] = next
      index += 1
      continue
    }

    flags[name] = true
  }

  return { command: positionals.shift(), positionals, flags }
}

export const booleanFlag = (flags: ParsedArgs['flags'], name: string): boolean => flags[name] === true

export const stringFlag = (flags: ParsedArgs['flags'], name: string): string | undefined => {
  const value = flags[name]
  return typeof value === 'string' ? value : undefined
}

const toCamelCase = (value: string): string => value.replace(/-([a-z])/g, (_, letter: string) => letter.toUpperCase())
