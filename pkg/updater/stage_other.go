//go:build !windows

package updater

import "fmt"

func stageWindowsReplacement(int, string, string, string, string) error {
	return fmt.Errorf("Windows replacement helper is unavailable on this platform")
}
