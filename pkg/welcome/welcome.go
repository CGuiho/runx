package welcome

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

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

	// ansiAccent is RunX's primary brand colour, used for the GUIHO
	// attribution, the version and the website in the Author detail.
	ansiAccent = ansiBrightRed + ansiBold
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

	tagline := "Documented command catalog"
	logoColor := ansiBrightRed
	taglineColor := ansiSlate

	var b strings.Builder
	// One blank line before, so the welcome never touches the prompt line
	b.WriteString("\n")

	for _, line := range runxLogo {
		b.WriteString(colorize(withColor, line, logoColor+ansiBold))
		b.WriteString("\n")
	}
	b.WriteString(colorize(withColor, tagline, taglineColor+ansiDim))
	b.WriteString("\n")
	b.WriteString(colorize(withColor, "GUIHO \u00b7 Crist\u00f3v\u00e3o GUIHO", ansiAccent))
	b.WriteString("\n\n")

	// Details — platform, version, author
	writeDetail(&b, withColor, "Platform", colorize(withColor, platform+" "+architecture, ansiOffWhite))
	writeDetail(&b, withColor, "Version", colorize(withColor, "v"+version, ansiAccent))
	writeDetail(&b, withColor, "Author",
		colorize(withColor, "GUIHO (", ansiOffWhite)+
			colorize(withColor, "www.guiho.co", ansiAccent)+
			colorize(withColor, ")", ansiOffWhite))
	b.WriteString("\n")

	// Help hint
	b.WriteString("Run ")
	b.WriteString(colorize(withColor, "runx --help", ansiOffWhite+ansiBold))
	b.WriteString(" to see available commands.")
	b.WriteString("\n")

	// Update notice (below the help hint)
	if strings.TrimSpace(updateNotice) != "" {
		// updateNotice is expected as "⚠ New version available: vX.Y.Z\n  run runx upgrade to update"
		parts := strings.Split(strings.TrimRight(updateNotice, "\n"), "\n")
		b.WriteString("\n")
		// First line triangle in the brand accent
		if len(parts) > 0 {
			b.WriteString(colorize(withColor, parts[0], ansiAccent))
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
		b.WriteString("\n")
	}
	// One blank line after, so the welcome never touches the next prompt
	b.WriteString("\n")
	return b.String()
}

// detailLabelWidth fits the longest detail label ("Platform") so values align.
const detailLabelWidth = 8

// writeDetail writes one left-aligned detail row; value is already
// colorized by the caller when colour is enabled.
func writeDetail(b *strings.Builder, withColor bool, label, value string) {
	padded := colorize(withColor, fmt.Sprintf("%-*s", detailLabelWidth, label), ansiSlate)
	fmt.Fprintf(b, "%s  %s\n", padded, value)
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
