/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { readManifest, resolveCommand } from './configuration.js'

const directories: string[] = []
const originalFetch = globalThis.fetch

afterEach(async () => {
  globalThis.fetch = originalFetch
  await Promise.all(directories.splice(0).map(directory => rm(directory, { recursive: true, force: true })))
})

describe('RunX manifest v2 composition', () => {
  test('colocates nested groups and mounts a reciprocally declared local child under an alias', async () => {
    const root = await temporaryDirectory()
    const child = join(root, 'packages', 'worker')
    await Bun.write(join(root, 'runx.yaml'), catalog('root-catalog', `
  - group: tools
    summary: Root tools.
    commands:
      - uid: root-check
        id: check
        summary: Check root.
        description: Check the root catalog.
        command: echo root
  - group: worker-alias
    summary: Worker catalog.
    runx: packages/worker/runx.yaml`))
    await Bun.write(join(child, 'runx.yaml'), catalog('worker', `
parent: ../../runx.yaml
commands:
  - group: build
    summary: Build commands.
    commands:
      - uid: worker-compile
        id: compile
        summary: Compile worker.
        description: Compile the worker package.
        command: echo worker`))

    const { manifest, path } = await readManifest(root)
    expect(manifest.namespace).toBe('root-catalog')
    expect(manifest.commands.map(command => command.selector)).toEqual(['tools/check', 'worker-alias/build/compile'])
    expect(manifest.children).toEqual([{
      namespace: 'worker-alias',
      declaredNamespace: 'worker',
      path: join(child, 'runx.yaml'),
      source: 'local',
      parent: join(root, 'runx.yaml'),
    }])
    const resolved = resolveCommand(manifest, path, 'worker-alias/build/compile')
    expect(resolved.uid).toBe('worker-compile')
    expect(resolved.cwd).toBe(child)
    expect(resolved.manifestPath).toBe(join(child, 'runx.yaml'))
  })

  test('validates a child parent declaration without silently replacing the child catalog', async () => {
    const root = await temporaryDirectory()
    const child = join(root, 'child')
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: child
    summary: Child catalog.
    runx: child/runx.yaml`))
    await Bun.write(join(child, 'runx.yaml'), catalog('child-catalog', `
parent: ../runx.yaml
commands:
  - uid: child-test
    id: test
    summary: Test child.
    description: Test only the child catalog.
    command: echo child`))

    const { manifest } = await readManifest(child)
    expect(manifest.namespace).toBe('child-catalog')
    expect(manifest.commands.map(command => command.selector)).toEqual(['test'])
    expect(manifest.parent).toEqual({ path: join(root, 'runx.yaml'), source: 'local' })
  })

  test('rejects legacy split groups, sibling collisions, namespace collisions, and non-GitHub foreign catalogs', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), `version: "1.0.0"\nproject:\n  name: legacy\nscripts:\n  directory: scripts\ngroups:\n  public:\n    summary: Legacy.\ncommands: []\n`)
    await expect(readManifest(root)).rejects.toThrow('Invalid RunX configuration')

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: first
    id: same
    summary: First.
    description: First command.
    command: echo first
  - group: same
    summary: Conflicting group.
    commands: []`))
    await expect(readManifest(root)).rejects.toThrow('Duplicate command or group name "same"')

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: namespace-conflict
    id: root
    summary: Conflict.
    description: Namespace conflict.
    command: echo conflict`))
    await expect(readManifest(root)).rejects.toThrow('Namespace "root" conflicts')

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: foreign
    summary: Foreign catalog.
    runx: https://example.com/runx.yaml`))
    await expect(readManifest(root)).rejects.toThrow('full GitHub URL')
  })

  test('normalizes bounded GitHub children and validates their upstream parent while keeping a local execution base', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('local-root', `
  - group: remote-alias
    summary: Remote child.
    runx: https://github.com/example/catalog/blob/main/child/runx.yaml`))
    const requests: string[] = []
    globalThis.fetch = (async (input) => {
      const url = String(input)
      requests.push(url)
      if (url.endsWith('/child/runx.yaml')) {
        return new Response(catalog('remote-child', `
parent: ../runx.yaml
commands:
  - uid: remote-test
    id: test
    summary: Test remote.
    description: Test the remote child.
    command: echo remote`), { headers: { 'content-length': '240' } })
      }
      return new Response(catalog('remote-root', `
  - group: child
    summary: Child.
    runx: child/runx.yaml`), { headers: { 'content-length': '180' } })
    }) as typeof fetch

    const { manifest, path } = await readManifest(root)
    expect(requests).toEqual([
      'https://raw.githubusercontent.com/example/catalog/main/child/runx.yaml',
      'https://raw.githubusercontent.com/example/catalog/main/runx.yaml',
    ])
    expect(manifest.children[0]).toMatchObject({ namespace: 'remote-alias', declaredNamespace: 'remote-child', source: 'foreign' })
    const command = resolveCommand(manifest, path, 'remote-alias/test')
    expect(command.catalogSource).toBe('foreign')
    expect(command.cwd).toBe(root)
  })
})

function catalog(namespace: string, body: string): string {
  const normalized = body.trimStart()
  const includesCommands = normalized.startsWith('commands:') || normalized.includes('\ncommands:')
  return `version: "2.0.0"
namespace: ${namespace}
scripts:
  directory: scripts
${includesCommands ? normalized : `commands:\n${body}`}
`
}

async function temporaryDirectory(): Promise<string> {
  const directory = await mkdtemp(join(tmpdir(), 'runx-composition-'))
  directories.push(directory)
  return directory
}
