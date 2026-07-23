/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'
import { $ } from 'bun'

describe('RunX direct installers', () => {
  test('README publishes the exact simplified POSIX bootstrap', async () => {
    const readme = await Bun.file(new URL('../README.md', import.meta.url)).text()
    const command = 'curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash'
    expect(readme).toContain(command)
    expect(readme).not.toContain("curl --proto '=https' --tlsv1.2")
  })

  test('PowerShell installer exposes the complete RFC sequence', async () => {
    const script = await Bun.file(new URL('./install.ps1', import.meta.url)).text()
    expect(script).toContain('Initiating GUIHO CLI Upgrade / Installation Sequence...')
    expect(script).toContain('guiho-s-runx.md')
    expect(script).toContain('guiho-i-runx.md')
    expect(script).toContain('.agents\\skills\\guiho-s-runx')
    expect(script).toContain('.claude\\skills\\guiho-s-runx')
    expect(script).toContain('[char]0x2014')
    expect(script).toContain('Test-NativeBinary')
    expect(script).toContain('Install-Transactional')
    expect(script).toContain('function Test-MarkdownAsset')
    expect(script).toContain("Test-MarkdownAsset -Path $skillAsset -ExpectedName 'guiho-s-runx'")
    expect(script).toContain("Test-MarkdownAsset -Path $promptAsset -ExpectedName 'guiho-i-runx'")
    expect(script).toContain('Windows executable header')
    expect(script).toContain('binary NUL bytes')
    expect(script).toContain('not valid UTF-8 text')
    expect(script).toContain('function Read-Utf8Text')
    expect(script).toContain('[System.IO.File]::WriteAllText')
    expect(script).toContain("RUNX_DISABLE_UPDATE_WORKER'] = '1'")
    expect(script).toContain("RUNX_DISABLE_AGENT_MAINTENANCE_WORKER'] = '1'")
  })

  test('Bash installer selects Darwin assets and resolves latest assets without parsing redirect tags', async () => {
    const script = await Bun.file(new URL('./install.sh', import.meta.url)).text()
    expect(script).toStartWith('#!/usr/bin/env bash\nset -euo pipefail')
    expect(script).toContain("Darwin) printf 'darwin")
    expect(script).toContain('--progress-bar')
    expect(script).toContain('/releases/latest/download/')
    expect(script).not.toContain("write-out '%{url_effective}'")
    expect(script).toContain('guiho-s-runx.md')
    expect(script).toContain('guiho-i-runx.md')
    expect(script).toContain('.agents/skills/guiho-s-runx')
    expect(script).toContain('.claude/skills/guiho-s-runx')
    expect(script).toContain('verify_markdown_asset()')
    expect(script).toContain("verify_markdown_asset \"$TMP/guiho-s-runx.md\" 'guiho-s-runx'")
    expect(script).toContain("verify_markdown_asset \"$TMP/guiho-i-runx.md\" 'guiho-i-runx'")
    expect(script).toContain('Windows executable header')
    expect(script).toContain('binary NUL bytes')
    expect(script).toContain('iconv -f UTF-8 -t UTF-8')
    expect(script).toContain('BEGIN RUNX — DO NOT EDIT THIS SECTION')
    expect(script).not.toContain(' bun ')
  })

  test('Bash installer passes syntax, piped startup, exact versions, and executable verification', async () => {
    const bash = await bashExecutable()
    const script = Bun.fileURLToPath(new URL('./install.sh', import.meta.url))
    const syntax = Bun.spawnSync([bash, '-n', script])
    expect(syntax.exitCode).toBe(0)

    const piped = Bun.spawnSync([
      bash,
      '-c',
      'cat "$1" | RUNX_INSTALLER_SOURCE_ONLY=1 bash -s --',
      'bash',
      script,
    ])
    expect(piped.exitCode).toBe(0)
    expect(piped.stderr.toString()).toBe('')

    const functions = Bun.spawnSync([
      bash,
      '-c',
      [
        'RUNX_INSTALLER_SOURCE_ONLY=1 source "$1"',
        'stable=$(normalize_version v1.2.3)',
        'prerelease=$(normalize_version @guiho/runx@1.3.0-alpha.2)',
        'test "$stable" = 1.2.3',
        'test "$prerelease" = 1.3.0-alpha.2',
        'VERSION=latest',
        'TARGET_VERSION=$(resolve_target_version)',
        'latest=$(build_asset_url runx-linux-x64-baseline)',
        'test "$latest" = https://github.com/CGuiho/runx/releases/latest/download/runx-linux-x64-baseline',
        'VERSION=@guiho/runx@1.3.0-alpha.2',
        'TARGET_VERSION=$(resolve_target_version)',
        'exact=$(build_asset_url runx-linux-x64-baseline)',
        'test "$exact" = https://github.com/CGuiho/runx/releases/download/%40guiho%2Frunx%401.3.0-alpha.2/runx-linux-x64-baseline',
        'TARGET_VERSION=$(bun --version)',
        'verify_installed_version "$(command -v bun)"',
      ].join('; '),
      'bash',
      script,
    ])
    expect(functions.exitCode).toBe(0)
  })

  test('platform installer rejects a PE payload disguised as Markdown', async () => {
    const temporaryDirectory = `${Bun.env.TEMP ?? Bun.env.TMPDIR ?? '/tmp'}/runx-markdown-${crypto.randomUUID()}`
    const validAsset = `${temporaryDirectory}/valid.md`
    const executableAsset = `${temporaryDirectory}/payload.md`
    await $`mkdir -p ${temporaryDirectory}`
    try {
      await Bun.write(validAsset, '---\nname: guiho-s-runx\n---\n\n# Skill\n')
      await Bun.write(executableAsset, new Uint8Array([0x4d, 0x5a, 0x00, 0x00]))

      if (process.platform === 'win32') {
        const script = Bun.fileURLToPath(new URL('./install.ps1', import.meta.url))
        const valid = Bun.spawnSync([
          'powershell',
          '-NoLogo',
          '-NoProfile',
          '-NonInteractive',
          '-ExecutionPolicy',
          'Bypass',
          '-Command',
          `$env:RUNX_INSTALLER_SOURCE_ONLY='1'; . '${script.replaceAll("'", "''")}'; Test-MarkdownAsset -Path '${validAsset.replaceAll("'", "''")}' -ExpectedName 'guiho-s-runx'`,
        ])
        const invalid = Bun.spawnSync([
          'powershell',
          '-NoLogo',
          '-NoProfile',
          '-NonInteractive',
          '-ExecutionPolicy',
          'Bypass',
          '-Command',
          `$env:RUNX_INSTALLER_SOURCE_ONLY='1'; . '${script.replaceAll("'", "''")}'; Test-MarkdownAsset -Path '${executableAsset.replaceAll("'", "''")}' -ExpectedName 'guiho-s-runx'`,
        ])
        expect(valid.exitCode).toBe(0)
        expect(invalid.exitCode).not.toBe(0)
        expect(invalid.stderr.toString()).toContain('Windows executable header')
      } else {
        const script = Bun.fileURLToPath(new URL('./install.sh', import.meta.url))
        const bash = await bashExecutable()
        const valid = Bun.spawnSync([
          bash,
          '-c',
          `RUNX_INSTALLER_SOURCE_ONLY=1 . "$1"; verify_markdown_asset "$2" guiho-s-runx`,
          'bash',
          script,
          validAsset,
        ])
        const invalid = Bun.spawnSync([
          bash,
          '-c',
          `RUNX_INSTALLER_SOURCE_ONLY=1 . "$1"; verify_markdown_asset "$2" guiho-s-runx`,
          'bash',
          script,
          executableAsset,
        ])
        expect(valid.exitCode).toBe(0)
        expect(invalid.exitCode).not.toBe(0)
        expect(invalid.stderr.toString()).toContain('Windows executable header')
      }
    } finally {
      await $`rm -rf ${temporaryDirectory}`
    }
  })

  if (process.platform === 'win32') {
    test('PowerShell reconciliation preserves UTF-8 and converges duplicate markers idempotently', async () => {
      const directory = `${Bun.env.TEMP ?? 'C:/tmp'}/runx-instructions-${crypto.randomUUID()}`
      const instruction = `${directory}/AGENTS.md`
      const prompt = `${directory}/prompt.md`
      await $`mkdir -p ${directory}`
      try {
        await Bun.write(instruction, `# Café guidance — preserve exactly

<!-- BEGIN RUNX \u00e2\u20ac\u201d DO NOT EDIT THIS SECTION -->
corrupted block
<!-- END RUNX -->

<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->
duplicate block
<!-- END RUNX -->
`)
        await Bun.write(prompt, '## RunX instructions\n\nPreserve naïve Unicode — exactly.\n')
        const script = Bun.fileURLToPath(new URL('./install.ps1', import.meta.url)).replaceAll("'", "''")
        const command = [
          "$env:RUNX_INSTALLER_SOURCE_ONLY='1'",
          `. '${script}'`,
          `$instruction = '${instruction.replaceAll("'", "''")}'`,
          `$prompt = Read-Utf8Text -Path '${prompt.replaceAll("'", "''")}'`,
          'Reconcile-InstructionFile -Path $instruction -Prompt $prompt',
          '$first = [Convert]::ToBase64String([IO.File]::ReadAllBytes($instruction))',
          'Reconcile-InstructionFile -Path $instruction -Prompt $prompt',
          '$second = [Convert]::ToBase64String([IO.File]::ReadAllBytes($instruction))',
          "if ($first -cne $second) { throw 'second reconciliation changed bytes' }",
          '$text = Read-Utf8Text -Path $instruction',
          "if (-not $text.Contains('# Café guidance — preserve exactly')) { throw 'Unicode guidance changed' }",
          "if (-not $text.Contains('Preserve naïve Unicode — exactly.')) { throw 'Unicode prompt changed' }",
          "if (([regex]::Matches($text, '<!-- BEGIN RUNX')).Count -ne 1) { throw 'managed marker did not converge' }",
          "if ($text.Contains('corrupted block') -or $text.Contains('duplicate block')) { throw 'stale managed content remained' }",
          '$bytes = [IO.File]::ReadAllBytes($instruction)',
          "if ($bytes.Length -ge 3 -and $bytes[0] -eq 0xEF -and $bytes[1] -eq 0xBB -and $bytes[2] -eq 0xBF) { throw 'unexpected UTF-8 BOM' }",
        ].join('; ')
        const result = Bun.spawnSync([
          'powershell',
          '-NoLogo',
          '-NoProfile',
          '-NonInteractive',
          '-ExecutionPolicy',
          'Bypass',
          '-Command',
          command,
        ])
        expect(result.exitCode).toBe(0)
        expect(result.stderr.toString()).toBe('')
      } finally {
        await $`rm -rf ${directory}`
      }
    })

    test('PowerShell installer accepts exact stable and prerelease versions and verifies the executable', () => {
      const script = Bun.fileURLToPath(new URL('./install.ps1', import.meta.url)).replaceAll("'", "''")
      const bun = process.execPath.replaceAll("'", "''")
      const command = [
        "$env:RUNX_INSTALLER_SOURCE_ONLY='1'",
        `. '${script}'`,
        "$stable = Resolve-TargetVersion -RequestedVersion 'v1.2.3'",
        "$prerelease = Resolve-TargetVersion -RequestedVersion '@guiho/runx@1.3.0-alpha.2'",
        "if ($stable -ne '1.2.3' -or $prerelease -ne '1.3.0-alpha.2') { throw 'exact version normalization failed' }",
        `Test-InstalledVersion -Path '${bun}' -ExpectedVersion '${process.versions.bun}'`,
        "$mismatchRejected = $false",
        `try { Test-InstalledVersion -Path '${bun}' -ExpectedVersion '0.0.0-impossible' } catch { $mismatchRejected = $true }`,
        "if (-not $mismatchRejected) { throw 'version mismatch was not rejected' }",
      ].join('; ')
      const result = Bun.spawnSync([
        'powershell',
        '-NoLogo',
        '-NoProfile',
        '-NonInteractive',
        '-ExecutionPolicy',
        'Bypass',
        '-Command',
        command,
      ])
      expect(result.exitCode).toBe(0)
    }, 15_000)
  }
})

async function bashExecutable(): Promise<string> {
  const gitBash = 'C:\\Program Files\\Git\\bin\\bash.exe'
  if (process.platform === 'win32' && await Bun.file(gitBash).exists()) return gitBash
  const bash = Bun.which('bash')
  if (!bash) throw new Error('Bash is required to validate devops/install.sh.')
  return bash
}
