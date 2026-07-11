const targets = [
  ['bun-linux-x64', 'runx-linux-x64'],
  ['bun-linux-arm64', 'runx-linux-arm64'],
  ['bun-darwin-x64', 'runx-darwin-x64'],
  ['bun-darwin-arm64', 'runx-darwin-arm64'],
  ['bun-windows-x64', 'runx-windows-x64.exe'],
  ['bun-windows-arm64', 'runx-windows-arm64.exe'],
] as const

for (const [target, output] of targets) {
  const child = Bun.spawn([
    process.execPath,
    'build',
    'source/guiho-runx-native-bin.ts',
    '--compile',
    `--target=${target}`,
    '--minify',
    '--outfile',
    `bin/${output}`,
  ], { stdin: 'inherit', stdout: 'inherit', stderr: 'inherit' })
  const code = await child.exited
  if (code !== 0) process.exit(code)
}
