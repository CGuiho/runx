/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { $ } from 'bun'
import { directoryName, homeDirectory, joinPath } from './path-utils.js'

export {
  copyTextFile,
  ensureDirectory,
  globalRunXDirectory,
  movePath,
  pathExists,
  readTextIfExists,
  removePath,
  writeTextFile,
}

function globalRunXDirectory(): string {
  return joinPath(homeDirectory(), '.guiho', 'runx')
}

async function pathExists(path: string): Promise<boolean> {
  return Bun.file(path).exists()
}

async function ensureDirectory(path: string): Promise<void> {
  await $`mkdir -p ${path}`.quiet()
}

async function removePath(path: string): Promise<void> {
  await $`rm -rf ${path}`.quiet()
}

async function movePath(from: string, to: string): Promise<void> {
  if (process.platform === 'win32') {
    const child = Bun.spawn(['cmd.exe', '/d', '/s', '/c', 'move', '/y', from, to], { stdout: 'ignore', stderr: 'pipe' })
    const [code, error] = await Promise.all([child.exited, new Response(child.stderr).text()])
    if (code !== 0) throw new Error(error.trim() || `Could not move ${from} to ${to}.`)
    return
  }
  await $`mv ${from} ${to}`.quiet()
}

async function readTextIfExists(path: string): Promise<string | null> {
  return await pathExists(path) ? Bun.file(path).text() : null
}

async function writeTextFile(path: string, value: string): Promise<void> {
  await ensureDirectory(directoryName(path))
  await Bun.write(path, value)
}

async function copyTextFile(source: string, target: string): Promise<void> {
  await writeTextFile(target, await Bun.file(source).text())
}
