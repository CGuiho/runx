package devops

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/CGuiho/runx/pkg/lifecycle"
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
		powerShellCache := t.TempDir()
		command.Env = append(os.Environ(), "LOCALAPPDATA="+powerShellCache, "APPDATA="+powerShellCache)
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

func TestPowerShellInstallerAndWholeReleaseUpgrade(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("the native PowerShell installer integration test requires Windows")
	}
	powerShell, err := exec.LookPath("powershell.exe")
	if err != nil {
		t.Skip("Windows PowerShell is unavailable")
	}

	const version = "9.8.7"
	const upgradeVersion = "9.8.8"
	repositoryRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	fixture := t.TempDir()
	buildWindowsFixture(t, repositoryRoot, fixture, version)

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		name := filepath.Base(request.URL.Path)
		data, readErr := os.ReadFile(filepath.Join(fixture, name))
		if readErr != nil {
			http.NotFound(response, request)
			return
		}
		response.Header().Set("Content-Length", fmt.Sprint(len(data)))
		if request.Method != http.MethodHead {
			_, _ = response.Write(data)
		}
	}))
	defer server.Close()

	home := t.TempDir()
	project := t.TempDir()
	script := filepath.Join(repositoryRoot, "devops", "install.ps1")
	runInstaller := func(label string, injectFailure bool) {
		t.Helper()
		command := exec.Command(
			powerShell,
			"-NoLogo", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
			"-File", script, "-Version", version,
		)
		command.Dir = project
		command.Env = append(os.Environ(),
			"HOME="+home,
			"USERPROFILE="+home,
			"LOCALAPPDATA="+filepath.Join(home, "AppData", "Local"),
			"APPDATA="+filepath.Join(home, "AppData", "Roaming"),
			"RUNX_INSTALL_TEST_MODE=1",
			"RUNX_INSTALL_TEST_DOWNLOAD_ROOT="+server.URL+"/releases/download",
			"RUNX_DISABLE_UPDATE_WORKER=1",
			"RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1",
		)
		if injectFailure {
			command.Env = append(command.Env, "RUNX_INSTALL_TEST_FAIL_AFTER_ACTIVATION=1")
		}
		output, runErr := command.CombinedOutput()
		if injectFailure {
			if runErr == nil {
				t.Fatalf("%s unexpectedly succeeded:\n%s", label, output)
			}
			if !strings.Contains(string(output), "was rolled back") {
				t.Fatalf("%s did not report rollback:\n%s", label, output)
			}
			return
		}
		if runErr != nil {
			t.Fatalf("%s failed: %v\n%s", label, runErr, output)
		}
		if !strings.Contains(string(output), "[OK] Installed and verified RunX "+version) {
			t.Fatalf("%s did not report verified installation:\n%s", label, output)
		}
	}

	assertInstallation := func(label, wantVersion string) {
		t.Helper()
		pointerBytes, readErr := os.ReadFile(installstate.CurrentPointerPathIn(home))
		if readErr != nil {
			t.Fatalf("%s read current.json: %v", label, readErr)
		}
		if len(pointerBytes) >= 3 && pointerBytes[0] == 0xef && pointerBytes[1] == 0xbb && pointerBytes[2] == 0xbf {
			t.Fatalf("%s wrote current.json with a UTF-8 BOM: %x", label, pointerBytes[:3])
		}
		pointer, pointerErr := installstate.ReadPointerIn(home)
		if pointerErr != nil {
			t.Fatalf("%s strict pointer decode failed: %v", label, pointerErr)
		}
		if pointer == nil || pointer.Active != wantVersion {
			t.Fatalf("%s active pointer = %#v, want %s", label, pointer, wantVersion)
		}

		launcher := installstate.LauncherPathIn(home)
		command := exec.Command(launcher, "--version")
		command.Env = append(os.Environ(),
			"HOME="+home,
			"USERPROFILE="+home,
			"LOCALAPPDATA="+filepath.Join(home, "AppData", "Local"),
			"APPDATA="+filepath.Join(home, "AppData", "Roaming"),
			"RUNX_DISABLE_UPDATE_WORKER=1",
			"RUNX_DISABLE_AGENT_MAINTENANCE_WORKER=1",
		)
		output, runErr := command.CombinedOutput()
		if runErr != nil {
			t.Fatalf("%s launcher failed: %v\n%s", label, runErr, output)
		}
		if got := strings.TrimSpace(string(output)); got != wantVersion {
			t.Fatalf("%s launcher version = %q, want %q", label, got, wantVersion)
		}
	}

	runInstaller("clean install", false)
	assertInstallation("clean install", version)
	runInstaller("injected same-version activation failure", true)
	assertInstallation("rollback after injected same-version activation failure", version)
	runInstaller("same-version reinstall", false)
	assertInstallation("same-version reinstall", version)

	upgradeFixture := t.TempDir()
	buildWindowsFixture(t, repositoryRoot, upgradeFixture, upgradeVersion)
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if request.URL.Host == "api.github.com" {
			catalog := []map[string]any{{
				"tag_name": "runx/v" + upgradeVersion,
				"draft":    false, "prerelease": false,
				"assets": []map[string]string{
					{"name": "runx-payload-windows-amd64.exe", "browser_download_url": "https://fixture/runx-payload-windows-amd64.exe"},
					{"name": "checksums.txt", "browser_download_url": "https://fixture/checksums.txt"},
				},
			}}
			body, marshalErr := json.Marshal(catalog)
			if marshalErr != nil {
				return nil, marshalErr
			}
			return httpResponse(request, http.StatusOK, body), nil
		}
		name := filepath.Base(request.URL.Path)
		body, readErr := os.ReadFile(filepath.Join(upgradeFixture, name))
		if readErr != nil {
			return httpResponse(request, http.StatusNotFound, nil), nil
		}
		return httpResponse(request, http.StatusOK, body), nil
	})
	client := &http.Client{Transport: transport}
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	requireChdir(t, project)
	defer requireChdir(t, originalDirectory)

	result, err := lifecycle.UpgradeWholeRelease(lifecycle.Options{
		CurrentVersion: version,
		BuildTarget:    "runx-payload-windows-amd64",
		HTTPClient:     client,
		HomeDir:        func() (string, error) { return home, nil },
	})
	if err != nil {
		t.Fatalf("whole-release upgrade failed: %v", err)
	}
	if result.Outcome != "upgraded" || !result.Verified {
		t.Fatalf("whole-release result = %#v, want upgraded and verified", result)
	}
	assertInstallation("whole-release upgrade with standard checksum manifest", upgradeVersion)

	sameVersionResult, err := lifecycle.UpgradeWholeRelease(lifecycle.Options{
		CurrentVersion:   upgradeVersion,
		RequestedVersion: upgradeVersion,
		BuildTarget:      "runx-payload-windows-amd64",
		HTTPClient:       client,
		HomeDir:          func() (string, error) { return home, nil },
	})
	if err != nil {
		t.Fatalf("same-version whole-release check failed: %v", err)
	}
	if sameVersionResult.Outcome != "up-to-date" || !sameVersionResult.Verified {
		t.Fatalf("same-version result = %#v, want verified up-to-date", sameVersionResult)
	}
	assertInstallation("same-version whole-release idempotence", upgradeVersion)
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func httpResponse(request *http.Request, status int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Request:    request,
	}
}

