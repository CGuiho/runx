import type { Static } from '@sinclair/typebox'
import type { CommandSchema, ManifestSchema } from './manifest.js'

export type RunXManifest = Static<typeof ManifestSchema>
export type RunXCommand = Static<typeof CommandSchema>

export type ResolvedCommand = RunXCommand & {
  index: number
  selector: string
  manifestPath: string
  cwd: string
}

export type OutputFormat = 'text' | 'json'

export type CliOptions = {
  cwd: string
  file?: string
  format: OutputFormat
  verbose: boolean
}

export type AgentTool = 'agents' | 'claude' | 'all'

export type AgentScope = 'local' | 'global'

export type UpdateResult = {
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
  url?: string
}

export type RunXUpgradeResult = UpdateResult & {
  executablePath: string
  scheduled: boolean
  upToDate: boolean
}
