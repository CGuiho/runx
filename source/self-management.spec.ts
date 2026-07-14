import { afterEach, describe, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { upgradeSelf } from './self-management.js'

const directories: string[] = []
const originalFetch = globalThis.fetch
const originalSelfPath = process.env['RUNX_SELF_PATH']

afterEach(async () => {
  globalThis.fetch = originalFetch
  if (originalSelfPath === undefined) delete process.env['RUNX_SELF_PATH']
  else process.env['RUNX_SELF_PATH'] = originalSelfPath
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX self-management', () => {
  if (process.platform === 'win32') {
    test('replaces a running Windows executable before reporting success', async () => {
      const directory = await temporaryDirectory()
      const executablePath = join(directory, 'runx.exe')
      const powershellPath = Bun.which('powershell.exe')
      if (!powershellPath) throw new Error('powershell.exe is required for the Windows replacement test')
      await Bun.write(executablePath, Bun.file(powershellPath))
      const runningExecutable = Bun.spawn([executablePath, '-NoProfile', '-Command', 'Start-Sleep -Seconds 30'], {
        stdin: 'ignore',
        stdout: 'ignore',
        stderr: 'ignore',
      })

      try {
        await Bun.sleep(250)
        const targetVersion = process.versions.bun
        process.env['RUNX_SELF_PATH'] = executablePath
        mockUpgrade(targetVersion, Bun.file(process.execPath))

        const result = await upgradeSelf(false)

        expect(result.latestVersion).toBe(targetVersion)
        expect(result.scheduled).toBe(false)
        expect(result.upToDate).toBe(false)
        expect(await Bun.file(`${executablePath}.new`).exists()).toBe(false)
        expect(await executableVersion(executablePath)).toBe(targetVersion)
      } finally {
        runningExecutable.kill()
        await runningExecutable.exited
        await waitForRemoval(`${executablePath}.old`)
      }
    }, 15_000)

    test('restores the previous Windows executable when verification fails', async () => {
      const directory = await temporaryDirectory()
      const executablePath = join(directory, 'runx.exe')
      const originalBytes = new Uint8Array([0x4d, 0x5a, 0x01, 0x02])
      await Bun.write(executablePath, originalBytes)
      process.env['RUNX_SELF_PATH'] = executablePath
      mockUpgrade('99.0.0', new Blob(['not-an-executable']))

      await expect(upgradeSelf(false)).rejects.toThrow('the previous executable was restored')
      expect(await Bun.file(executablePath).bytes()).toEqual(originalBytes)
      expect(await Bun.file(`${executablePath}.new`).exists()).toBe(false)
      expect(await Bun.file(`${executablePath}.old`).exists()).toBe(false)
    })
  }
})

const temporaryDirectory = async (): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), 'runx-self-management-'))
  directories.push(directory)
  return directory
}

const mockUpgrade = (version: string, binary: Blob): void => {
  let request = 0
  globalThis.fetch = (async () => {
    request += 1
    if (request === 1) {
      return Response.json({
        tag_name: `@guiho/runx@${version}`,
        assets: [{ name: `runx-windows-${process.arch}.exe`, browser_download_url: 'https://example.test/runx.exe' }],
      })
    }
    return new Response(binary)
  }) as typeof fetch
}

const executableVersion = async (executablePath: string): Promise<string> => {
  const process = Bun.spawn([executablePath, '--version'], { stdin: 'ignore', stdout: 'pipe', stderr: 'pipe' })
  const [stdout, exitCode] = await Promise.all([new Response(process.stdout).text(), process.exited])
  if (exitCode !== 0) throw new Error(`Replacement verification exited with code ${exitCode}`)
  return stdout.trim()
}

const waitForRemoval = async (path: string): Promise<void> => {
  for (let attempt = 0; attempt < 100; attempt += 1) {
    if (!await Bun.file(path).exists()) return
    await Bun.sleep(100)
  }
  throw new Error(`Timed out waiting for cleanup: ${path}`)
}
