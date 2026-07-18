import { afterAll, describe, expect, test } from 'bun:test'
import { chmod, mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { join } from 'node:path'
import { createRecoveryInstructions } from '../source/recovery.js'

const root = Bun.fileURLToPath(new URL('..', import.meta.url))
const powershellInstaller = join(root, 'devops', 'install.ps1')
const shellInstaller = join(root, 'devops', 'install.sh')
const directories: string[] = []
let fixturePromise: Promise<string> | null = null

afterAll(async () => {
  await Promise.all(directories.map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('direct installer contracts', () => {
  test('platform installers pass native syntax and help smokes', async () => {
    const powershell = Bun.which('powershell.exe') ?? Bun.which('pwsh')
    if (powershell) {
      const syntax = Bun.spawn([powershell, '-NoProfile', '-Command', `[scriptblock]::Create((Get-Content -Raw -LiteralPath '${powershellInstaller.replaceAll("'", "''")}')) | Out-Null`], { stdout: 'pipe', stderr: 'pipe' })
      const [syntaxError, syntaxExit] = await Promise.all([new Response(syntax.stderr).text(), syntax.exited])
      expect(syntaxExit, syntaxError).toBe(0)
      const help = Bun.spawn([powershell, '-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', powershellInstaller, '-Help'], { stdout: 'pipe', stderr: 'pipe' })
      const [stdout, stderr, exitCode] = await Promise.all([new Response(help.stdout).text(), new Response(help.stderr).text(), help.exited])
      expect(exitCode, stderr).toBe(0)
      expect(stdout).toContain('Exact stable or prerelease version')
    }

    const shell = Bun.which('sh')
    if (shell) {
      const scriptPath = shellPath(shellInstaller, shell)
      const syntax = Bun.spawn([shell, '-n', scriptPath], { stdout: 'pipe', stderr: 'pipe' })
      const [syntaxError, syntaxExit] = await Promise.all([new Response(syntax.stderr).text(), syntax.exited])
      expect(syntaxExit, syntaxError).toBe(0)
      const help = Bun.spawn([shell, scriptPath, '--help'], { stdout: 'pipe', stderr: 'pipe' })
      const [stdout, exitCode] = await Promise.all([new Response(help.stdout).text(), help.exited])
      expect(exitCode).toBe(0)
      expect(stdout).toContain('exact stable or prerelease version')
    }
  }, 20_000)

  test('executes stable and prerelease recovery commands against controlled release assets', async () => {
    const fixture = await compiledFixture()
    for (const targetVersion of ['4.5.6', '4.6.0-alpha.3']) {
      const installDirectory = await temporaryDirectory(`runx recovery ${targetVersion}`)
      const server = Bun.serve({
        port: 0,
        idleTimeout: 120,
        fetch: async (request) => request.url.endsWith('/install.ps1')
          ? new Response(Bun.file(powershellInstaller))
          : request.url.endsWith('/install.sh')
            ? new Response(Bun.file(shellInstaller))
            : new Response(Bun.file(fixture)),
      })
      try {
        const os = process.platform === 'win32' ? 'windows' : process.platform === 'darwin' ? 'macos' : 'linux'
        const recovery = createRecoveryInstructions(targetVersion, os, 'resolved', `http://127.0.0.1:${server.port}`)
        const result = await runRecovery(recovery.installCommand, {
          RUNX_DOWNLOAD_BASE_URL: `http://127.0.0.1:${server.port}`,
          RUNX_INSTALL_DIR: installDirectory,
          RUNX_SKIP_PATH_UPDATE: '1',
          RUNX_FIXTURE_VERSION: targetVersion,
        })
        expect(result.exitCode, result.stderr).toBe(0)
        expect(result.stdout).toContain(`Installed and verified RunX ${targetVersion}`)
        expect(await installedVersion(installDirectory, targetVersion)).toBe(targetVersion)
        expect(recovery.stopProcessCommand).toMatch(/runx/)
      } finally {
        server.stop(true)
      }
    }
  }, 240_000)

  test('falls back only on 404 and preserves paths containing spaces', async () => {
    const fixture = await compiledFixture()
    const installDirectory = await temporaryDirectory('runx install with spaces')
    const requested: string[] = []
    const server = Bun.serve({
      port: 0,
      idleTimeout: 120,
      fetch: (request) => {
        requested.push(new URL(request.url).pathname)
        return requested.length === 1 ? new Response('missing', { status: 404 }) : new Response(Bun.file(fixture))
      },
    })
    try {
      const result = await runInstaller('5.0.0', installDirectory, `http://127.0.0.1:${server.port}`, '5.0.0')
      expect(result.exitCode, result.stderr).toBe(0)
      expect(requested).toHaveLength(2)
      expect(requested[0]).toContain('baseline')
      expect(await installedVersion(installDirectory, '5.0.0')).toBe('5.0.0')
    } finally {
      server.stop(true)
    }
  }, 120_000)

  test('fails corrupt downloads without trying another candidate and preserves the installation', async () => {
    const installDirectory = await temporaryDirectory('runx corrupt rollback')
    const destination = installedPath(installDirectory)
    const original = new Uint8Array([1, 9, 8, 4])
    await Bun.write(destination, original)
    const requested: string[] = []
    const server = Bun.serve({ port: 0, idleTimeout: 120, fetch: (request) => { requested.push(request.url); return new Response('<html>corrupt</html>') } })
    try {
      const result = await runInstaller('5.0.1', installDirectory, `http://127.0.0.1:${server.port}`, '5.0.1')
      expect(result.exitCode).not.toBe(0)
      expect(result.stderr).toMatch(/not a native/i)
      expect(requested).toHaveLength(1)
      expect(await Bun.file(destination).bytes()).toEqual(original)
    } finally {
      server.stop(true)
    }
  }, 120_000)

  test('rolls back an exact-version mismatch after canonical replacement', async () => {
    const fixture = await compiledFixture()
    const installDirectory = await temporaryDirectory('runx mismatch rollback')
    const destination = installedPath(installDirectory)
    const original = new Uint8Array([8, 6, 7, 5, 3, 0, 9])
    await Bun.write(destination, original)
    const server = Bun.serve({ port: 0, idleTimeout: 120, fetch: () => new Response(Bun.file(fixture)) })
    try {
      const result = await runInstaller('6.0.0', installDirectory, `http://127.0.0.1:${server.port}`, '6.0.1')
      expect(result.exitCode).not.toBe(0)
      expect(`${result.stdout}\n${result.stderr}`).toMatch(/restored/i)
      expect(await Bun.file(destination).bytes()).toEqual(original)
    } finally {
      server.stop(true)
    }
  }, 120_000)

  test('distinguishes network failure from a missing candidate and preserves the installation', async () => {
    const installDirectory = await temporaryDirectory('runx network failure')
    const destination = installedPath(installDirectory)
    const original = new Uint8Array([4, 2, 4, 2])
    await Bun.write(destination, original)
    const result = await runInstaller('6.1.0', installDirectory, 'http://127.0.0.1:1', '6.1.0')
    expect(result.exitCode).not.toBe(0)
    expect(`${result.stdout}\n${result.stderr}`).toMatch(/download failed|Failed to download/i)
    expect(`${result.stdout}\n${result.stderr}`).not.toMatch(/not available; trying/i)
    expect(await Bun.file(destination).bytes()).toEqual(original)
  }, 60_000)
})

const compiledFixture = (): Promise<string> => {
  fixturePromise ??= (async () => {
    const directory = await temporaryDirectory('runx fixture')
    const source = join(directory, 'fixture.ts')
    const output = join(directory, process.platform === 'win32' ? 'runx-fixture.exe' : 'runx-fixture')
    await Bun.write(source, `process.stdout.write(process.env.RUNX_FIXTURE_VERSION ?? '0.0.0')\n`)
    const build = Bun.spawn([process.execPath, 'build', source, '--compile', '--outfile', output], { stdout: 'pipe', stderr: 'pipe' })
    const [stderr, exitCode] = await Promise.all([new Response(build.stderr).text(), build.exited])
    if (exitCode !== 0) throw new Error(`Could not compile installer fixture: ${stderr}`)
    if (process.platform !== 'win32') await chmod(output, 0o755)
    return output
  })()
  return fixturePromise
}

const runInstaller = async (version: string, installDirectory: string, downloadBaseUrl: string, reportedVersion: string) => {
  const environment = {
    ...process.env,
    RUNX_DOWNLOAD_BASE_URL: downloadBaseUrl,
    RUNX_INSTALL_DIR: installDirectory,
    RUNX_SKIP_PATH_UPDATE: '1',
    RUNX_FIXTURE_VERSION: reportedVersion,
  }
  const command = process.platform === 'win32'
    ? ['powershell.exe', '-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', powershellInstaller, '-Version', version, '-InstallDir', installDirectory]
    : ['sh', shellInstaller, '--version', version, '--install-dir', installDirectory]
  return runProcess(command, environment)
}

const runRecovery = async (command: string, additions: Record<string, string>) => {
  const executable = process.platform === 'win32' ? 'powershell.exe' : 'sh'
  return runProcess([executable, '-c', command], { ...process.env, ...additions })
}

const runProcess = async (command: string[], env: Record<string, string | undefined>) => {
  const child = Bun.spawn(command, { env, stdout: 'pipe', stderr: 'pipe' })
  const [stdout, stderr, exitCode] = await Promise.all([new Response(child.stdout).text(), new Response(child.stderr).text(), child.exited])
  return { stdout, stderr, exitCode }
}

const installedVersion = async (directory: string, version: string): Promise<string> => {
  const child = Bun.spawn([installedPath(directory), '--version'], { env: { ...process.env, RUNX_FIXTURE_VERSION: version }, stdout: 'pipe', stderr: 'pipe' })
  const [stdout, exitCode] = await Promise.all([new Response(child.stdout).text(), child.exited])
  expect(exitCode).toBe(0)
  return stdout.trim()
}

const installedPath = (directory: string): string => join(directory, process.platform === 'win32' ? 'runx.exe' : 'runx')

const temporaryDirectory = async (prefix: string): Promise<string> => {
  const directory = await mkdtemp(join(tmpdir(), `${prefix.replaceAll(/[^a-z0-9]+/gi, '-')}-`))
  directories.push(directory)
  return directory
}

const shellPath = (path: string, shell: string): string => {
  const normalized = path.replaceAll('\\', '/')
  return shell.toLowerCase().includes('/system32/') || shell.toLowerCase().includes('\\system32\\')
    ? `/mnt/${normalized[0]?.toLowerCase()}${normalized.slice(2)}`
    : normalized
}
