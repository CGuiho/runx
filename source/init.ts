/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { rename, rm, writeFile } from 'node:fs/promises'
import { basename, isAbsolute, relative, resolve } from 'node:path'
import { createInterface } from 'node:readline/promises'
import { RunXError } from './errors.js'
import { readManifest } from './manifest.js'

import type { RunXManifest } from './types.js'

export type RunXInitPrompter = {
  intro(title: string): void
  text(options: { message: string, initialValue: string, validate: (value: string | undefined) => string | undefined }): Promise<string | undefined>
  confirm(options: { message: string, initialValue: boolean }): Promise<boolean | undefined>
  preview(manifest: string): void
  cancel(message: string): void
  outro(message: string): void
  close?(): void
}

export type RunXInitResult =
  | { status: 'cancelled' }
  | { status: 'created', path: string, manifest: RunXManifest }

type RunXInitOptions = {
  cwd: string
  isInteractive?: boolean
  prompter?: RunXInitPrompter
}

const manifestVersion = '1.0.0'

export const initializeRunXManifest = async ({ cwd, isInteractive = process.stdin.isTTY === true && process.stdout.isTTY === true, prompter }: RunXInitOptions): Promise<RunXInitResult> => {
  if (!isInteractive) {
    throw new RunXError('runx init requires an interactive terminal. Run it directly in a terminal with stdin and stdout attached.')
  }

  const activePrompter = prompter ?? createTerminalPrompter()

  const root = resolve(cwd)
  const path = resolve(root, 'runx.yaml')
  const defaultProjectName = basename(root) || 'my-project'

  activePrompter.intro('RunX project setup')

  if (await Bun.file(path).exists()) {
    const overwrite = await activePrompter.confirm({
      message: 'runx.yaml already exists. Replace it?',
      initialValue: false,
    })
    if (overwrite !== true) return cancelInitialization(activePrompter)
  }

  const projectName = await activePrompter.text({
    message: 'Project name',
    initialValue: defaultProjectName,
    validate: (value) => value?.trim() ? undefined : 'Enter a project name.',
  })
  if (projectName === undefined) return cancelInitialization(activePrompter)

  const scriptsDirectory = await activePrompter.text({
    message: 'Scripts directory',
    initialValue: 'scripts',
    validate: (value) => validateScriptsDirectory(value, root),
  })
  if (scriptsDirectory === undefined) return cancelInitialization(activePrompter)

  const manifest = createInitialManifest(projectName.trim(), scriptsDirectory.trim())
  const yaml = renderInitialManifest(manifest)
  activePrompter.preview(yaml)

  const confirmed = await activePrompter.confirm({
    message: 'Create this runx.yaml?',
    initialValue: true,
  })
  if (confirmed !== true) return cancelInitialization(activePrompter)

  try {
    await writeValidatedManifest(root, path, yaml)
  } catch (error) {
    activePrompter.close?.()
    throw error
  }
  activePrompter.outro(`Created ${path}`)
  activePrompter.close?.()
  return { status: 'created', path, manifest }
}

export const createInitialManifest = (projectName: string, scriptsDirectory: string): RunXManifest => ({
  version: manifestVersion,
  project: { name: projectName },
  scripts: { directory: scriptsDirectory },
  groups: {
    public: { summary: 'Default public project commands.' },
  },
  commands: [],
})

export const renderInitialManifest = (manifest: RunXManifest): string => [
  `version: ${JSON.stringify(manifest.version)}`,
  '',
  'project:',
  `  name: ${JSON.stringify(manifest.project?.name ?? '')}`,
  '',
  'scripts:',
  `  directory: ${JSON.stringify(manifest.scripts.directory)}`,
  '',
  'groups:',
  '  public:',
  `    summary: ${JSON.stringify(manifest.groups.public.summary)}`,
  '',
  'commands: []',
  '',
].join('\n')

const createTerminalPrompter = (): RunXInitPrompter => {
  const terminal = createInterface({ input: process.stdin, output: process.stdout, terminal: true })
  const accent = (value: string): string => process.env['NO_COLOR'] ? value : `\u001b[36m${value}\u001b[0m`
  const success = (value: string): string => process.env['NO_COLOR'] ? value : `\u001b[32m${value}\u001b[0m`
  const muted = (value: string): string => process.env['NO_COLOR'] ? value : `\u001b[2m${value}\u001b[0m`
  const error = (value: string): string => process.env['NO_COLOR'] ? value : `\u001b[31m${value}\u001b[0m`
  const write = (value: string): void => {
    process.stdout.write(value)
  }

  const ask = async (label: string): Promise<string | undefined> => {
    try {
      return await terminal.question(label)
    } catch {
      return undefined
    }
  }

  return {
    intro: (title) => write(`\n${accent('◆')} ${title}\n${muted('  A documented command catalog for this project.')}\n\n`),
    text: async (options) => {
      while (true) {
        const value = await ask(`${accent('◆')} ${options.message} ${muted(`[${options.initialValue}]`)}\n  `)
        if (value === undefined) return undefined
        const selected = value.trim() || options.initialValue
        const validation = options.validate(selected)
        if (!validation) return selected
        write(`${error('▲')} ${validation}\n`)
      }
    },
    confirm: async (options) => {
      while (true) {
        const defaultLabel = options.initialValue ? 'Y/n' : 'y/N'
        const value = await ask(`${accent('◆')} ${options.message} ${muted(`[${defaultLabel}]`)}\n  `)
        if (value === undefined) return undefined
        const answer = value.trim().toLowerCase()
        if (!answer) return options.initialValue
        if (answer === 'y' || answer === 'yes') return true
        if (answer === 'n' || answer === 'no') return false
        write(`${error('▲')} Enter yes or no.\n`)
      }
    },
    preview: (manifest) => write(`\n${accent('┌─ runx.yaml preview ─────────────────────────────────')}\n${manifest.trimEnd().split('\n').map((line) => `${accent('│')} ${line}`).join('\n')}\n${accent('└────────────────────────────────────────────────────')}\n\n`),
    cancel: (message) => write(`${error('■')} ${message}\n`),
    outro: (message) => {
      write(`${success('◆')} ${message}\n`)
    },
    close: () => terminal.close(),
  }
}

const cancelInitialization = (prompter: RunXInitPrompter): RunXInitResult => {
  prompter.close?.()
  prompter.cancel('Initialization cancelled. No files were changed.')
  return { status: 'cancelled' }
}

const validateScriptsDirectory = (value: string | undefined, root: string): string | undefined => {
  const directory = value?.trim()
  if (!directory) return 'Enter a scripts directory.'
  const resolvedDirectory = resolve(root, directory)
  const relativeDirectory = relative(root, resolvedDirectory)
  if (relativeDirectory === '' || relativeDirectory === '.' || isAbsolute(relativeDirectory) || relativeDirectory.startsWith('..')) {
    return 'Use a relative subdirectory inside this project.'
  }
  return undefined
}

const writeValidatedManifest = async (root: string, path: string, yaml: string): Promise<void> => {
  const temporaryPath = resolve(root, `.runx-${crypto.randomUUID()}.yaml.tmp`)
  try {
    await writeFile(temporaryPath, yaml, { encoding: 'utf8', flag: 'wx' })
    await readManifest(root, temporaryPath)
    await rename(temporaryPath, path)
  } finally {
    await rm(temporaryPath, { force: true })
  }
}
