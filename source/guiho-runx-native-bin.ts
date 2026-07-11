#!/usr/bin/env bun
import { registerEmbeddedResources } from './embedded-resources.js'

registerEmbeddedResources()

const { runCliWithErrorHandling } = await import('./cli.js')
await runCliWithErrorHandling()
