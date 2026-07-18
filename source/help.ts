import packageJson from '../package.json' with { type: 'json' }

export const readVersion = (): string => typeof packageJson.version === 'string' ? packageJson.version : '0.0.0'

export const showHome = (): string => `RunX ${readVersion()}

A documented, local command catalog for runx.yaml manifests.

Usage:
  runx init [--cwd <path>]
  runx list [--file <path>] [--format <text|json>]
  runx describe <selector>
  runx run <selector> [--dry-run] [--yes]
  runx r <selector> [--dry-run] [--yes]
  runx <selector>
  runx upgrade [--dry-run] [--format <text|json>]
  runx upgrade list [--format <text|json>]

Start here:
  runx init                 Interactively create an empty runx.yaml catalog.
  runx list                 List every command in the nearest manifest.
  runx upgrade              Install and verify the latest stable native release.
  runx upgrade list         List all stable and prerelease releases.
  runx --help-tree          Show the complete command tree.
  runx --help-docs          Show manifest and agent documentation guidance.
`

export const showHelpTree = (): string => [
  'runx',
  '|- list',
  '|- describe <selector>',
  '|- run <selector>',
  '|  `- alias: r',
  '|- check',
  '|- init',
  '|- agents',
  '|  |- install <local|global>',
  '|  `- instructions',
  '|- upgrade [check|list]',
  '`- uninstall',
  '',
].join('\n')

export const showHelpDocs = (): string => `RunX documentation

Manifest: runx.yaml, discovered from the current directory upward or selected with --file.
Create an empty catalog with runx init. Its manifest uses SemVer 1.x, configures a scripts directory, and always includes the public group.
Required command fields: uid, id, group, summary, description, command.
Optional command fields: cwd, shell, tags, confirm.

Selectors resolve in this order: UID, group/id, one-based index, then an unambiguous ID.
Use UID values for automation. Use runx describe <selector> and runx run <selector> --dry-run before unfamiliar execution.

Agent skill: runx agents install local installs guiho-s-runx under .agents/skills.

Native upgrades: runx upgrade prints its plan before download, verifies the canonical executable, rolls back failure, and always prints a pinned recovery install command.
`
