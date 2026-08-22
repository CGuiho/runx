package welcome

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const welcomeWidth = 58

const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b[1m"
	ansiDim   = "\x1b[2m"

	// Palette from user: 2B2D42, 8D99AE, EDF2F4, EF233C, D90429
	ansiDark      = "\x1b[38;2;43;45;66m"
	ansiSlate     = "\x1b[38;2;141;153;174m"
	ansiOffWhite  = "\x1b[38;2;237;242;244m"
	ansiBrightRed = "\x1b[38;2;239;35;60m"
	ansiDarkRed   = "\x1b[38;2;217;4;41m"
)

var platformLabels = map[string]string{
	"darwin":  "macOS",
	"linux":   "Linux",
	"windows": "Windows",
	"win32":   "Windows",
}

var archLabels = map[string]string{
	"amd64": "x64",
	"386":   "x86",
	"arm64": "arm64",
	"arm":   "arm",
}

func platformLabel(goos string) string {
	if label, ok := platformLabels[goos]; ok {
		return label
	}
	return goos
}

func archLabel(goarch string) string {
	if label, ok := archLabels[goarch]; ok {
		return label
	}
	return goarch
}

// RUNX logo in ANSI Shadow style (6 rows, block letters)
var runxLogo = []string{
	"██████╗ ██╗   ██╗███╗  ██╗██╗  ██╗",
	"██╔══██╗██║   ██║████╗ ██║╚██╗██╔╝",
	"██████╔╝██║   ██║██╔██╗██║ ╚███╔╝ ",
	"██╔══██╗██║   ██║██║╚████║ ██╔██╗ ",
	"██║  ██║╚██████╔╝██║ ╚███║██╔╝╚██╗",
	"╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚══╝╚═╝  ╚═╝",
}

func colorize(enabled bool, s, code string) string {
	if !enabled || code == "" {
		return s
	}
	return code + s + ansiReset
}

func useColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	// Default to color when stdout is a terminal; tests override via -no-color
	return true
}

// Render returns the deterministic RunX welcome window without ANSI.
// It is used by tests and non-terminal output.
func Render(platform, architecture, version, updateNotice string) string {
	return renderWithColor(platform, architecture, version, updateNotice, false)
}

// RenderWithColor returns the welcome window with optional ANSI palette.
func RenderWithColor(platform, architecture, version, updateNotice string, withColor bool) string {
	return renderWithColor(platform, architecture, version, updateNotice, withColor)
}

