/**
 * @copyright Copyright © 2026 GUIHO Technologies as represented by Cristóvão GUIHO. All Rights Reserved.
 */

// @ts-expect-error Bun text imports embed the bundled skill in native binaries.
import skill from '../skills/guiho-s-runx/SKILL.md' with { type: 'text' }
// @ts-expect-error Bun text imports embed the bundled instruction prompt in native binaries.
import prompt from '../prompts/guiho-i-runx.md' with { type: 'text' }

export {
  registerEmbeddedResources,
}

declare global {
  var __RUNX_EMBEDDED_RESOURCES__: { skill: string, prompt: string } | undefined
}

function registerEmbeddedResources(): void {
  globalThis.__RUNX_EMBEDDED_RESOURCES__ = { skill, prompt }
}
