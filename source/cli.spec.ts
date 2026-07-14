/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { readVersion } from './help.js'

type CliResult = {
  exitCode: number
  stdout: string
  stderr: string
}

const directories: string[] = []
const cliPath = Bun.fileURLToPath(new URL('./guiho-runx-bin.ts', import.meta.url))

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX Citty CLI', () => {
  test('shows version and root help without discovering a manifest', async () => {
    const cwd = await emptyDirectory()

    for (const flag of ['-v', '--version']) {
      const result = await cli([flag], cwd)
      expect(result.exitCode).toBe(0)
      expect(result.stdout.trim()).toBe(readVersion())
      expect(result.stderr).toBe('')
    }

    for (const flag of ['-h', '--help']) {
      const result = await cli([flag], cwd)
      expect(result.exitCode).toBe(0)
      expect(result.stdout).toContain('USAGE')
      expect(result.stdout).toContain('runx [OPTIONS]')
      expect(result.stderr).toBe('')
    }
  })

  test('shows the home page, help tree, and documentation without a manifest', async () => {
    const cwd = await emptyDirectory()
    const home = await cli([], cwd)
    const tree = await cli(['--help-tree'], cwd)
    const docs = await cli(['--help-docs'], cwd)

    expect(home.exitCode).toBe(0)
    expect(home.stdout).toContain(`RunX ${readVersion()}`)
    expect(tree.exitCode).toBe(0)
    expect(tree.stdout).toContain('|- agents')
    expect(docs.exitCode).toBe(0)
    expect(docs.stdout).toContain('Manifest: runx.yaml')
  })

  test('renders Citty help for every public command and nested route', async () => {
    const cwd = await emptyDirectory()
    const paths = [
      ['list'],
      ['describe'],
      ['run'],
      ['r'],
      ['check'],
      ['agents'],
      ['agents', 'install'],
      ['agents', 'instructions'],
      ['upgrade'],
      ['upgrade', 'check'],
      ['upgrade', 'list'],
      ['uninstall'],
    ]

    for (const path of paths) {
      const result = await cli([...path, '--help'], cwd)
      expect(result.exitCode).toBe(0)
      expect(result.stdout).toContain('USAGE')
      expect(result.stdout).toContain(`runx ${path.join(' ')}`)
      expect(result.stderr).toBe('')
    }
  })

  test('lists, checks, and describes a manifest with text and JSON output', async () => {
    const cwd = await projectDirectory()
    const manifestPath = join(cwd, 'runx.yaml')
    const list = await cli(['--file', manifestPath, 'list', '--format', 'json'], cwd)
    const check = await cli(['check', '--file', manifestPath], cwd)
    const describe = await cli(['describe', 'dev-start', '--file', manifestPath, '--format=json'], cwd)

    expect(list.exitCode).toBe(0)
    expect(JSON.parse(list.stdout).manifest.commands).toHaveLength(3)
    expect(check.exitCode).toBe(0)
    expect(check.stdout).toContain('valid: true')
    expect(JSON.parse(describe.stdout).uid).toBe('dev-start')
  })

  test('runs the explicit command, alias, and root selector shorthand', async () => {
    const cwd = await projectDirectory()
    const manifestPath = join(cwd, 'runx.yaml')

    for (const args of [
      ['run', 'dev-start', '--file', manifestPath],
      ['r', 'dev-start', '--file', manifestPath],
      ['dev-start', '--file', manifestPath],
    ]) {
      const result = await cli(args, cwd)
      expect(result.exitCode).toBe(0)
      expect(result.stdout).toContain('Running dev-start')
      expect(result.stdout).toContain('start')
    }

    const prototypeNamedSelector = await cli(['constructor', '--file', manifestPath], cwd)
    expect(prototypeNamedSelector.exitCode).toBe(0)
    expect(prototypeNamedSelector.stdout).toContain('Running constructor')
  })

  test('keeps dry runs read-only and confirmation gates explicit', async () => {
    const cwd = await projectDirectory()
    const manifestPath = join(cwd, 'runx.yaml')
    const dryRun = await cli(['run', 'confirmed', '--file', manifestPath, '--dry-run'], cwd)
    const blocked = await cli(['run', 'confirmed', '--file', manifestPath], cwd)
    const approved = await cli(['run', 'confirmed', '--file', manifestPath, '--yes'], cwd)

    expect(dryRun.exitCode).toBe(0)
    expect(dryRun.stdout).toContain('echo confirmed')
    expect(dryRun.stdout).not.toContain('Running confirmed')
    expect(blocked.exitCode).toBe(1)
    expect(blocked.stderr).toContain('requires confirmation')
    expect(approved.exitCode).toBe(0)
    expect(approved.stdout).toContain('confirmed')
  })

  test('reports Citty usage for unknown options and missing arguments', async () => {
    const cwd = await emptyDirectory()
    const unknown = await cli(['--unknown-option'], cwd)
    const missing = await cli(['describe'], cwd)
    const nested = await cli(['agents', 'unknown'], cwd)

    expect(unknown.exitCode).toBe(1)
    expect(unknown.stderr).toContain('USAGE')
    expect(unknown.stderr).toContain('Unknown option --unknown-option')
    expect(unknown.stderr).not.toContain('No runx.yaml found')
    expect(missing.exitCode).toBe(1)
    expect(missing.stderr).toContain('runx describe')
    expect(missing.stderr).toContain('Missing required positional argument: SELECTOR')
    expect(nested.exitCode).toBe(1)
    expect(nested.stderr).toContain('runx agents')
    expect(nested.stderr).toContain('Unknown command')
  })

  test('routes local agent installation and managed instructions', async () => {
    const cwd = await emptyDirectory()
    const install = await cli(['agents', 'install', 'local', '--cwd', cwd, '--tool', 'all', '--format', 'json'], cwd)
    const instructions = await cli(['agents', 'instructions', '--cwd', cwd, '--format=json'], cwd)
    const installed = JSON.parse(install.stdout).installed as string[]

    expect(install.exitCode).toBe(0)
    expect(installed).toHaveLength(2)
    expect(await Bun.file(join(cwd, '.agents', 'skills', 'guiho-s-runx', 'SKILL.md')).exists()).toBe(true)
    expect(await Bun.file(join(cwd, '.claude', 'skills', 'guiho-s-runx', 'SKILL.md')).exists()).toBe(true)
    expect(instructions.exitCode).toBe(0)
    expect(await Bun.file(join(cwd, 'AGENTS.md')).text()).toContain('BEGIN RUNX AGENT INSTRUCTIONS')
  })
})

