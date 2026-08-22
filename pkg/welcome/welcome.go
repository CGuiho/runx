package welcome

import (
	"fmt"
	"runtime"
	"strings"
)

const welcomeWidth = 58

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

func borderLine(value string, width int) string {
	// "║  <value padded>║"
	return fmt.Sprintf("║  %-*s║", width, value)
}

// Render returns the deterministic RunX welcome window.
// It contains no ANSI codes and ends with a trailing newline.
// updateNotice, when non-empty, is appended after the body exactly as provided;
// it should already be in the two-line convention form.
func Render(platform, architecture, version, updateNotice string) string {
	if platform == "" {
		platform = runtime.GOOS
	}
	if architecture == "" {
		architecture = runtime.GOARCH
	}
	platform = platformLabel(platform)
	architecture = archLabel(architecture)
	contentWidth := welcomeWidth - 4
	bordered := []string{
		fmt.Sprintf("╔%s╗", strings.Repeat("═", welcomeWidth-2)),
		borderLine("RUNX", contentWidth),
		borderLine("Documented command catalog", contentWidth),
		borderLine("GUIHO \u00b7 Crist\u00f3v\u00e3o GUIHO", contentWidth),
		fmt.Sprintf("╚%s╝", strings.Repeat("═", welcomeWidth-2)),
	}
	lines := []string{}
	lines = append(lines, bordered...)
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  platform      %s %s", platform, architecture))
	lines = append(lines, fmt.Sprintf("  version       v%s", strings.TrimPrefix(version, "v")))
	lines = append(lines, "")
	lines = append(lines, "  Run `runx --help` to see available commands.")
	if strings.TrimSpace(updateNotice) != "" {
		lines = append(lines, "")
		// Preserve the notice's internal newlines exactly; do not re-wrap.
		lines = append(lines, strings.Split(strings.TrimRight(updateNotice, "\n"), "\n")...)
	}
	return strings.Join(lines, "\n") + "\n"
}
