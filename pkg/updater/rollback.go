package updater

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func GenerateUUID() string {
	var uuid [16]byte
	_, err := io.ReadFull(rand.Reader, uuid[:])
	if err != nil {
		return fmt.Sprintf("%d", os.Getpid())
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func ValidateNativeBinary(data []byte, osName string) error {
	if len(data) == 0 {
		return fmt.Errorf("downloaded executable is empty")
	}

	switch osName {
	case "windows":
		if len(data) < 2 || data[0] != 0x4D || data[1] != 0x5A { // MZ
			return fmt.Errorf("downloaded file is not a native Windows executable")
		}
	case "linux":
		if len(data) < 4 || data[0] != 0x7F || data[1] != 'E' || data[2] != 'L' || data[3] != 'F' { // ELF
			return fmt.Errorf("downloaded file is not a native Linux executable")
		}
	case "darwin":
		if !isMachO(data) {
			return fmt.Errorf("downloaded file is not a native macOS executable")
		}
	}
	return nil
}

func isMachO(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	magics := [][]byte{
		{0xFE, 0xED, 0xFA, 0xCE},
		{0xFE, 0xED, 0xFA, 0xCF},
		{0xCE, 0xFA, 0xED, 0xFE},
		{0xCF, 0xFA, 0xED, 0xFE},
		{0xCA, 0xFE, 0xBA, 0xBE},
		{0xBE, 0xFA, 0xEC, 0xFA},
	}
	for _, m := range magics {
		if bytes.Equal(data[:4], m) {
			return true
		}
	}
	return false
}

func VerifyExecutableVersion(execPath string, expectedVersion string) error {
	cmd := exec.Command(execPath, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("replacement executable failed version check: %w (%s)", err, string(out))
	}
	output := string(out)
	if expectedVersion != "" && !bytes.Contains(out, []byte(expectedVersion)) {
		return fmt.Errorf("replacement reported version %q; expected %q", output, expectedVersion)
	}
	return nil
}

func PerformReplacementAndRollback(
	execPath string,
	binaryData []byte,
	targetVersion string,
	osName string,
	fileOps FileOperations,
	verifyFunc func(path, expectedVersion string) error,
) (*UpgradeResult, string, error) {
	if fileOps == nil {
		fileOps = SystemFileOps{}
	}
	if verifyFunc == nil {
		verifyFunc = VerifyExecutableVersion
	}

	pid := os.Getpid()
	uid := GenerateUUID()
	tempPath := fmt.Sprintf("%s.new-%d-%s", execPath, pid, uid)
	backupPath := fmt.Sprintf("%s.old-%d-%s", execPath, pid, uid)

	// Write new binary
	if err := os.WriteFile(tempPath, binaryData, 0755); err != nil {
		return nil, "download_invalid", fmt.Errorf("failed to write temporary executable: %w", err)
	}
	defer func() {
		_ = fileOps.Remove(tempPath)
	}()

	if err := fileOps.MakeExecutable(tempPath); err != nil {
		return nil, "download_invalid", fmt.Errorf("failed to make executable: %w", err)
	}

	// Two-phase rename
	// Phase 1: Move current binary to backup path
	if err := fileOps.Rename(execPath, backupPath); err != nil {
		return nil, "backup_failed", fmt.Errorf("failed to backup existing binary: %w", err)
	}

	// Phase 2: Move new binary to execPath
	if err := fileOps.Rename(tempPath, execPath); err != nil {
		// Rollback phase 1
		_ = fileOps.Rename(backupPath, execPath)
		return nil, "replace_failed", fmt.Errorf("failed to replace binary: %w", err)
	}

	// Verify replacement
	if err := verifyFunc(execPath, targetVersion); err != nil {
		// Rollback replacement
		_ = fileOps.Remove(execPath)
		rbErr := fileOps.Rename(backupPath, execPath)
		if rbErr != nil {
			return nil, "rollback_failed", fmt.Errorf("verification failed (%v); rollback also failed: %w", err, rbErr)
		}
		return nil, "verification_failed", fmt.Errorf("replacement verification failed: %w", err)
	}

	// Cleanup backup
	cleanupDeferred := false
	if err := fileOps.Remove(backupPath); err != nil {
		if osName == "windows" {
			ScheduleWindowsBackupCleanup(backupPath)
			cleanupDeferred = true
		} else {
			return nil, "replace_failed", fmt.Errorf("could not delete old executable: %w", err)
		}
	}

	return &UpgradeResult{
		InstalledVersion: targetVersion,
		CleanupDeferred:  cleanupDeferred,
	}, "", nil
}
