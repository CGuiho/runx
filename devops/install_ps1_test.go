package devops

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPowerShellInstallerConfiguresGitBashIdempotently(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("PowerShell installer behavior is Windows-specific")
	}

	installer, err := filepath.Abs("install.ps1")
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	profile := filepath.Join(home, ".bashrc")
	const existing = "# Existing Git Bash configuration.\r\n"
	if err := os.WriteFile(profile, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	installDirectory := filepath.Join(home, ".local", "bin")

	command := strings.Join([]string{
		"$env:RUNX_INSTALLER_SOURCE_ONLY = '1'",
		". '" + quotePowerShellLiteral(installer) + "'",
		"$installDirectory = '" + quotePowerShellLiteral(installDirectory) + "'",
		"$first = Add-GitBashPath -InstallDirectory $installDirectory -HomeDirectory '" + quotePowerShellLiteral(home) + "'",
		"$second = Add-GitBashPath -InstallDirectory $installDirectory -HomeDirectory '" + quotePowerShellLiteral(home) + "'",
		"if (-not $first.Changed) { throw 'first Git Bash update did not change the profile' }",
		"if ($second.Changed) { throw 'second Git Bash update was not idempotent' }",
		`if ($first.ExportLine -cne 'export PATH="$HOME/.local/bin:$PATH"') { throw "unexpected export line: $($first.ExportLine)" }`,
		`$mixedCasePath = 'C:\Other;' + $installDirectory.ToUpperInvariant() + '\'`,
		"if (-not (Test-PathValueContains -PathValue $mixedCasePath -Directory $installDirectory)) { throw 'normalized Windows Path lookup failed' }",
	}, "; ")

	result := exec.Command(
		"powershell.exe",
		"-NoLogo",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy",
		"Bypass",
		"-Command",
		command,
	)
	output, err := result.CombinedOutput()
	if err != nil {
		t.Fatalf("PowerShell path helper failed: %v\n%s", err, output)
	}

	updated, err := os.ReadFile(profile)
	if err != nil {
		t.Fatal(err)
	}
	text := string(updated)
	if !strings.HasPrefix(text, existing) {
		t.Fatalf("existing Git Bash configuration changed:\n%s", text)
	}
	const exportLine = `export PATH="$HOME/.local/bin:$PATH"`
	if count := strings.Count(text, exportLine); count != 1 {
		t.Fatalf("export line count = %d, want 1:\n%s", count, text)
	}
	if !strings.Contains(text, "\r\n"+exportLine+"\r\n") {
		t.Fatalf("existing CRLF style was not preserved:\n%q", text)
	}
	if strings.HasPrefix(text, "\xEF\xBB\xBF") {
		t.Fatal("Git Bash profile unexpectedly gained a UTF-8 BOM")
	}
}

func TestREADMEPublishesSimplifiedPowerShellBootstrap(t *testing.T) {
	readme, err := os.ReadFile(filepath.Join("..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(readme)
	const command = "irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex"
	if !strings.Contains(text, command) {
		t.Fatalf("README does not contain the canonical PowerShell bootstrap: %s", command)
	}
	if strings.Contains(text, "& ([scriptblock]::Create((Invoke-RestMethod") {
		t.Fatal("README still contains the complex PowerShell bootstrap")
	}
}

func TestREADMEDocumentsLegacyNativeMigration(t *testing.T) {
	readme, err := os.ReadFile(filepath.Join("..", "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(readme)
	for _, expected := range []string{
		"### Migrate From RunX 0.8",
		"curl -fsSL https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.sh | bash",
		"irm https://raw.githubusercontent.com/CGuiho/runx/main/devops/install.ps1 | iex",
		"hash -r",
		"runx --version",
	} {
		if !strings.Contains(text, expected) {
			t.Errorf("README legacy migration guidance does not contain %q", expected)
		}
	}
	if strings.Contains(text, "--version '0.8.0'") {
		t.Fatal("README legacy migration guidance still recommends the pinned 0.8 installer")
	}
}

func quotePowerShellLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
