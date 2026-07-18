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

process.stdout.write(`${JSON.stringify({ count: observed.length, assets: observed }, null, 2)}\n`)
