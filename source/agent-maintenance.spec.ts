/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { afterEach, describe, expect, test } from 'bun:test'
import { mkdtemp, rm } from 'node:fs/promises'
import { tmpdir } from 'node:os'
import { dirname, join } from 'node:path'
import { agentMaintenanceWorkerCwd, spawnAgentMaintenanceWorker } from './agent-maintenance.js'
import { maintainAgentIntegration } from './agents.js'

const directories: string[] = []

afterEach(async () => {
  await Promise.all(directories.splice(0).map((directory) => rm(directory, { recursive: true, force: true })))
})

describe('RunX automatic agent maintenance', () => {
  test.serial('installs missing resources, refreshes stale resources, and leaves current resources unchanged', async () => {
    const home = await temporaryDirectory()
    const project = await temporaryDirectory()
    const nested = join(project, 'packages', 'cli')
    await Bun.write(join(nested, '.keep'), '')
    await Bun.write(join(project, 'AGENTS.md'), '# Existing agent guidance\n')

    await withHome(home, async () => {
      const installed = await maintainAgentIntegration(nested)
      expect(installed.skills).toHaveLength(2)
      expect(installed.instructions).toEqual([join(project, 'AGENTS.md')])

      const skill = await Bun.file(new URL('../skills/guiho-s-runx/SKILL.md', import.meta.url)).text()
      for (const tool of ['.agents', '.claude']) {
        expect(await Bun.file(join(home, tool, 'skills', 'guiho-s-runx', 'SKILL.md')).text()).toBe(skill)
      }
      const agents = await Bun.file(join(project, 'AGENTS.md')).text()
      expect(agents).toContain('# Existing agent guidance')
      expect(agents.match(/BEGIN RUNX — DO NOT EDIT THIS SECTION/g)).toHaveLength(1)
      expect(agents).toContain('creating or\nupdating catalog entries')

      expect(await maintainAgentIntegration(nested)).toEqual({ skills: [], instructions: [] })

      await Bun.write(join(home, '.agents', 'skills', 'guiho-s-runx', 'SKILL.md'), 'legacy skill')
      await Bun.write(join(project, 'AGENTS.md'), `# Existing agent guidance

<!-- BEGIN RUNX AGENT INSTRUCTIONS -->
outdated block
<!-- END RUNX AGENT INSTRUCTIONS -->

User-authored ending.
`)
      const refreshed = await maintainAgentIntegration(nested)
      expect(refreshed.skills).toEqual([join(home, '.agents', 'skills', 'guiho-s-runx', 'SKILL.md')])
      expect(refreshed.instructions).toEqual([join(project, 'AGENTS.md')])
      const refreshedAgents = await Bun.file(join(project, 'AGENTS.md')).text()
      expect(refreshedAgents).toContain('# Existing agent guidance')
      expect(refreshedAgents).toContain('User-authored ending.')
      expect(refreshedAgents).not.toContain('BEGIN RUNX AGENT INSTRUCTIONS')
      expect(refreshedAgents.match(/BEGIN RUNX — DO NOT EDIT THIS SECTION/g)).toHaveLength(1)
    })
  })

  test.serial('concurrent workers converge on one instruction block and identical skills', async () => {
    const home = await temporaryDirectory()
    const project = await temporaryDirectory()
    await Bun.write(join(project, 'AGENTS.md'), '# Concurrent project\n')

    await withHome(home, async () => {
      await Promise.all([
        maintainAgentIntegration(project),
        maintainAgentIntegration(project),
      ])

      const agents = await Bun.file(join(project, 'AGENTS.md')).text()
      expect(agents.match(/BEGIN RUNX — DO NOT EDIT THIS SECTION/g)).toHaveLength(1)
      const standard = await Bun.file(join(home, '.agents', 'skills', 'guiho-s-runx', 'SKILL.md')).text()
      const claude = await Bun.file(join(home, '.claude', 'skills', 'guiho-s-runx', 'SKILL.md')).text()
      expect(standard).toBe(claude)
    })
  })

  test('validates hidden worker input and isolates spawn failures', async () => {
    const cwd = await temporaryDirectory()
    expect(agentMaintenanceWorkerCwd(['list'])).toBeNull()
    expect(agentMaintenanceWorkerCwd(['--maintain-agent-integration-worker', cwd])).toBe(cwd)
    expect(() => agentMaintenanceWorkerCwd(['--maintain-agent-integration-worker'])).toThrow(
      'Invalid agent-maintenance worker invocation',
    )

    const originalDisabled = Bun.env.RUNX_DISABLE_AGENT_MAINTENANCE_WORKER
    delete Bun.env.RUNX_DISABLE_AGENT_MAINTENANCE_WORKER
    try {
      let detached = false
      const spawned = spawnAgentMaintenanceWorker(cwd, (command, options) => {
        expect(command).toContain('--maintain-agent-integration-worker')
        expect(command.at(-1)).toBe(cwd)
        expect(options.stdout).toBe('ignore')
        expect(options.stderr).toBe('ignore')
        expect(options.cwd).toBe(dirname(process.execPath))
        return { unref: () => { detached = true } }
      })
      expect(spawned).toBe(true)
      expect(detached).toBe(true)
      expect(spawnAgentMaintenanceWorker(cwd, () => { throw new Error('spawn failure') })).toBe(false)
    } finally {
      restoreEnvironment('RUNX_DISABLE_AGENT_MAINTENANCE_WORKER', originalDisabled)
    }
  })
})

async function temporaryDirectory(): Promise<string> {
  const directory = await mkdtemp(join(tmpdir(), 'runx-agent-maintenance-'))
  directories.push(directory)
  return directory
}

async function withHome(home: string, action: () => Promise<void>): Promise<void> {
  const originalHome = Bun.env.HOME
  const originalProfile = Bun.env.USERPROFILE
  try {
    Bun.env.HOME = home
    Bun.env.USERPROFILE = home
    await action()
  } finally {
    restoreEnvironment('HOME', originalHome)
    restoreEnvironment('USERPROFILE', originalProfile)
  }
}

function restoreEnvironment(name: string, value: string | undefined): void {
  if (value === undefined) delete Bun.env[name]
  else Bun.env[name] = value
}
