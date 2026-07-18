/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, spyOn, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { runCliWithErrorHandling } from './cli.js'
import { readVersion } from './help.js'

type CliResult = { exitCode: number, stdout: string, stderr: string }

const directories: string[] = []
const cliPath = Bun.fileURLToPath(new URL('./guiho-runx-bin.ts', import.meta.url))

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX RFC 0034 CLI', () => {
  test('prints exact banner and manifest-free help/version', async () => {
    const cwd = await temporaryDirectory()
    expect((await cli([], cwd)).stdout).toBe(`Hello Windows - runx v${readVersion()}\n`)
    for (const flag of ['-v', '--version']) expect((await cli([flag], cwd)).stdout.trim()).toBe(readVersion())
    for (const flag of ['-h', '--help']) expect((await cli([flag], cwd)).stdout).toContain('USAGE')
  })

  test('prints the cached update notice before the no-argument banner without network work', async () => {
    const cwd = await temporaryDirectory()
    const home = await temporaryDirectory()
    await Bun.write(join(home, '.guiho', 'runx', 'cache.json'), JSON.stringify({
      newVersionAvailable: true,
      latestVersion: '99.0.0',
      upgradeCommand: 'runx upgrade',
      lastCheck: '2026-07-18T00:00:00.000Z',
    }))

    const result = await cli([], cwd, home)

    expect(result.exitCode).toBe(0)
    expect(result.stdout).toBe(`New version available. Run this command to upgrade: runx upgrade\nHello Windows - runx v${readVersion()}\n`)
    expect(result.stderr).toBe('')
  })

  test('renders every public scope with all Developer Context help modes', async () => {
    const cwd = await temporaryDirectory()
    const paths = [
      [], ['list'], ['describe'], ['run'], ['check'], ['init'], ['agent'],
      ['agent', 'skill'], ['agent', 'skill', 'install'], ['agent', 'skill', 'uninstall'],
      ['agent', 'skill', 'update'], ['agent', 'skill', 'list'], ['agent', 'skill', 'show'],
      ['agent', 'instruction'], ['agent', 'instruction', 'apply'], ['agent', 'instruction', 'remove'],
      ['agent', 'instruction', 'update'], ['agent', 'instruction', 'show'],
      ['agent', 'prompt'], ['agent', 'prompt', 'list'], ['agent', 'prompt', 'show'],
      ['upgrade'], ['upgrade', 'check'], ['upgrade', 'list'], ['uninstall'],
    ]
    for (const path of paths) {
      expect((await cli([...path, '--help'], cwd)).stdout).toContain('USAGE')
      expect((await cli([...path, '--help-tree'], cwd)).stdout).toStartWith('COMMAND TREE\n\n')
      expect((await cli([...path, '--help-docs'], cwd)).stdout).toStartWith(`# runx${path.length ? ` ${path.join(' ')}` : ''}\n`)
    }
    const depth = await cli(['agent', '--help-tree-depth', '1'], cwd)
    expect(depth.stdout).toContain('├── skill')
    expect(depth.stdout).not.toContain('install')
    expect((await cli(['--help-tree-depth', '0'], cwd)).exitCode).toBe(2)
  }, 30_000)

  test('uses exact YAML precedence and reports the loaded absolute path', async () => {
    const cwd = await temporaryDirectory()
    const home = await temporaryDirectory()
    const explicit = join(cwd, 'selected.yaml')
    await Bun.write(join(cwd, 'runx.yaml'), manifest('cwd-command'))
    await Bun.write(explicit, manifest('explicit-command'))
    await Bun.write(join(home, '.guiho', 'runx', 'runx.yaml'), manifest('global-command'))

    const selected = await cli(['list', '--config', explicit, '--format', 'json'], cwd, home)
    expect(selected.exitCode).toBe(0)
    expect(selected.stderr).toBe(`configuration file loaded: ${explicit}\n`)
    expect(JSON.parse(selected.stdout).manifest.commands[0].uid).toBe('explicit-command')

    const local = await cli(['list', '--format', 'json'], cwd, home)
    expect(JSON.parse(local.stdout).manifest.commands[0].uid).toBe('cwd-command')
  })

  test('does not search parents and maps configuration errors to exit 3', async () => {
    const parent = await temporaryDirectory()
    const child = join(parent, 'child')
    await Bun.write(join(parent, 'runx.yaml'), manifest('parent-command'))
    await Bun.write(join(child, '.keep'), '')
    expect((await cli(['list'], child, await temporaryDirectory())).exitCode).toBe(3)
    await Bun.write(join(child, 'runx.yaml'), 'version: [')
    expect((await cli(['check'], child)).exitCode).toBe(3)
  })

  test('keeps inspection and dry-run read-only while preserving delegated exit codes', async () => {
    const cwd = await temporaryDirectory()
    await Bun.write(join(cwd, 'runx.yaml'), manifest('selected', 'exit 7'))
    expect((await cli(['list'], cwd)).stdout).not.toContain('Running selected')
    expect((await cli(['run', 'selected', '--dry-run'], cwd)).exitCode).toBe(0)
    expect((await cli(['run', 'selected'], cwd)).exitCode).toBe(7)
  })

  test('rejects removed aliases, arbitrary short flags, unknown commands, and invalid enums', async () => {
    const cwd = await temporaryDirectory()
    await Bun.write(join(cwd, 'runx.yaml'), manifest('selected'))
    for (const args of [['r', 'selected'], ['selected'], ['agents'], ['list', '--file', 'runx.yaml'], ['-x']]) {
      expect((await cli(args, cwd)).exitCode).toBe(2)
    }
    expect((await cli(['list', '--format', 'xml'], cwd)).exitCode).toBe(2)
  })

  test('implements local dual-tool skills, instructions, and raw prompts idempotently', async () => {
    const cwd = await temporaryDirectory()
    expect((await cli(['agent', 'skill', 'install', '--local', '--cwd', cwd], cwd)).exitCode).toBe(0)
    expect(await Bun.file(join(cwd, '.agents', 'skills', 'guiho-s-runx', 'SKILL.md')).exists()).toBe(true)
    expect(await Bun.file(join(cwd, '.claude', 'skills', 'guiho-s-runx', 'SKILL.md')).exists()).toBe(true)
    await Bun.write(join(cwd, 'AGENTS.md'), '# Agent\n')
    await Bun.write(join(cwd, 'CLAUDE.md'), '# Claude\n')
    await cli(['agent', 'instruction', 'apply', '--cwd', cwd], cwd)
    await cli(['agent', 'instruction', 'update', '--cwd', cwd], cwd)
    for (const name of ['AGENTS.md', 'CLAUDE.md']) {
      const text = await Bun.file(join(cwd, name)).text()
      expect(text.match(/BEGIN RUNX — DO NOT EDIT THIS SECTION/g)).toHaveLength(1)
    }
    expect((await cli(['agent', 'prompt', 'list', '--names'], cwd)).stdout).toBe('guiho-i-runx\n')
    expect(JSON.parse((await cli(['agent', 'prompt', 'list', '--names', '--format', 'json'], cwd)).stdout)).toEqual(['guiho-i-runx'])
    expect((await cli(['agent', 'prompt', 'show', 'guiho-i-runx'], cwd)).stdout).toContain('# RunX Agent Instruction')
  })

  test('routes public upgrade flags to the upgrade command without root version interception', async () => {
    const cwd = await temporaryDirectory()
    const originalFetch = globalThis.fetch
    const originalHome = Bun.env.HOME
    const originalSelfPath = Bun.env.RUNX_SELF_PATH
    const originalWorker = Bun.env.RUNX_DISABLE_UPDATE_WORKER
    const originalExitCode = process.exitCode
    let stdout = ''
    let stderr = ''
    const stdoutSpy = spyOn(process.stdout, 'write').mockImplementation((chunk) => {
      stdout += String(chunk)
      return true
    })
    const stderrSpy = spyOn(process.stderr, 'write').mockImplementation((chunk) => {
      stderr += String(chunk)
      return true
    })

    try {
      Bun.env.HOME = cwd
      Bun.env.RUNX_SELF_PATH = join(cwd, process.platform === 'win32' ? 'runx.exe' : 'runx')
      Bun.env.RUNX_DISABLE_UPDATE_WORKER = '1'
      globalThis.fetch = releaseFetch(readVersion())
      process.exitCode = 0

      await runCliWithErrorHandling(['upgrade', '--version', readVersion(), '--arch', 'x64', '--variant', 'baseline', '--dry-run', '--format', 'json'])

      expect(process.exitCode).toBe(0)
      const output = JSON.parse(stdout)
      expect(output.command).toBe('runx upgrade')
      expect(output.outcome).toBe('up-to-date')
      expect(output.plan.targetVersion).toBe(readVersion())
      expect(stdout.trim()).not.toBe(readVersion())
      expect(stderr).toBe('')
    } finally {
      stdoutSpy.mockRestore()
      stderrSpy.mockRestore()
      globalThis.fetch = originalFetch
      restoreEnvironment('HOME', originalHome)
      restoreEnvironment('RUNX_SELF_PATH', originalSelfPath)
      restoreEnvironment('RUNX_DISABLE_UPDATE_WORKER', originalWorker)
      process.exitCode = originalExitCode
    }
  })
})

