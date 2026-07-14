#!/usr/bin/env bun

const args = process.argv.slice(2)
const binaryPath = new URL(`../vendor/runx${process.platform === 'win32' ? '.exe' : ''}`, import.meta.url)
const libraryPath = new URL('../library/guiho-runx-bin.js', import.meta.url)
const sourcePath = new URL('../source/guiho-runx-bin.ts', import.meta.url)

if (await Bun.file(binaryPath).exists()) {
  const child = Bun.spawn([Bun.fileURLToPath(binaryPath), ...args], { stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' })
  process.exit(await child.exited)
}

if (await Bun.file(libraryPath).exists()) {
  const child = Bun.spawn([process.execPath, Bun.fileURLToPath(libraryPath), ...args], { stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' })
  process.exit(await child.exited)
}

if (await Bun.file(sourcePath).exists()) {
  const child = Bun.spawn([process.execPath, Bun.fileURLToPath(sourcePath), ...args], { stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' })
  process.exit(await child.exited)
}

console.error('error: RunX executable is unavailable. Install an official native release or reinstall @guiho/runx.')
process.exit(1)
