/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

export {
  renderWelcome,
}

type WelcomeOptions = {
  readonly architecture?: string
  readonly platform?: string
  readonly updateNotice?: string | null
  readonly version: string
}

const platformLabels: Readonly<Record<string, string>> = {
  darwin: 'macOS',
  linux: 'Linux',
  win32: 'Windows',
}

const welcomeWidth = 58

function renderWelcome(options: WelcomeOptions): string {
  const platform = platformLabels[options.platform ?? process.platform] ?? (options.platform ?? process.platform)
  const architecture = options.architecture ?? process.arch
  const contentWidth = welcomeWidth - 4
  const bordered = [
    `╔${'═'.repeat(welcomeWidth - 2)}╗`,
    borderLine('RUNX', contentWidth),
    borderLine('Documented command catalog', contentWidth),
    borderLine('GUIHO · Cristóvão GUIHO', contentWidth),
    `╚${'═'.repeat(welcomeWidth - 2)}╝`,
  ]
  const lines = [
    ...bordered,
    '',
    `  platform      ${platform} ${architecture}`,
    `  version       v${options.version}`,
    '',
    '  Run `runx --help` to see available commands.',
  ]
  if (options.updateNotice) lines.push('', ...options.updateNotice.split('\n'))
  return `${lines.join('\n')}\n`
}

function borderLine(value: string, width: number): string {
  return `║  ${value.padEnd(width)}║`
}
