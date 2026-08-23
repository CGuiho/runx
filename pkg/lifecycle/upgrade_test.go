package lifecycle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/CGuiho/runx/pkg/installstate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func fixedHome(home string) func() (string, error) {
	return func() (string, error) { return home, nil }
}

func TestIsProtocolV1Installation(t *testing.T) {
	home := t.TempDir()
	assert.False(t, IsProtocolV1Installation(fixedHome(home)))
	require.NoError(t, installstate.WritePointerIn(home, installstate.Pointer{Protocol: 1, Active: "1.0.0"}))
	assert.True(t, IsProtocolV1Installation(fixedHome(home)))
}

func TestUpgradeWithoutPointerFailsClosedBeforeNetwork(t *testing.T) {
	home := t.TempDir()
	_, err := UpgradeWholeRelease(Options{HomeDir: fixedHome(home), BuildTarget: "runx-windows-amd64"})
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrLegacyInstallation), "want ErrLegacyInstallation, got %v", err)
}

func TestVerifyChecksumsDoesNotRequireManifestToHashItself(t *testing.T) {
	staging := t.TempDir()
	payload := []byte("verified payload")
	require.NoError(t, os.WriteFile(filepath.Join(staging, "payload"), payload, 0o644))
	sum := sha256.Sum256(payload)
	manifest := hex.EncodeToString(sum[:]) + "  payload\n"
	require.NoError(t, os.WriteFile(filepath.Join(staging, "checksums.txt"), []byte(manifest), 0o644))

	require.NoError(t, verifyChecksums(staging, []string{"payload", "checksums.txt"}))
	require.ErrorContains(t, verifyChecksums(staging, []string{"payload", "missing", "checksums.txt"}), "checksum entry missing for missing")
}
