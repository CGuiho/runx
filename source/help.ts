/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import packageJson from '../package.json' with { type: 'json' }

import type { CommandDef } from 'citty'

export {
  readVersion,
  renderHelpDocs,
  renderHelpTree,
}

function readVersion(): string {
  return typeof packageJson.version === 'string' ? packageJson.version : '0.0.0'
}

function renderHelpTree(command: CommandDef<any>, depth = Number.POSITIVE_INFINITY): string {
  const lines = ['COMMAND TREE', '', commandMeta(command).name ?? 'runx']
  appendChildren(lines, command, '', depth)
  return `${lines.join('\n')}\n`
}

function appendChildren(lines: string[], command: CommandDef<any>, prefix: string, depth: number): void {
  if (depth <= 0) return
  const subCommands = Object.entries(command.subCommands ?? {}).filter(([name, child]) => !name.startsWith('_') && !commandMeta(child).hidden)
  const flags = Object.entries(command.args ?? {}) as Array<[string, { type: string, valueHint?: string, description?: string }]>
  const visibleFlags = flags.filter(([, value]) => value.type !== 'positional')
  const nodes = [
    ...subCommands.map(([name, child]) => ({ label: name, description: commandMeta(child).description ?? '', child })),
    ...visibleFlags.map(([name, value]) => ({
      label: `--${name}${value.type === 'string' ? ` <${value.valueHint ?? 'value'}>` : ''}`,
      description: value.description ?? '',
      child: null,
    })),
  ]
  const width = Math.max(0, ...nodes.map((node) => node.label.length))
  nodes.forEach((node, index) => {
    const last = index === nodes.length - 1
    lines.push(`${prefix}${last ? '└── ' : '├── '}${node.label.padEnd(width)}${node.description ? `  ${node.description}` : ''}`)
    if (node.child) appendChildren(lines, node.child, `${prefix}${last ? '    ' : '│   '}`, depth - 1)
  })
}

function renderHelpDocs(command: CommandDef<any>): string {
  const meta = commandMeta(command)
  const name = meta.name ?? 'runx'
  const description = meta.description ?? ''
  const args = Object.entries(command.args ?? {}) as Array<[string, { type: string, valueHint?: string, description?: string }]>
  const positionals = args.filter(([, value]) => value.type === 'positional')
  const flags = args.filter(([, value]) => value.type !== 'positional')
  const children = Object.entries(command.subCommands ?? {}).filter(([key, value]) => !key.startsWith('_') && !commandMeta(value).hidden)
  const syntax = [
    name,
    ...positionals.map(([key]) => `<${key}>`),
    children.length ? '<command>' : '',
    flags.length ? '[options]' : '',
  ].filter(Boolean).join(' ')
  const lines = [`# ${name}`, '', description, '', '## Syntax', '', '```text', syntax, '```']
  if (positionals.length) {
    lines.push('', '## Positionals', '')
    for (const [key, value] of positionals) lines.push(`- \`${key}\` — ${value.description ?? ''}`)
  }
  if (flags.length) {
    lines.push('', '## Flags', '')
    for (const [key, value] of flags) lines.push(`- \`--${key}${value.type === 'string' ? ` <${value.valueHint ?? 'value'}>` : ''}\` — ${value.description ?? ''}`)
  }
  if (children.length) {
    lines.push('', '## Subcommands', '')
    for (const [key, value] of children) lines.push(`- \`${key}\` — ${commandMeta(value).description ?? ''}`)
  }
  lines.push('', '## Examples', '', '```text', `${name} --help`, `${name} --help-tree`, `${name} --help-docs`, '```', '')
  return lines.join('\n')
}

function commandMeta(command: CommandDef<any>): { name?: string, description?: string, hidden?: boolean } {
  return (command.meta ?? {}) as { name?: string, description?: string, hidden?: boolean }
}
