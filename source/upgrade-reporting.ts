/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import type { RecoveryInstructions, ReleaseCatalog, UpgradeEnvelope, UpgradeEvent, UpgradePlan } from './upgrade-types.js'

export {
  renderReleaseCatalog,
  renderUpgradeEvent,
  renderUpgradeHeading,
  renderUpgradePlan,
  renderUpgradeResult,
}

const renderUpgradeHeading = (): string => [
  '------------------------------------------------------------',
  '  Upgrading the CLI',
  '------------------------------------------------------------',
].join('\n') + '\n'

const renderUpgradePlan = (plan: UpgradePlan): string => [
  `  current : ${plan.currentVersion}`,
  `  target  : ${plan.targetVersion}`,
  `  os      : ${plan.os}`,
  `  arch    : ${plan.arch}`,
  `  binary  : ${plan.assetName}`,
  `  path    : ${plan.executablePath.replaceAll('\\', '/')}`,
  `  url     : ${plan.assetUrl}`,
  '------------------------------------------------------------',
].join('\n') + '\n'

const renderUpgradeEvent = (event: UpgradeEvent): string => {
  if (event.status !== 'started' || event.phase === 'plan') return ''
  const labels: Partial<Record<UpgradeEvent['phase'], string>> = {
    download: 'Downloading...',
    validate: 'Validating...',
    replace: 'Replacing...',
    verify: 'Verifying...',
    cleanup: 'Cleaning up...',
  }
  const label = labels[event.phase]
  return label ? `${label}\n` : ''
}

const renderUpgradeResult = (result: UpgradeEnvelope): string => {
  const summary = result.outcome === 'upgraded'
    ? `Upgrade complete: ${result.plan?.currentVersion} -> ${result.plan?.targetVersion}`
    : result.outcome === 'up-to-date'
      ? `Already up to date: ${result.result?.installedVersion ?? result.recovery.targetVersion}`
      : result.outcome === 'dry-run'
        ? `Dry run complete: ${result.plan?.currentVersion} -> ${result.plan?.targetVersion}`
        : result.outcome === 'rolled-back'
          ? `Upgrade failed during ${result.error?.phase ?? 'replace'}; restored RunX ${result.result?.installedVersion ?? result.recovery.targetVersion}.`
          : `Upgrade failed during ${result.error?.phase ?? 'plan'}.`
  return `${summary}\n\n${renderRecovery(result.recovery, result.outcome === 'failed' || result.outcome === 'rolled-back')}\n`
}

const renderRecovery = (recovery: RecoveryInstructions, repair: boolean): string => [
  recovery.targetSource === 'fallback-current'
    ? `Repair reinstall (target lookup failed; pinned to installed RunX ${recovery.targetVersion}):`
    : `${repair ? 'To repair RunX' : 'If the new version is not active'}, install RunX ${recovery.targetVersion} directly:`,
  `  ${recovery.installCommand}`,
  '',
  'If RunX is still running and blocks installation, stop it first:',
  `  ${recovery.stopProcessCommand}`,
].join('\n')

const renderReleaseCatalog = (catalog: ReleaseCatalog): string => {
  const headers = ['VERSION', 'CHANNEL', 'PUBLISHED', 'CURRENT', 'LATEST', 'ASSET'] as const
  const rows = catalog.releases.map((release) => [
    release.version,
    release.channel,
    release.publishedAt?.slice(0, 10) ?? '-',
    release.current ? 'yes' : '',
    release.latestStable ? 'yes' : '',
    release.compatibleAsset ? 'yes' : 'no',
  ])
  const widths = headers.map((header, index) => Math.max(header.length, ...rows.map((row) => row[index]?.length ?? 0)))
  const line = (values: readonly string[]): string => values.map((value, index) => value.padEnd(widths[index] ?? value.length)).join('  ').trimEnd()
  return ['AVAILABLE RUNX VERSIONS', '', line(headers), ...rows.map(line)].join('\n') + '\n'
}
