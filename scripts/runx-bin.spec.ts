/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'

let directory = ''

afterEach(async () => {
  if (directory) await rm(directory, { recursive: true, force: true })
  directory = ''
})

test('packed npm bootstrap downloads and delegates with Node while Bun is absent from PATH', async () => {
  const node = Bun.which('node')
  if (!node) throw new Error('Node is required for the npm bootstrap test.')
  directory = await mkdtemp(join(tmpdir(), 'runx-npm-'))
  const root = Bun.fileURLToPath(new URL('..', import.meta.url))
  const packed = Bun.spawn([process.execPath, 'pm', 'pack', '--destination', directory, '--ignore-scripts'], {
    cwd: root,
    stdout: 'pipe',
    stderr: 'pipe',
  })
  const [packCode, packOutput, packError] = await Promise.all([
    packed.exited,
    new Response(packed.stdout).text(),
    new Response(packed.stderr).text(),
  ])
  expect(packCode, packError).toBe(0)
  const archiveName = [...new Bun.Glob('*.tgz').scanSync({ cwd: directory, onlyFiles: true })][0]
  expect(archiveName, packOutput).toBeDefined()
  const extraction = Bun.spawn(['tar', '-xzf', join(directory, archiveName), '-C', directory])
  expect(await extraction.exited).toBe(0)

  const server = Bun.serve({ port: 0, fetch: () => new Response(Bun.file(process.execPath)) })
  try {
    const wrapper = join(directory, 'package', 'scripts', 'runx-bin.mjs')
    const child = Bun.spawn([node, wrapper, '--version'], {
      cwd: join(directory, 'package'),
      env: {
        ...process.env,
        PATH: join(node, '..'),
        RUNX_NATIVE_CACHE: join(directory, 'native-cache'),
        RUNX_DOWNLOAD_BASE_URL: `http://127.0.0.1:${server.port}`,
      },
      stdout: 'pipe',
      stderr: 'pipe',
    })
    const [exitCode, stdout, stderr] = await Promise.all([
      child.exited,
      new Response(child.stdout).text(),
      new Response(child.stderr).text(),
    ])
    expect(exitCode, stderr).toBe(0)
    expect(stdout.trim()).toBe(process.versions.bun)
  } finally {
    server.stop(true)
  }
})
