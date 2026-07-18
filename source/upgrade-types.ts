/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

export type UpgradeOs = 'windows' | 'darwin' | 'linux'
export type UpgradeArch = 'x64' | 'arm64'
export type UpgradeVariant = 'baseline' | 'default' | 'modern'

export type ReleaseAsset = {
  name: string
  url: string
}

export type ReleaseCatalogEntry = {
  tag: string
  version: string
  channel: string
  prerelease: boolean
  publishedAt: string | null
  current: boolean
  latestStable: boolean
  compatibleAsset: ReleaseAsset | null
}

export type ReleaseCatalog = {
  schemaVersion: 1
  command: 'runx upgrade list'
  currentVersion: string
  latestStableVersion: string | null
  releases: ReleaseCatalogEntry[]
}

export type RecoveryInstructions = {
  targetVersion: string
  targetSource: 'resolved' | 'fallback-current'
  installCommand: string
  stopProcessCommand: string
}

export type UpgradePlan = {
  currentVersion: string
  targetVersion: string
  os: UpgradeOs
  arch: UpgradeArch
  assetName: string
  assetUrl: string
  executablePath: string
}

export type UpgradePhase = 'plan' | 'download' | 'validate' | 'replace' | 'verify' | 'cache' | 'cleanup'
export type UpgradePhaseStatus = 'started' | 'succeeded' | 'skipped' | 'failed'

export type UpgradeEvent = {
  sequence: number
  phase: UpgradePhase
  status: UpgradePhaseStatus
  message?: string
}

export type UpgradeOutcome = 'upgraded' | 'up-to-date' | 'dry-run' | 'rolled-back' | 'failed'

export type UpgradeResult = {
  installedVersion: string
  cleanupDeferred: boolean
}

export type UpgradeErrorCode =
  | 'release_lookup_failed'
  | 'release_payload_invalid'
  | 'no_compatible_asset'
  | 'download_failed'
  | 'download_invalid'
  | 'backup_failed'
  | 'replace_failed'
  | 'verification_failed'
  | 'rollback_failed'
  | 'installer_verification_failed'

export type UpgradeError = {
  code: UpgradeErrorCode
  phase: UpgradePhase
  message: string
}

export type UpgradeEnvelope = {
  schemaVersion: 1
  command: 'runx upgrade'
  outcome: UpgradeOutcome
  plan: UpgradePlan | null
  events: UpgradeEvent[]
  result: UpgradeResult | null
  recovery: RecoveryInstructions
  error: UpgradeError | null
}
