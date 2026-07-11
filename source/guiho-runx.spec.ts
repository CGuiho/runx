import { afterEach, describe, expect, test } from 'bun:test'
import { mkdir, mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { RunXError } from './errors.js'
import { booleanFlag, parseArgs, stringFlag } from './flags.js'
import { readManifest, resolveCommand } from './manifest.js'

const directories: string[] = []

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX manifests', () => {
  test('discovers a parent manifest and resolves stable selectors', async () => {
    const root = await fixture(manifest())
    const nested = join(root, 'nested')
    await mkdir(nested)
    const { manifest: loaded, path } = await readManifest(nested)

    expect(resolveCommand(loaded, path, 'dev-start').selector).toBe('development/start')
    expect(resolveCommand(loaded, path, 'development/check').uid).toBe('check')
    expect(resolveCommand(loaded, path, '2').uid).toBe('check')
  })

  test('rejects ambiguous IDs', async () => {
    const root = await fixture(`${manifest()}\n  - uid: maintenance-start\n    id: start\n    group: maintenance\n    summary: Start maintenance.\n    description: Starts local maintenance.\n    command: echo maintenance\n`.replace('  development:\n    summary: Development commands.', '  development:\n    summary: Development commands.\n  maintenance:\n    summary: Maintenance commands.'))
    const { manifest: loaded, path } = await readManifest(root)
    expect(() => resolveCommand(loaded, path, 'start')).toThrow(RunXError)
  })

  test('parses global flags and the run alias', () => {
    const parsed = parseArgs(['r', 'dev-start', '--dry-run', '--format=json'])
    expect(parsed.command).toBe('r')
    expect(parsed.positionals).toEqual(['dev-start'])
    expect(booleanFlag(parsed.flags, 'dryRun')).toBe(true)
    expect(stringFlag(parsed.flags, 'format')).toBe('json')
  })

  test('runs the r alias and keeps dry runs read-only', async () => {
    const root = await fixture(manifest())
    const manifestPath = join(root, 'runx.yaml')
    const cliPath = Bun.fileURLToPath(new URL('./guiho-runx-bin.ts', import.meta.url))
    const dryRun = Bun.spawn([process.execPath, cliPath, 'r', 'dev-start', '--file', manifestPath, '--dry-run'], { stdout: 'pipe', stderr: 'pipe' })
    expect(await dryRun.exited).toBe(0)
    expect(await new Response(dryRun.stdout).text()).toContain('uid: dev-start')

    const run = Bun.spawn([process.execPath, cliPath, 'r', 'dev-start', '--file', manifestPath], { stdout: 'pipe', stderr: 'pipe' })
    expect(await run.exited).toBe(0)
    expect(await new Response(run.stdout).text()).toContain('start')
  })
})

const fixture = async (content: string): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), 'runx-'))
  directories.push(directory)
  await Bun.write(join(directory, 'runx.yaml'), content)
  return directory
}

const manifest = (): string => `version: 1
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
  - uid: check
    id: check
    group: development
    summary: Check source.
    description: Checks source files.
    command: echo check
`
