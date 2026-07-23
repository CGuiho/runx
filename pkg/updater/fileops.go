package updater

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type SystemFileOps struct{}

func (s SystemFileOps) Rename(from, to string) error {
	if runtime.GOOS == "windows" {
		err := os.Rename(from, to)
		if err == nil {
			return nil
		}
		cmd := exec.Command("cmd.exe", "/d", "/s", "/c", "move", "/y", from, to)
		out, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			return fmt.Errorf("move failed (%v): %s", cmdErr, string(out))
		}
		return nil
	}
	return os.Rename(from, to)
}

func (s SystemFileOps) Remove(path string) error {
	if runtime.GOOS == "windows" {
		err := os.Remove(path)
		if err == nil || os.IsNotExist(err) {
			return nil
		}
		cmd := exec.Command("cmd.exe", "/d", "/s", "/c", "del", "/f", "/q", path)
		out, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil && PathExists(path) {
			return fmt.Errorf("remove failed (%v): %s", cmdErr, string(out))
		}
		return nil
	}
	return os.RemoveAll(path)
}

func (s SystemFileOps) MakeExecutable(path string) error {
	if runtime.GOOS == "windows" {
		return nil
	}
	return os.Chmod(path, 0755)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ScheduleWindowsBackupCleanup(backupPath string) {
	cmdStr := "for ($attempt = 0; $attempt -lt 300; $attempt += 1) { if (-not (Test-Path -LiteralPath $env:RUNX_BACKUP_PATH)) { exit 0 }; try { Remove-Item -LiteralPath $env:RUNX_BACKUP_PATH -Force -ErrorAction Stop; exit 0 } catch { Start-Sleep -Milliseconds 100 } }; exit 1"
	cmd := exec.Command("powershell.exe", "-NoLogo", "-NoProfile", "-NonInteractive", "-WindowStyle", "Hidden", "-Command", cmdStr)
	cmd.Env = append(os.Environ(), "RUNX_BACKUP_PATH="+backupPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	_ = cmd.Start()
}
