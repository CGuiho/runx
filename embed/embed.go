package embed

import (
	"embed"
)

// FS holds the embedded agent skills and documentation assets.
//
//go:embed skills/* prompts/*
var FS embed.FS
