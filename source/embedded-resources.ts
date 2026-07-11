// @ts-expect-error Bun text imports embed the bundled agent skill in native binaries.
import skill from '../skills/guiho-s-runx/SKILL.md' with { type: 'text' }

declare global {
  var __RUNX_EMBEDDED_RESOURCES__: { skill: string } | undefined
}

export const registerEmbeddedResources = (): void => {
  globalThis.__RUNX_EMBEDDED_RESOURCES__ = { skill }
}
