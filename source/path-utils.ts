/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

export {
  baseName,
  directoryName,
  homeDirectory,
  isAbsolutePath,
  joinPath,
  normalizePath,
  relativePath,
  resolvePath,
}

const windowsDrivePattern = /^[A-Za-z]:[\\/]/

function pathSeparator(value = process.platform): '/' | '\\' {
  return value === 'win32' ? '\\' : '/'
}

function normalizePath(value: string): string {
  const separator = pathSeparator()
  const normalized = value.replace(/[\\/]+/g, separator)
  const drive = windowsDrivePattern.test(normalized) ? normalized.slice(0, 2) : ''
  const rooted = drive !== '' || normalized.startsWith(separator)
  const parts = normalized.slice(drive.length).split(separator)
  const output: string[] = []

  for (const part of parts) {
    if (!part || part === '.') continue
    if (part === '..') {
      if (output.length > 0 && output.at(-1) !== '..') output.pop()
      else if (!rooted) output.push(part)
      continue
    }
    output.push(part)
  }

  const prefix = drive ? `${drive}${separator}` : rooted ? separator : ''
  return `${prefix}${output.join(separator)}` || (rooted ? prefix : '.')
}

function isAbsolutePath(value: string): boolean {
  return windowsDrivePattern.test(value) || value.startsWith('/') || value.startsWith('\\')
}

function resolvePath(...values: string[]): string {
  let current = ''
  for (const value of values) {
    if (!value) continue
    current = isAbsolutePath(value) ? value : current ? joinPath(current, value) : value
  }
  if (!isAbsolutePath(current)) current = joinPath(process.cwd(), current)
  return normalizePath(current)
}

function joinPath(...values: string[]): string {
  return normalizePath(values.filter(Boolean).join(pathSeparator()))
}

function directoryName(value: string): string {
  const normalized = normalizePath(value)
  const separator = pathSeparator()
  const index = normalized.lastIndexOf(separator)
  if (index < 0) return '.'
  if (index === 0) return separator
  if (index === 2 && windowsDrivePattern.test(normalized)) return normalized.slice(0, 3)
  return normalized.slice(0, index)
}

function baseName(value: string): string {
  const normalized = normalizePath(value)
  const index = normalized.lastIndexOf(pathSeparator())
  return index < 0 ? normalized : normalized.slice(index + 1)
}

function relativePath(from: string, to: string): string {
  const separator = pathSeparator()
  const fromParts = normalizePath(from).split(separator)
  const toParts = normalizePath(to).split(separator)
  const insensitive = process.platform === 'win32'
  let index = 0
  while (
    index < fromParts.length
    && index < toParts.length
    && (insensitive ? fromParts[index]!.toLowerCase() === toParts[index]!.toLowerCase() : fromParts[index] === toParts[index])
  ) index += 1
  return [...Array(fromParts.length - index).fill('..'), ...toParts.slice(index)].join(separator) || '.'
}

function homeDirectory(): string {
  const home = Bun.env.HOME ?? Bun.env.USERPROFILE
  if (!home) throw new Error('RunX requires HOME or USERPROFILE to resolve global storage.')
  return resolvePath(home)
}
