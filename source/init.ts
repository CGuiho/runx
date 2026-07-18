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

function createInitialManifest(projectName: string, scriptsDirectory: string): RunXManifest {
  return {
    version: '1.0.0',
    project: { name: projectName },
    scripts: { directory: scriptsDirectory },
    groups: { public: { summary: 'Default public project commands.' } },
    commands: [],
  }
}

function renderInitialManifest(manifest: RunXManifest): string {
  return [
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
    `    summary: ${JSON.stringify(manifest.groups.public!.summary)}`,
    '',
    'commands: []',
    '',
  ].join('\n')
}
