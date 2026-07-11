export class RunXError extends Error {
  readonly exitCode: number

  constructor(message: string, exitCode = 1) {
    super(message)
    this.name = 'RunXError'
    this.exitCode = exitCode
  }
}

export const invariant: (condition: unknown, message: string) => asserts condition = (condition, message) => {
  if (!condition) throw new RunXError(message)
}
