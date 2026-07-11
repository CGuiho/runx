import type { ResolvedCommand, RunXManifest } from './types.js'

export const renderJson = (value: unknown): string => JSON.stringify(value, null, 2) + '\n'

export const renderList = (manifest: RunXManifest, manifestPath: string): string => {
  const rows = manifest.commands.map((command, index) => ({
    index: index + 1,
    uid: command.uid,
    selector: `${command.group}/${command.id}`,
    summary: command.summary,
    confirm: command.confirm ?? 'never',
  }))
  const width = (key: keyof (typeof rows)[number], label: string): number => Math.max(label.length, ...rows.map((row) => String(row[key]).length))
  const indexWidth = width('index', 'IDX')
  const uidWidth = width('uid', 'UID')
  const selectorWidth = width('selector', 'SELECTOR')
  const lines = [`RunX commands from ${manifestPath}`, '', `${'IDX'.padEnd(indexWidth)}  ${'UID'.padEnd(uidWidth)}  ${'SELECTOR'.padEnd(selectorWidth)}  SUMMARY`]

  for (const row of rows) {
    const confirmation = row.confirm === 'always' ? ' [confirm]' : ''
    lines.push(`${String(row.index).padEnd(indexWidth)}  ${row.uid.padEnd(uidWidth)}  ${row.selector.padEnd(selectorWidth)}  ${row.summary}${confirmation}`)
  }

  return lines.join('\n') + '\n'
}

export const renderDescription = (command: ResolvedCommand): string => [
  `${command.uid} (${command.selector})`,
  '',
  command.description,
  '',
  `index: ${command.index}`,
  `cwd: ${command.cwd}`,
  `shell: ${command.shell ?? 'auto'}`,
  `confirmation: ${command.confirm ?? 'never'}`,
  `tags: ${(command.tags ?? []).join(', ') || 'none'}`,
  `command: ${command.command}`,
].join('\n') + '\n'

export const renderExecutionPlan = (command: ResolvedCommand): string => [
  `uid: ${command.uid}`,
  `selector: ${command.selector}`,
  `manifest: ${command.manifestPath}`,
  `cwd: ${command.cwd}`,
  `shell: ${command.shell ?? 'auto'}`,
  `command: ${command.command}`,
].join('\n') + '\n'
