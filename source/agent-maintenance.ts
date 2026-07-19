/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'
import { maintainAgentIntegration } from './agents.js'
import { RunXError } from './errors.js'
import { directoryName, resolvePath } from './path-utils.js'

export {
  agentMaintenanceWorkerCwd,
  runAgentMaintenanceWorker,
  spawnAgentMaintenanceWorker,
}

type AgentWorkerSpawnOptions = {
  readonly detached: true
  readonly stdout: 'ignore'
  readonly stderr: 'ignore'
  readonly stdin: 'ignore'
  readonly cwd: string
  readonly env: Record<string, string | undefined>
}
type AgentWorkerSpawner = (
  command: string[],
  options: AgentWorkerSpawnOptions,
) => { unref(): unknown }

const workerFlag = '--maintain-agent-integration-worker'
const workerArgsSchema = Type.Tuple([
  Type.Literal(workerFlag),
  Type.String({ minLength: 1 }),
])

function agentMaintenanceWorkerCwd(rawArgs: string[]): string | null {
  if (rawArgs[0] !== workerFlag) return null
  try {
    return resolvePath(Value.Decode(workerArgsSchema, rawArgs)[1])
  } catch (error) {
    throw new RunXError(`Invalid agent-maintenance worker invocation: ${error instanceof Error ? error.message : String(error)}`, 2)
  }
}

async function runAgentMaintenanceWorker(cwd: string): Promise<void> {
  await maintainAgentIntegration(resolvePath(cwd))
}

function spawnAgentMaintenanceWorker(
  cwd: string,
  spawn: AgentWorkerSpawner = (command, options) => Bun.spawn(command, options),
): boolean {
  if (Bun.env.RUNX_DISABLE_AGENT_MAINTENANCE_WORKER === '1') return false
  const effectiveCwd = resolvePath(cwd)
  const command = process.execPath.toLowerCase().includes('bun')
    ? [process.execPath, Bun.main, workerFlag, effectiveCwd]
    : [process.execPath, workerFlag, effectiveCwd]
  try {
    spawn(command, {
      detached: true,
      stdout: 'ignore',
      stderr: 'ignore',
      stdin: 'ignore',
      cwd: directoryName(process.execPath),
      env: { ...process.env, RUNX_DISABLE_AGENT_MAINTENANCE_WORKER: '1' },
    }).unref()
    return true
  } catch {
    return false
  }
}
