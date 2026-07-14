import packageJson from '../package.json' with { type: 'json' }

export const readVersion = (): string => typeof packageJson.version === 'string' ? packageJson.version : '0.0.0'

export const showHome = (): string => `RunX ${readVersion()}

A documented, local command catalog for runx.yaml manifests.

Usage:
  runx list [--file <path>] [--format <text|json>]
  runx describe <selector>
  runx run <selector> [--dry-run] [--yes]
  runx r <selector> [--dry-run] [--yes]
  runx <selector>

Start here:
  runx list                 List every command in the nearest manifest.
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
  '|- agents',
  '|  |- install <local|global>',
  '|  `- instructions',
  '|- upgrade [check|list]',
  '`- uninstall',
  '',
].join('\n')

export const showHelpDocs = (): string => `RunX documentation

Manifest: runx.yaml, discovered from the current directory upward or selected with --file.
Required command fields: uid, id, group, summary, description, command.
Optional command fields: cwd, shell, tags, confirm.

Selectors resolve in this order: UID, group/id, one-based index, then an unambiguous ID.
Use UID values for automation. Use runx describe <selector> and runx run <selector> --dry-run before unfamiliar execution.

Agent skill: runx agents install local installs guiho-s-runx under .agents/skills.
`
