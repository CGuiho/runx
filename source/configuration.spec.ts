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
      declaredParent: join(root, 'runx.yaml'),
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
      'https://raw.githubusercontent.com/example/catalog/main/child/runx.yaml',
    ])
    expect(manifest.children[0]).toMatchObject({ namespace: 'remote-alias', declaredNamespace: 'remote-child', source: 'foreign' })
    const command = resolveCommand(manifest, path, 'remote-alias/test')
    expect(command.catalogSource).toBe('foreign')
    expect(command.cwd).toBe(root)
  })

  test('maps foreign transport failures to configuration errors', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: foreign
    summary: Foreign catalog.
    runx: https://github.com/example/catalog/blob/main/runx.yaml`))
    globalThis.fetch = (async () => { throw new Error('network unavailable') }) as typeof fetch
    await expect(readManifest(root)).rejects.toMatchObject({ exitCode: 3 })
    await expect(readManifest(root)).rejects.toThrow('network unavailable')
  })

  test('uses one collision domain without rejecting aliases owned by the same leaf', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: test
    id: test
    summary: Same leaf aliases.
    description: UID and selector belong to one command.
    command: echo test`))
    expect((await readManifest(root)).manifest.commands).toHaveLength(1)

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: shared
    id: first
    summary: First.
    description: First command.
    command: echo first
  - group: nested
    summary: Nested.
    commands:
      - uid: second
        id: shared
        summary: Second.
        description: Second command.
        command: echo second`))
    await expect(readManifest(root)).rejects.toThrow('ID shorthand "shared" conflicts')

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: nested
    summary: Nested.
    commands:
      - uid: first
        id: shared
        summary: First.
        description: First command.
        command: echo first
  - uid: shared
    id: second
    summary: Second.
    description: Second command.
    command: echo second`))
    await expect(readManifest(root)).rejects.toThrow('UID "shared" conflicts')
  })

  test('validates child namespaces before applying their mount prefixes', async () => {
    const root = await temporaryDirectory()
    const child = join(root, 'child')
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: alias
    summary: Child.
    runx: child/runx.yaml`))
    await Bun.write(join(child, 'runx.yaml'), catalog('child-name', `
parent: ../runx.yaml
commands:
  - uid: child-conflict
    id: child-name
    summary: Conflict.
    description: Child namespace conflict.
    command: echo conflict`))
    await expect(readManifest(root)).rejects.toThrow('Namespace "child-name" conflicts')
  })

  test('rejects cross-mount identity shadows and semantically invalid declared parents', async () => {
    const root = await temporaryDirectory()
    const child = join(root, 'child')
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: shared
    id: root-command
    summary: Root.
    description: Root command.
    command: echo root
  - group: child
    summary: Child.
    runx: child/runx.yaml`))
    await Bun.write(join(child, 'runx.yaml'), catalog('child-catalog', `
parent: ../runx.yaml
commands:
  - uid: child-command
    id: shared
    summary: Child.
    description: Child command.
    command: echo child`))
    await expect(readManifest(root)).rejects.toThrow('ID shorthand "shared" conflicts')

    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: child
    summary: Invalid child mount.
    runx: child/runx.yaml
    commands: []`))
    await expect(readManifest(child)).rejects.toThrow('must define exactly one of commands or runx')
  })

  test('accepts 32 nested groups and rejects the 33rd across inline depth', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), nestedCatalog(32))
    expect((await readManifest(root)).manifest.commands[0]?.selector?.split('/')).toHaveLength(33)
    await Bun.write(join(root, 'runx.yaml'), nestedCatalog(33))
    await expect(readManifest(root)).rejects.toThrow('depth exceeds 32')
  })

  test('rejects escaping cwd during composition for root, local child, and foreign child commands', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - uid: root-escape
    id: escape
    summary: Escape.
    description: Escape root.
    command: echo escape
    cwd: ..`))
    await expect(readManifest(root)).rejects.toThrow('cwd outside its catalog directory')

    const child = join(root, 'child')
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: child
    summary: Child.
    runx: child/runx.yaml`))
    await Bun.write(join(child, 'runx.yaml'), catalog('child-catalog', `
parent: ../runx.yaml
commands:
  - uid: child-escape
    id: escape
    summary: Escape.
    description: Escape child.
    command: echo escape
    cwd: ..`))
    await expect(readManifest(root)).rejects.toThrow('cwd outside its catalog directory')

    globalThis.fetch = foreignGraphFetch({ commandCwd: '..' })
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: foreign
    summary: Foreign.
    runx: https://github.com/example/catalog/blob/main/child/runx.yaml`))
    await expect(readManifest(root)).rejects.toThrow('cwd outside its catalog directory')
  })

  test('validates scripts.directory for foreign catalogs and streams the one-MiB limit', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: foreign
    summary: Foreign.
    runx: https://github.com/example/catalog/blob/main/child/runx.yaml`))
    for (const scriptsDirectory of ['.', '..', 'C:\\absolute']) {
      globalThis.fetch = foreignGraphFetch({ scriptsDirectory })
      await expect(readManifest(root)).rejects.toThrow('Scripts directory must be a relative subdirectory')
    }

    let pulls = 0
    globalThis.fetch = (async () => new Response(new ReadableStream({
      pull(controller) {
        pulls += 1
        controller.enqueue(new Uint8Array(600_000))
        if (pulls === 2) controller.close()
      },
    }))) as typeof fetch
    await expect(readManifest(root)).rejects.toThrow('exceeds 1048576 bytes')
    expect(pulls).toBe(2)
  })

  test('rejects foreign-relative references that escape the GitHub owner, repository, or ref root', async () => {
    const root = await temporaryDirectory()
    await Bun.write(join(root, 'runx.yaml'), catalog('root', `
  - group: foreign
    summary: Foreign.
    runx: https://github.com/example/catalog/blob/main/child/runx.yaml`))
    globalThis.fetch = foreignGraphFetch({ nestedReference: '../../../../other/repository/main/runx.yaml' })
    await expect(readManifest(root)).rejects.toThrow('escapes its GitHub owner/repository/ref root')
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

function nestedCatalog(depth: number): string {
  let commands: unknown[] = [{
    uid: 'deep-command', id: 'leaf', summary: 'Deep command.', description: 'Deep command.', command: 'echo deep',
  }]
  for (let index = depth; index >= 1; index -= 1) {
    commands = [{ group: `g${index}`, summary: `Group ${index}.`, commands }]
  }
  return Bun.YAML.stringify({ version: '2.0.0', namespace: 'root', scripts: { directory: 'scripts' }, commands })
}

function foreignGraphFetch(options: { scriptsDirectory?: string, commandCwd?: string, nestedReference?: string }): typeof fetch {
  return (async (input) => {
    const url = String(input)
    if (url.endsWith('/child/runx.yaml')) {
      const nested = options.nestedReference
        ? [{ group: 'nested', summary: 'Nested.', runx: options.nestedReference }]
        : [{
            uid: 'foreign-command', id: 'test', summary: 'Foreign.', description: 'Foreign command.', command: 'echo foreign',
            ...(options.commandCwd ? { cwd: options.commandCwd } : {}),
          }]
      return Response.json({
        version: '2.0.0', namespace: 'foreign-child', scripts: { directory: options.scriptsDirectory ?? 'scripts' },
        parent: '../runx.yaml', commands: nested,
      }, { headers: { 'content-type': 'application/yaml' } })
    }
    return Response.json({
      version: '2.0.0', namespace: 'foreign-root', scripts: { directory: 'scripts' },
      commands: [{ group: 'child', summary: 'Child.', runx: 'child/runx.yaml' }],
    }, { headers: { 'content-type': 'application/yaml' } })
  }) as typeof fetch
}
