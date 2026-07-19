/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

export {
  extractReleaseNotes,
}

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function extractReleaseNotes(changelog: string, version: string): string {
  const normalized = changelog.replace(/\r\n/g, '\n')
  const heading = new RegExp(`^## ${escapeRegExp(version)} - .+$`, 'gm')
  const matches = [...normalized.matchAll(heading)]
  if (matches.length !== 1 || matches[0]?.index === undefined) {
    throw new Error(`Expected exactly one changelog heading for version ${version}.`)
  }

  const start = matches[0].index
  const nextHeading = normalized.indexOf('\n## ', start + matches[0][0].length)
  const section = normalized.slice(start, nextHeading === -1 ? normalized.length : nextHeading).trim()
  if (!section.includes('\n')) {
    throw new Error(`Changelog section for version ${version} has no release notes.`)
  }

  return `${section}\n`
}

if (import.meta.main) {
  const [version, outputPath] = Bun.argv.slice(2)
  if (!version || !outputPath) {
    throw new Error('Usage: bun run devops/extract-release-notes.ts <version> <output-path>')
  }

  const changelog = await Bun.file(new URL('../CHANGELOG.md', import.meta.url)).text()
  await Bun.write(outputPath, extractReleaseNotes(changelog, version))
  process.stdout.write(`wrote release notes for ${version}: ${outputPath}\n`)
}