func requireChdir(t *testing.T, directory string) {
	t.Helper()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("change directory to %s: %v", directory, err)
	}
}

func buildWindowsFixture(t *testing.T, repositoryRoot, fixture, version string) {
	t.Helper()
	build := func(output, packagePath string, ldflags string) {
		t.Helper()
		arguments := []string{"build", "-trimpath"}
		if ldflags != "" {
			arguments = append(arguments, "-ldflags", ldflags)
		}
		arguments = append(arguments, "-o", output, packagePath)
		command := exec.Command("go", arguments...)
		command.Dir = repositoryRoot
		command.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=windows", "GOARCH=amd64", "GOAMD64=v1")
		if outputBytes, err := command.CombinedOutput(); err != nil {
			t.Fatalf("build fixture %s: %v\n%s", packagePath, err, outputBytes)
		}
	}

	payload := filepath.Join(fixture, "runx-payload-windows-amd64.exe")
	launcher := filepath.Join(fixture, "runx-launcher-windows-amd64.exe")
	build(payload, ".", fmt.Sprintf("-X main.version=%s -X main.commit=installer-test -X main.buildDate=2026-01-01T00:00:00Z -X main.buildTarget=runx-windows-amd64", version))
	build(launcher, "./cmd/runx-launcher", "")

	archivePath := filepath.Join(fixture, "guiho-s-runx.zip")
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		t.Fatal(err)
	}
	archive := zip.NewWriter(archiveFile)
	skill, err := archive.Create("guiho-s-runx/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := skill.Write([]byte("---\nname: guiho-s-runx\n---\n# Fixture\n")); err != nil {
		t.Fatal(err)
	}
	if err := archive.Close(); err != nil {
		t.Fatal(err)
	}
	if err := archiveFile.Close(); err != nil {
		t.Fatal(err)
	}

	textFiles := map[string]string{
		"guiho-i-runx.md":           "# Fixture instruction\n",
		"guiho-p-runx.md":           "---\nname: guiho-p-runx\n---\n# Fixture prompt\n",
		"guiho-p-runx-uninstall.md": "---\nname: guiho-p-runx-uninstall\n---\n# Fixture uninstall prompt\n",
		"runx.schema.json":          "{}\n",
		"runx.global.schema.json":   "{}\n",
	}
	for name, content := range textFiles {
		if err := os.WriteFile(filepath.Join(fixture, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	assetNames := []string{
		filepath.Base(payload), filepath.Base(launcher), "guiho-s-runx.zip",
		"guiho-i-runx.md", "guiho-p-runx.md", "guiho-p-runx-uninstall.md",
		"runx.schema.json", "runx.global.schema.json",
	}
	type manifestEntry struct {
		File   string `json:"file"`
		SHA256 string `json:"sha256"`
	}
	manifest := struct {
		Schema    int             `json:"schema"`
		Protocol  int             `json:"protocol"`
		Name      string          `json:"name"`
		Version   string          `json:"version"`
		Artifacts []manifestEntry `json:"artifacts"`
	}{Schema: 1, Protocol: 1, Name: "runx", Version: version}
	for _, name := range assetNames {
		manifest.Artifacts = append(manifest.Artifacts, manifestEntry{File: name, SHA256: fixtureSHA256(t, filepath.Join(fixture, name))})
	}
	manifestBytes, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(fixture, "artifacts.json"), append(manifestBytes, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	assetNames = append(assetNames, "artifacts.json")

	var checksums strings.Builder
	for _, name := range assetNames {
		fmt.Fprintf(&checksums, "%s  %s\n", fixtureSHA256(t, filepath.Join(fixture, name)), name)
	}
	if err := os.WriteFile(filepath.Join(fixture, "checksums.txt"), []byte(checksums.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

func fixtureSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
