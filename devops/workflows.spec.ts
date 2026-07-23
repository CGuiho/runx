/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

import { describe, expect, test } from 'bun:test'

describe('RunX release workflow installer ownership', () => {
  test('ordinary CI smokes the latest public release without comparing it to source', async () => {
    const workflow = await Bun.file(new URL('../.github/workflows/ci.yml', import.meta.url)).text()
    const installerStep = workflow.slice(
      workflow.indexOf('- name: Verify public Linux installer'),
      workflow.indexOf('\n\n  test-windows:'),
    )

    expect(installerStep).toContain('runx/main/devops/install.sh | bash')
    expect(installerStep).toContain('installed_version=')
    expect(installerStep).not.toContain('package.json')
    expect(installerStep).not.toContain('expected_version')
  })

  test('publish accepts the exact tagged release only after its assets exist', async () => {
    const workflow = await Bun.file(new URL('../.github/workflows/publish.yml', import.meta.url)).text()
    const assets = workflow.indexOf('- name: Verify exact GitHub release asset set')
    const installer = workflow.indexOf('- name: Verify exact-version public Linux installer')
    const npm = workflow.indexOf('- name: Set up npm trusted publishing')
    const installerStep = workflow.slice(installer, npm)

    expect(assets).toBeGreaterThan(-1)
    expect(installer).toBeGreaterThan(assets)
    expect(npm).toBeGreaterThan(installer)
    expect(installerStep).toContain('${GITHUB_REF_NAME}/devops/install.sh')
    expect(installerStep).toContain('--version "$VERSION"')
    expect(installerStep).toContain('runx" --version)" = "$VERSION"')
  })
})
