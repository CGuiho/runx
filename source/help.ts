import packageJson from '../package.json' with { type: 'json' }

export const readVersion = (): string => typeof packageJson.version === 'string' ? packageJson.version : '0.0.0'

export const showHome = (): string => `RunX ${readVersion()}\n\nA documented, local command catalog for runx.yaml manifests.\n\nUsage:\n  runx list [--file <path>] [--format <text|json>]\n  runx describe <selector>\n  runx run <selector> [--dry-run] [--yes]\n  runx r <selector> [--dry-run] [--yes]\n  runx <selector>\n\nStart here:\n  runx list                 List every command in the nearest manifest.\n  runx --help-tree          Show the complete command tree.\n  runx --help-docs          Show manifest and agent documentation guidance.\n`

export const showHelpTree = (): string => `runx\n|- list\n|- describe <selector>\n|- run <selector>\n|  \\`- alias: r\n|- check\n|- agents\n|  |- install <local|global>\n|  \\`- instructions\n|- upgrade [check|list]\n\\`- uninstall\n`

export const showHelpDocs = (): string => `RunX documentation\n\nManifest: runx.yaml, discovered from the current directory upward or selected with --file.\nRequired command fields: uid, id, group, summary, description, command.\nOptional command fields: cwd, shell, tags, confirm.\n\nSelectors resolve in this order: UID, group/id, one-based index, then an unambiguous ID.\nUse UID values for automation. Use runx describe <selector> and runx run <selector> --dry-run before unfamiliar execution.\n\nAgent skill: runx agents install local installs guiho-s-runx under .agents/skills.\n`

export const showCommandHelp = (command: string): string => `RunX ${command}\n\nRun \`runx --help-tree\` for command structure and \`runx --help-docs\` for manifest guidance.\n`
