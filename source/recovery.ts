/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import type { RecoveryInstructions, UpgradeOs } from './upgrade-types.js'

export { createRecoveryInstructions }

const rawBase = 'https://raw.githubusercontent.com/CGuiho/runx/main/devops'

const createRecoveryInstructions = (targetVersion: string, os: UpgradeOs, targetSource: RecoveryInstructions['targetSource'] = 'resolved', installerBaseUrl = rawBase): RecoveryInstructions => {
  if (os === 'windows') {
    return {
      targetVersion,
      targetSource,
      installCommand: `powershell.exe -NoProfile -ExecutionPolicy Bypass -Command '& ([scriptblock]::Create((Invoke-RestMethod "${installerBaseUrl}/install.ps1"))) -Version "${escapePowerShell(targetVersion)}"'`,
      stopProcessCommand: 'powershell.exe -NoProfile -Command "Get-Process runx -ErrorAction SilentlyContinue | Stop-Process -Force"',
    }
  }
  return {
    targetVersion,
    targetSource,
    installCommand: `curl -fsSL ${installerBaseUrl}/install.sh | bash -s -- --version '${escapeShell(targetVersion)}'`,
    stopProcessCommand: 'pkill -x runx',
  }
}

const escapePowerShell = (value: string): string => value.replace(/`/g, '``').replace(/"/g, '`"')
const escapeShell = (value: string): string => value.replace(/'/g, `'"'"'`)
