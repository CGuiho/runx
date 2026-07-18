/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { $ } from 'bun'
import { joinPath } from '../source/path-utils.js'

export {
  agentAssetNames,
  binaryTargets,
  expectedReleaseAssetNames,
}

type BinaryTarget = {
  readonly bunTarget: string
  readonly assetName: string
}

const binaryTargets: readonly BinaryTarget[] = [
  { bunTarget: 'bun-linux-arm64', assetName: 'runx-linux-arm64' },
  { bunTarget: 'bun-linux-x64', assetName: 'runx-linux-x64' },
  { bunTarget: 'bun-linux-x64-baseline', assetName: 'runx-linux-x64-baseline' },
  { bunTarget: 'bun-linux-x64-modern', assetName: 'runx-linux-x64-modern' },
  { bunTarget: 'bun-darwin-arm64', assetName: 'runx-darwin-arm64' },
  { bunTarget: 'bun-darwin-x64', assetName: 'runx-darwin-x64' },
  { bunTarget: 'bun-darwin-x64-baseline', assetName: 'runx-darwin-x64-baseline' },
  { bunTarget: 'bun-darwin-x64-modern', assetName: 'runx-darwin-x64-modern' },
  { bunTarget: 'bun-windows-arm64', assetName: 'runx-windows-arm64.exe' },
  { bunTarget: 'bun-windows-x64', assetName: 'runx-windows-x64.exe' },
  { bunTarget: 'bun-windows-x64-baseline', assetName: 'runx-windows-x64-baseline.exe' },
  { bunTarget: 'bun-windows-x64-modern', assetName: 'runx-windows-x64-modern.exe' },
]

const agentAssetNames = ['guiho-s-runx', 'guiho-i-runx'] as const
const expectedReleaseAssetNames = [...binaryTargets.map((target) => target.assetName), ...agentAssetNames]

if (expectedReleaseAssetNames.length !== 14 || new Set(expectedReleaseAssetNames).size !== 14) {
  throw new Error('RunX release matrix must contain exactly fourteen unique assets.')
}

if (import.meta.main) {
  const root = Bun.fileURLToPath(new URL('..', import.meta.url))
  const bin = joinPath(root, 'bin')
  await $`rm -rf ${bin}`.quiet()
  await $`mkdir -p ${bin}`.quiet()

  for (const target of binaryTargets) {
    const output = joinPath(bin, target.assetName)
    const child = Bun.spawn([
      process.execPath,
      'build',
      'source/guiho-runx-native-bin.ts',
      '--compile',
      '--production',
      '--minify-whitespace',
      '--minify-syntax',
      '--target',
      target.bunTarget,
      '--outfile',
      output,
    ], { cwd: root, stdout: 'pipe', stderr: 'pipe' })
    const [stdout, stderr, exitCode] = await Promise.all([
      new Response(child.stdout).text(),
      new Response(child.stderr).text(),
      child.exited,
    ])
    if (exitCode !== 0) throw new Error([`Failed to build ${target.assetName}`, stdout, stderr].filter(Boolean).join('\n'))
    if ((await Bun.file(output).size) === 0) throw new Error(`Built binary is empty: ${target.assetName}`)
    process.stdout.write(`built: ${target.assetName}\n`)
  }

  await Bun.write(joinPath(bin, 'guiho-s-runx'), Bun.file(joinPath(root, 'skills', 'guiho-s-runx', 'SKILL.md')))
  await Bun.write(joinPath(bin, 'guiho-i-runx'), Bun.file(joinPath(root, 'prompts', 'guiho-i-runx.md')))

  const observed = [...new Bun.Glob('*').scanSync({ cwd: bin, onlyFiles: true })].sort()
  const expected = [...expectedReleaseAssetNames].sort()
  if (JSON.stringify(observed) !== JSON.stringify(expected)) {
    throw new Error(`Release asset mismatch.\nExpected: ${expected.join(', ')}\nObserved: ${observed.join(', ')}`)
  }
  process.stdout.write('verified exactly 14 release assets\n')
}
