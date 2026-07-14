/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { createInitialManifest, initializeRunXManifest, renderInitialManifest } from './init.js'
import { readManifest } from './manifest.js'

import type { RunXInitPrompter } from './init.js'

const directories: string[] = []

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX initialization', () => {
  test('creates the exact empty manifest without creating the scripts directory', async () => {
    const root = await projectDirectory()
    const prompter = fakePrompter(['sample-project', 'scripts'], [true])

    const result = await initializeRunXManifest({ cwd: root, isInteractive: true, prompter })
    const path = join(root, 'runx.yaml')

    expect(result).toEqual({ status: 'created', path, manifest: createInitialManifest('sample-project', 'scripts') })
    expect(await Bun.file(path).text()).toBe(renderInitialManifest(createInitialManifest('sample-project', 'scripts')))
    expect(await Bun.file(join(root, 'scripts')).exists()).toBe(false)
    await expect(readManifest(root)).resolves.toMatchObject({ manifest: { commands: [], groups: { public: { summary: 'Default public project commands.' } } } })
    expect(prompter.previews).toEqual([renderInitialManifest(createInitialManifest('sample-project', 'scripts'))])
  })

  test('leaves no manifest when the user cancels', async () => {
    const root = await projectDirectory()
    const prompter = fakePrompter([undefined], [])

    await expect(initializeRunXManifest({ cwd: root, isInteractive: true, prompter })).resolves.toEqual({ status: 'cancelled' })
    expect(await Bun.file(join(root, 'runx.yaml')).exists()).toBe(false)
    expect(prompter.cancelled).toEqual(['Initialization cancelled. No files were changed.'])
  })

  test('preserves an existing manifest until overwrite is explicitly confirmed', async () => {
    const root = await projectDirectory()
    const path = join(root, 'runx.yaml')
    await Bun.write(path, 'version: "1.0.0"\n')

    const rejected = fakePrompter([], [false])
    await expect(initializeRunXManifest({ cwd: root, isInteractive: true, prompter: rejected })).resolves.toEqual({ status: 'cancelled' })
    expect(await Bun.file(path).text()).toBe('version: "1.0.0"\n')

    const accepted = fakePrompter(['replaced-project', 'automation'], [true, true])
    await expect(initializeRunXManifest({ cwd: root, isInteractive: true, prompter: accepted })).resolves.toMatchObject({ status: 'created', path })
    await expect(readManifest(root)).resolves.toMatchObject({ manifest: { project: { name: 'replaced-project' }, scripts: { directory: 'automation' } } })
  })
})

type FakePrompter = RunXInitPrompter & { previews: string[], cancelled: string[] }

const fakePrompter = (textValues: Array<string | undefined>, confirmationValues: Array<boolean | undefined>): FakePrompter => {
  const previews: string[] = []
  const cancelled: string[] = []
  return {
    previews,
    cancelled,
    intro: () => undefined,
    text: async (options) => {
      const value = textValues.shift()
      if (value !== undefined) expect(options.validate(value)).toBeUndefined()
      return value
    },
    confirm: async () => confirmationValues.shift(),
    preview: (manifest) => previews.push(manifest),
    cancel: (message) => cancelled.push(message),
    outro: () => undefined,
  }
}

const projectDirectory = async (): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), 'runx-init-'))
  directories.push(directory)
  return directory
}
