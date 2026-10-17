package welcome

import (
	"strings"
	"testing"
)

func TestRenderIsLeftAlignedWithoutCreatorLine(t *testing.T) {
	output := Render("linux", "amd64", "1.2.3", "")
	var expected strings.Builder
	expected.WriteString("\n")
	for _, line := range runxLogo {
		expected.WriteString(line + "\n")
	}
	expected.WriteString("Documented command catalog\n")
	expected.WriteString("GUIHO \u00b7 Crist\u00f3v\u00e3o GUIHO\n")
	expected.WriteString("\n")
	expected.WriteString("Platform  Linux x64\n")
	expected.WriteString("Version   v1.2.3\n")
	expected.WriteString("Author    GUIHO (www.guiho.co)\n")
	expected.WriteString("\n")
	expected.WriteString("Run runx --help to see available commands.\n")
	expected.WriteString("\n")
	if output != expected.String() {
		t.Fatalf("welcome output mismatch:\n got: %q\nwant: %q", output, expected.String())
	}
	if strings.Contains(output, "Creator") {
		t.Fatalf("welcome output contains a Creator line:\n%s", output)
	}
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, " ") {
			t.Fatalf("welcome line is not left-aligned: %q", line)
		}
	}
}

func TestRenderWithColorKeepsLeftAlignment(t *testing.T) {
	output := RenderWithColor("windows", "arm64", "v1.2.3", "", true)
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, " ") {
			t.Fatalf("coloured welcome line is indented: %q", line)
		}
	}
	for _, required := range []string{
		ansiBrightRed + ansiBold + runxLogo[0] + ansiReset + "\n",
		ansiAccent + "GUIHO \u00b7 Crist\u00f3v\u00e3o GUIHO" + ansiReset + "\n\n" + ansiSlate + "Platform" + ansiReset,
		ansiSlate + "Version " + ansiReset + "  " + ansiAccent + "v1.2.3" + ansiReset + "\n" +
			ansiSlate + "Author  " + ansiReset + "  " +
			ansiOffWhite + "GUIHO (" + ansiReset + ansiAccent + "www.guiho.co" + ansiReset + ansiOffWhite + ")" + ansiReset + "\n\nRun ",
	} {
		if !strings.Contains(output, required) {
			t.Fatalf("coloured welcome omits %q:\n%q", output, required)
		}
	}
	if strings.Contains(output, "Creator") {
		t.Fatalf("coloured welcome still shows the Creator line: %q", output)
	}
}

func TestRenderUpdateNoticeKeepsOneTrailingBlankLine(t *testing.T) {
	output := Render("linux", "amd64", "1.2.3", "\u26a0 New version available: v1.2.4\n  run runx upgrade to update")
	for _, required := range []string{
		"Run runx --help to see available commands.\n\n\u26a0 New version available: v1.2.4\n  run runx upgrade to update\n\n",
	} {
		if !strings.HasSuffix(output, required) {
			t.Fatalf("welcome output does not end with %q:\n%q", required, output)
		}
	}
	if strings.Contains(output, "\x1b[") {
		t.Fatalf("plain welcome contains ANSI: %q", output)
	}
}
