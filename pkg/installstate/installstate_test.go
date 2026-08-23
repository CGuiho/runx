package installstate

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func isolateHome(t *testing.T) string {
	t.Helper()
	fake := t.TempDir()
	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", fake)
	} else {
		t.Setenv("HOME", fake)
	}
	return fake
}

func TestCanonicalLayout(t *testing.T) {
	fake := isolateHome(t)
	launcher, err := LauncherPath()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(fake, ".guiho", "bin", LauncherName()), launcher)

	cliDir, err := CLIDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(fake, ".guiho", "runx"), cliDir)

	versionDir, err := VersionDir("1.2.3")
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(fake, ".guiho", "runx", "versions", "1.2.3"), versionDir)

	staging, err := TempRoot()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(fake, ".guiho", ".temp"), staging)
}

func TestPointerRoundTripAndValidation(t *testing.T) {
	isolateHome(t)
	missing, err := ReadPointer()
	require.NoError(t, err)
	assert.Nil(t, missing)

	require.NoError(t, WritePointer(Pointer{Protocol: 1, Active: "1.2.3", Previous: "1.2.2"}))
	loaded, err := ReadPointer()
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, 1, loaded.Protocol)
	assert.Equal(t, "1.2.3", loaded.Active)
	assert.Equal(t, "1.2.2", loaded.Previous)

	err = WritePointer(Pointer{Protocol: 2, Active: "1.2.3"})
	require.Error(t, err)
	err = WritePointer(Pointer{Protocol: 1, Active: "../../etc"})
	require.Error(t, err)
	err = WritePointer(Pointer{Protocol: 1, Active: "1.2.3", Previous: "1.2.3"})
	require.Error(t, err)
}

func TestReadPointerAcceptsLegacyPowerShellUTF8BOM(t *testing.T) {
	home := t.TempDir()
	path := CurrentPointerPathIn(home)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
	data := append([]byte{0xef, 0xbb, 0xbf}, []byte(`{"protocol":1,"active":"1.2.3"}`)...)
	require.NoError(t, os.WriteFile(path, data, 0o644))

	pointer, err := ReadPointerIn(home)
	require.NoError(t, err)
	require.NotNil(t, pointer)
	assert.Equal(t, "1.2.3", pointer.Active)
}

func TestLedgerRejectsForeignPaths(t *testing.T) {
	isolateHome(t)
	launcherPath, err := LauncherPath()
	require.NoError(t, err)
	cliDir, err := CLIDir()
	require.NoError(t, err)

	valid := Ledger{Protocol: 1, Version: "1.0.0", Artifacts: []Artifact{
		{ID: "payload", Version: "1.0.0", Kind: "payload", Path: filepath.Join(cliDir, "versions", "1.0.0", PayloadName())},
		{ID: "launcher", Version: "1.0.0", Kind: "launcher", Path: launcherPath},
	}}
	require.NoError(t, valid.Validate())
	require.NoError(t, WriteLedger(valid))
	loaded, err := ReadLedger()
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Len(t, loaded.Artifacts, 2)

	duplicate := Ledger{Protocol: 1, Version: "1.0.0", Artifacts: []Artifact{
		{ID: "a", Version: "1.0.0", Kind: "payload", Path: filepath.Join(cliDir, "x")},
		{ID: "a", Version: "1.0.0", Kind: "payload", Path: filepath.Join(cliDir, "y")},
	}}
	require.Error(t, duplicate.Validate())

	foreign := Ledger{Protocol: 1, Version: "1.0.0", Artifacts: []Artifact{
		{ID: "evil", Version: "1.0.0", Kind: "payload", Path: filepath.Join(cliDir, "..", "other-cli", "file")},
	}}
	require.Error(t, foreign.Validate())

	traversal := Ledger{Protocol: 1, Version: "1.0.0", Artifacts: []Artifact{
		{ID: "evil", Version: "1.0.0", Kind: "payload", Path: `C:\Windows\system32\evil.exe`},
	}}
	require.Error(t, traversal.Validate())
}

func TestWriteFileAtomicNeverLeavesPartialReads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "file.json")
	require.NoError(t, WriteFileAtomic(path, []byte(`{"ok":true}`), 0o644))
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, `{"ok":true}`, string(data))
	entries, err := os.ReadDir(filepath.Dir(path))
	require.NoError(t, err)
	assert.Len(t, entries, 1, "no temporary siblings may survive")
}

func TestStagingDirIsUnderSharedTempAndCleansUp(t *testing.T) {
	isolateHome(t)
	dir, cleanup, err := StagingDir("install")
	require.NoError(t, err)
	root, err := TempRoot()
	require.NoError(t, err)
	assert.Contains(t, dir, root)
	assert.Contains(t, filepath.Base(dir), "runx-install-")
	cleanup()
	_, err = os.Stat(dir)
	assert.True(t, os.IsNotExist(err))

	_, _, err = StagingDir("../escape")
	require.Error(t, err)
}
