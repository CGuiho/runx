/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import type { RunXCommand, RunXManifest } from './configuration.js'

export type {
  AgentScope,
  CliOptions,
  OutputFormat,
  ResolvedCommand,
  RunXCommand,
  RunXManifest,
  UpdateResult,
}

type ResolvedCommand = RunXCommand & {
  index: number
  selector: string
  manifestPath: string
  cwd: string
}

type OutputFormat = 'text' | 'json'

type CliOptions = {
  cwd: string
  config?: string
  format: OutputFormat
  verbose: boolean
}

type AgentScope = 'local' | 'global'

type UpdateResult = {
  currentVersion: string
  latestVersion: string
  updateAvailable: boolean
  url?: string
}