async function cli(args: string[], cwd: string, home = cwd): Promise<CliResult> {
  const child = Bun.spawn([process.execPath, cliPath, ...args], {
    cwd,
    env: { ...process.env, HOME: home, USERPROFILE: home, NO_COLOR: '1', FORCE_COLOR: '0', RUNX_DISABLE_UPDATE_WORKER: '1' },
    stdout: 'pipe',
    stderr: 'pipe',
  })
  const [exitCode, stdout, stderr] = await Promise.all([child.exited, new Response(child.stdout).text(), new Response(child.stderr).text()])
  return { exitCode, stdout, stderr }
}

async function temporaryDirectory(): Promise<string> {
  const directory = await mkdtemp(join(tmpdir(), 'runx-rfc-'))
  directories.push(directory)
  return directory
}

function manifest(uid: string, command = 'echo safe'): string {
  return `version: "1.0.0"
scripts:
  directory: scripts
groups:
  public:
    summary: Public commands.
commands:
  - uid: ${uid}
    id: selected
    group: public
    summary: Selected command.
    description: Selected command for tests.
    command: ${command}
`
}

function releaseFetch(version: string): typeof fetch {
  return (async () => Response.json([{
    tag_name: `@guiho/runx@${version}`,
    draft: false,
    prerelease: false,
    published_at: '2026-07-18T00:00:00Z',
    assets: [{
      name: process.platform === 'win32' ? 'runx-windows-x64-baseline.exe' : process.platform === 'darwin' ? 'runx-darwin-x64-baseline' : 'runx-linux-x64-baseline',
      browser_download_url: 'https://example.test/runx',
    }],
  }])) as typeof fetch
}

function restoreEnvironment(name: string, value: string | undefined): void {
  if (value === undefined) delete Bun.env[name]
  else Bun.env[name] = value
}