func renderWithColor(platform, architecture, version, updateNotice string, withColor bool) string {
	if platform == "" {
		platform = runtime.GOOS
	}
	if architecture == "" {
		architecture = runtime.GOARCH
	}
	platform = platformLabel(platform)
	architecture = archLabel(architecture)
	version = strings.TrimPrefix(version, "v")

	// Box width: widest logo line + 4 padding
	innerWidth := 0
	for _, line := range runxLogo {
		if l := len([]rune(line)); l > innerWidth {
			innerWidth = l
		}
	}
	innerWidth += 4
	if innerWidth < 56 {
		innerWidth = 56
	}
	// Ensure tagline fits
	tagline := "Documented command catalog"
	if len([]rune(tagline))+4 > innerWidth {
		innerWidth = len([]rune(tagline)) + 4
	}

	var b strings.Builder
	logoColor := ansiBrightRed
	taglineColor := ansiSlate

	// Logo centered — no border per user request
	for _, line := range runxLogo {
		lineLen := len([]rune(line))
		pad := innerWidth - lineLen
		left := pad / 2
		right := pad - left
		inner := strings.Repeat(" ", left) + line + strings.Repeat(" ", right)
		b.WriteString(colorize(withColor, strings.TrimRight(inner, " "), logoColor+ansiBold))
		b.WriteString("\n")
	}
	// Tagline centered, dim slate — no border
	tagPad := innerWidth - len([]rune(tagline))
	tagLeft := tagPad / 2
	tagRight := tagPad - tagLeft
	tagInner := strings.Repeat(" ", tagLeft) + tagline + strings.Repeat(" ", tagRight)
	b.WriteString(colorize(withColor, strings.TrimRight(tagInner, " "), taglineColor+ansiDim))
	b.WriteString("\n")
	// GUIHO line centered — no border
	guihoLine := "GUIHO \u00b7 Crist\u00f3v\u00e3o GUIHO"
	guihoPad := innerWidth - len([]rune(guihoLine))
	guihoLeft := guihoPad / 2
	guihoRight := guihoPad - guihoLeft
	guihoInner := strings.Repeat(" ", guihoLeft) + guihoLine + strings.Repeat(" ", guihoRight)
	b.WriteString(colorize(withColor, strings.TrimRight(guihoInner, " "), taglineColor))
	b.WriteString("\n")

	// ── Outside box ─────────────────────────
	b.WriteString("\n\n")
	b.WriteString(colorize(withColor, "GUIHO", ansiOffWhite+ansiBold))
	b.WriteString(colorize(withColor, "  —  Documented command catalog", ansiSlate))
	b.WriteString("\n\n")

	// Creator
	b.WriteString(colorize(withColor, fmt.Sprintf("%-9s", "Creator"), ansiSlate))
	b.WriteString("  ")
	b.WriteString(colorize(withColor, "Cristóvão GUIHO", ansiOffWhite+ansiBold))
	b.WriteString("\n")
	// Platform
	if withColor {
		b.WriteString("  ")
		b.WriteString(colorize(true, "platform", ansiSlate))
		b.WriteString("      ")
		b.WriteString(colorize(true, platform+" "+architecture, ansiOffWhite))
		b.WriteString("\n")
	} else {
		b.WriteString(fmt.Sprintf("  platform      %s %s\n", platform, architecture))
	}
	// Version
	if withColor {
		b.WriteString("  ")
		b.WriteString(colorize(true, "version", ansiSlate))
		b.WriteString("       ")
		b.WriteString(colorize(true, "v"+version, ansiBrightRed+ansiBold))
		b.WriteString("\n\n")
	} else {
		b.WriteString(fmt.Sprintf("  version       v%s\n\n", version))
	}

	// Help hint
	b.WriteString("Run ")
	b.WriteString(colorize(withColor, "runx --help", ansiOffWhite+ansiBold))
	b.WriteString(" to see available commands.")

	// Update notice (inside hello, below help)
	if strings.TrimSpace(updateNotice) != "" {
		// updateNotice is expected as "⚠ New version available: vX.Y.Z\n  run runx upgrade to update"
		parts := strings.Split(strings.TrimRight(updateNotice, "\n"), "\n")
		b.WriteString("\n\n")
		// First line triangle in bright red
		if len(parts) > 0 {
			b.WriteString(colorize(withColor, parts[0], ansiBrightRed+ansiBold))
		}
		if len(parts) > 1 {
			b.WriteString("\n")
			// Second line: "  run runx upgrade to update" -> color "run" slate, command off-white
			second := parts[1]
			// Try to split to color command
			if idx := strings.Index(second, "runx"); idx >= 0 {
				prefix := second[:idx]
				rest := second[idx:]
				b.WriteString(colorize(withColor, prefix, ansiSlate))
				b.WriteString(colorize(withColor, rest, ansiOffWhite+ansiBold))
			} else {
				b.WriteString(colorize(withColor, second, ansiSlate))
			}
		}
		if len(parts) > 2 {
			for _, p := range parts[2:] {
				b.WriteString("\n")
				b.WriteString(colorize(withColor, p, ansiSlate))
			}
		}
	}
	b.WriteString("\n")
	return b.String()
}

// ShouldUseColor reports whether colored output should be used for the current terminal.
func ShouldUseColor(isTerminal bool) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return isTerminal
}
