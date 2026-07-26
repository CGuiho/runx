#!/usr/bin/env node

/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { spawn } from 'node:child_process'
import { chmod, mkdir, readFile, rename, rm, stat, writeFile } from 'node:fs/promises'
import { createHash } from 'node:crypto'
import { arch, homedir, platform } from 'node:os'
import { join } from 'node:path'

import packageJson from '../package.json' with { type: 'json' }

const os = platform() === 'win32' ? 'windows' : platform() === 'darwin' ? 'darwin' : platform() === 'linux' ? 'linux' : null
const cpu = arch() === 'x64' ? 'amd64' : arch() === 'arm64' ? 'arm64' : null
if (!os || !cpu) fail(`Unsupported platform: ${platform()} ${arch()}`)

const suffix = os === 'windows' ? '.exe' : ''
const asset = `runx-${os}-${cpu}${suffix}`
const cache = process.env.RUNX_NATIVE_CACHE || join(homedir(), '.guiho', 'runx', 'npm', packageJson.version)
const executable = join(cache, `runx${suffix}`)
const tag = encodeURIComponent(`@guiho/runx/v${packageJson.version}`)
const downloadBase = (process.env.RUNX_DOWNLOAD_BASE_URL || 'https://github.com/CGuiho/runx/releases/download').replace(/\/$/, '')
const url = `${downloadBase}/${tag}/${asset}`

try {
  await stat(executable)
} catch {
  await mkdir(cache, { recursive: true })
  const response = await fetch(url, { redirect: 'follow' })
  if (!response.ok) fail(`Could not download ${asset}: HTTP ${response.status}`)
  const temporary = `${executable}.${process.pid}.tmp`
  await writeFile(temporary, new Uint8Array(await response.arrayBuffer()), { flag: 'wx' })
  const checksumsResponse = await fetch(`${downloadBase}/${tag}/checksums.txt`, { redirect: 'follow' })
  if (!checksumsResponse.ok) fail(`Could not download checksums.txt: HTTP ${checksumsResponse.status}`)
  const checksums = await checksumsResponse.text()
  const expected = checksums.split(/\r?\n/).map((line) => line.trim().split(/\s+/)).find((fields) => fields[1] === asset)?.[0]
  if (!expected) fail(`Checksum entry missing for ${asset}`)
  const actual = createHash('sha256').update(await readFile(temporary)).digest('hex')
  if (actual.toLowerCase() !== expected.toLowerCase()) fail(`Checksum verification failed for ${asset}`)
  if (os !== 'windows') await chmod(temporary, 0o755)
  await rename(temporary, executable)
}

const child = spawn(executable, process.argv.slice(2), { env: process.env, stdio: 'inherit' })
child.once('error', (error) => fail(error.message))
child.once('exit', (code, signal) => {
  if (signal) process.kill(process.pid, signal)
  else process.exit(code ?? 1)
})

function fail(message) {
  process.stderr.write(`error: ${message}\n`)
  void rm(executable, { force: true }).catch(() => undefined)
  process.exit(1)
}