async function cli(args: string[], cwd: string): Promise<CliResult> {
  const subprocess = Bun.spawn([process.execPath, cliPath, ...args], {
    cwd,
    env: { ...process.env, NO_COLOR: '1', FORCE_COLOR: '0' },
    stdout: 'pipe',
    stderr: 'pipe',
  })
  const [exitCode, stdout, stderr] = await Promise.all([
    subprocess.exited,
    new Response(subprocess.stdout).text(),
    new Response(subprocess.stderr).text(),
  ])
  return { exitCode, stdout, stderr }
}

async function emptyDirectory(): Promise<string> {
  const directory = await mkdtemp(join(tmpdir(), 'runx-cli-'))
  directories.push(directory)
  return directory
}

async function projectDirectory(): Promise<string> {
  const directory = await emptyDirectory()
  await Bun.write(join(directory, 'runx.yaml'), manifest())
  return directory
}

function manifest(): string {
  return `version: 1
groups:
  development:
    summary: Development commands.
commands:
  - uid: dev-start
    id: start
    group: development
    summary: Start development.
    description: Starts local development.
    command: echo start
  - uid: confirmed
    id: confirmed
    group: development
    summary: Run a confirmed command.
    description: Exercises the confirmation gate.
    command: echo confirmed
    confirm: always
  - uid: constructor
    id: constructor
    group: development
    summary: Exercise selector compatibility.
    description: Verifies that object prototype names remain valid selectors.
    command: echo constructor
`
}
