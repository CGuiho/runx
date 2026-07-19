/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { expectedReleaseAssetNames } from './build-binaries.js'
import { joinPath } from '../source/path-utils.js'

const root = Bun.fileURLToPath(new URL('..', import.meta.url))
const bin = joinPath(root, 'bin')
const observed = [...new Bun.Glob('*').scanSync({ cwd: bin, onlyFiles: true })].sort()
const expected = [...expectedReleaseAssetNames].sort()

if (JSON.stringify(observed) !== JSON.stringify(expected)) {
  throw new Error(`Expected exactly ${expected.length} release assets.\nExpected: ${expected.join(', ')}\nObserved: ${observed.join(', ')}`)
}

for (const [assetName, expectedName] of [
  ['guiho-s-runx.md', 'guiho-s-runx'],
  ['guiho-i-runx.md', 'guiho-i-runx'],
] as const) {
  const bytes = await Bun.file(joinPath(bin, assetName)).bytes()
  if (bytes.length === 0 || (bytes[0] === 0x4d && bytes[1] === 0x5a) || bytes.includes(0)) {
    throw new Error(`Release Markdown asset is empty or binary: ${assetName}`)
  }
  const text = new TextDecoder('utf-8', { fatal: true }).decode(bytes).replace(/\r\n/g, '\n')
  if (!text.startsWith('---\n') || !new RegExp(`^name:\\s*${expectedName}\\s*$`, 'm').test(text)) {
    throw new Error(`Release Markdown asset has invalid frontmatter identity: ${assetName}`)
  }
}

process.stdout.write(`${JSON.stringify({ count: observed.length, assets: observed }, null, 2)}\n`)
