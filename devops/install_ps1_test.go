package devops

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// readInstaller loads a devops lifecycle script from this package directory.
func readInstaller(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestInstallersUseCanonicalLayout(t *testing.T) {
	sh := readInstaller(t, "install.sh")
	ps1 := readInstaller(t, "install.ps1")
	layoutElements := []string{".guiho", ".temp", "runx-install-", "versions", "current.json"}
	for _, expected := range layoutElements {
		if !strings.Contains(sh, expected) {
			t.Errorf("install.sh is missing canonical layout element %q", expected)
		}
		if !strings.Contains(ps1, expected) {
			t.Errorf("install.ps1 is missing canonical layout element %q", expected)
		}
	}
	for name, script := range map[string]string{"install.sh": sh, "install.ps1": ps1} {
		if !strings.Contains(script, "__self-test") {
			t.Errorf("%s does not run the hidden installation self-test", name)
		}
	}
}

func TestInstallersAcceptOnlyFullFlagNames(t *testing.T) {
	sh := readInstaller(t, "install.sh")
	ps1 := readInstaller(t, "install.ps1")
	shForbidden := []string{"-v)", "-h)", "--install-dir", "RUNX_INSTALL_DIR"}
	for _, forbidden := range shForbidden {
		if strings.Contains(sh, forbidden) {
			t.Errorf("install.sh contains forbidden alias or override %q", forbidden)
		}
	}
	ps1Forbidden := []string{"RUNX_INSTALL_DIR", "$InstallDir"}
	for _, forbidden := range ps1Forbidden {
		if strings.Contains(ps1, forbidden) {
			t.Errorf("install.ps1 contains forbidden alias or override %q", forbidden)
		}
	}
	if !strings.Contains(sh, "--channel") || !strings.Contains(sh, "--version") {
		t.Error("install.sh must support --version and --channel")
	}
	if !strings.Contains(ps1, "[string]$Channel") || !strings.Contains(ps1, "[string]$Version") {
		t.Error("install.ps1 must support -Version and -Channel")
	}
}

func TestUninstallersShareTheContract(t *testing.T) {
	cases := []struct{ name, text string }{
		{"uninstall.sh", readInstaller(t, "uninstall.sh")},
		{"uninstall.ps1", readInstaller(t, "uninstall.ps1")},
	}
	for _, testCase := range cases {
		lower := strings.ToLower(testCase.text)
		options := map[string][2]string{
			"preserve-config": {"--preserve-config", "-preserveconfig"},
			"preserve-data":   {"--preserve-data", "-preservedata"},
			"dry-run":         {"--dry-run", "-dryrun"},
		}
		for option, variants := range options {
			if !strings.Contains(lower, variants[0]) && !strings.Contains(lower, variants[1]) {
				t.Errorf("%s is missing the shared contract option %q", testCase.name, option)
			}
		}
		if !strings.Contains(lower, "yes") {
			t.Errorf("%s is missing the --yes confirmation option", testCase.name)
		}
		if !strings.Contains(testCase.text, "PRESERVE") || !strings.Contains(testCase.text, "REMOVE") {
			t.Errorf("%s must display a REMOVE/PRESERVE plan", testCase.name)
		}
	}
}

func TestPowerShellScriptsParse(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("PowerShell parsing is Windows-specific")
	}
	for _, name := range []string{"install.ps1", "uninstall.ps1"} {
		scriptPath, err := filepath.Abs(name)
		if err != nil {
			t.Fatal(err)
		}
		command := exec.Command(
			"powershell.exe",
			"-NoLogo",
			"-NoProfile",
			"-NonInteractive",
			"-ExecutionPolicy",
			"Bypass",
			"-Command",
			"try { [void][scriptblock]::Create((Get-Content -Raw -LiteralPath '"+scriptPath+"')); exit 0 } catch { Write-Error $_; exit 1 }",
		)
		if output, err := command.CombinedOutput(); err != nil {
			t.Errorf("%s does not parse: %v\n%s", name, err, output)
		}
	}
}

func TestBashScriptsParse(t *testing.T) {
	if runtime.GOOS == "windows" {
		bashPath, err := exec.LookPath("bash.exe")
		if err != nil {
			if _, err2 := exec.LookPath("bash"); err2 != nil {
				t.Skip("bash is unavailable for syntax validation")
			}
			bashPath = "bash"
		}
		for _, name := range []string{"install.sh", "uninstall.sh"} {
			command := exec.Command(bashPath, "-n", name)
			if output, err := command.CombinedOutput(); err != nil {
				t.Errorf("%s does not parse: %v\n%s", name, err, output)
			}
		}
		return
	}
	for _, name := range []string{"install.sh", "uninstall.sh"} {
		command := exec.Command("bash", "-n", name)
		if output, err := command.CombinedOutput(); err != nil {
			t.Errorf("%s does not parse: %v\n%s", name, err, output)
		}
	}
}

func quotePowerShellLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
