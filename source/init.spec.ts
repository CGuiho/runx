/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { initializeRunXManifest } from './init.js'

let directory = ''

afterEach(async () => {
  if (directory) await rm(directory, { recursive: true, force: true })
  directory = ''
})

test('init creates only the selected YAML configuration', async () => {
  directory = await mkdtemp(join(tmpdir(), 'runx-init-'))
  const result = await initializeRunXManifest({ cwd: directory, config: 'catalog/runx.yaml' })
  expect(result.status).toBe('created')
  expect(result.path).toBe(join(directory, 'catalog', 'runx.yaml'))
  const manifest = Bun.YAML.parse(await Bun.file(result.path).text())
  expect(manifest).toMatchObject({ version: '2.0.0', namespace: 'catalog', commands: [] })
  expect(manifest.project).toBeUndefined()
  expect(manifest.groups).toBeUndefined()
  expect(await Bun.file(join(directory, 'runx.toml')).exists()).toBe(false)
  expect(await Bun.file(join(directory, 'runx.json')).exists()).toBe(false)
})

test('init derives a legal deterministic namespace from the target catalog directory', async () => {
  directory = await mkdtemp(join(tmpdir(), 'runx-init-'))
  const cases = [
    { config: 'My Project/runx.yaml', namespace: 'my-project' },
    { config: '!!!/runx.yaml', namespace: 'runx' },
    { config: '42 APIs/runx.yaml', namespace: 'n-42-apis' },
    { config: 'packages/Worker API/runx.yaml', namespace: 'worker-api' },
  ]
  for (const entry of cases) {
    const result = await initializeRunXManifest({ cwd: directory, config: entry.config })
    const manifest = Bun.YAML.parse(await Bun.file(result.path).text())
    expect(manifest.namespace).toBe(entry.namespace)
  }
})

test('init refuses to overwrite an existing configuration', async () => {
  directory = await mkdtemp(join(tmpdir(), 'runx-init-'))
  await Bun.write(join(directory, 'runx.yaml'), 'existing')
  await expect(initializeRunXManifest({ cwd: directory })).rejects.toThrow('already exists')
  expect(await Bun.file(join(directory, 'runx.yaml')).text()).toBe('existing')
})
