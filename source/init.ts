/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { RunXError } from './errors.js'
import { baseName, resolvePath } from './path-utils.js'
import { writeTextFile } from './storage.js'

import type { RunXManifest } from './types.js'

export {
  createInitialManifest,
  initializeRunXManifest,
  renderInitialManifest,
}
export type {
  RunXInitResult,
}

type RunXInitResult = { status: 'created', path: string, manifest: RunXManifest }

async function initializeRunXManifest(options: { cwd: string, config?: string }): Promise<RunXInitResult> {
  const cwd = resolvePath(options.cwd)
  const path = resolvePath(cwd, options.config ?? 'runx.yaml')
  if (await Bun.file(path).exists()) throw new RunXError(`Configuration already exists: ${path}`, 5)
  const manifest = createInitialManifest(baseName(cwd) || 'my-project', 'scripts')
  await writeTextFile(path, renderInitialManifest(manifest))
  return { status: 'created', path, manifest }
}

function createInitialManifest(namespace: string, scriptsDirectory: string): RunXManifest {
  return {
    version: '2.0.0',
    namespace,
    scripts: { directory: scriptsDirectory },
    commands: [],
    groups: {},
    children: [],
  }
}

function renderInitialManifest(manifest: RunXManifest): string {
  return [
    `version: ${JSON.stringify(manifest.version)}`,
    '',
    `namespace: ${JSON.stringify(manifest.namespace)}`,
    '',
    'scripts:',
    `  directory: ${JSON.stringify(manifest.scripts.directory)}`,
    '',
    'commands: []',
    '',
  ].join('\n')
}
