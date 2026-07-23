package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUninstallDryRun(t *testing.T) {
	uninstallDryRun = true
	uninstallFormat = "text"
	jsonFmt = false

	err := uninstallCmd.RunE(uninstallCmd, []string{})
	assert.NoError(t, err)

	uninstallFormat = "json"
	err = uninstallCmd.RunE(uninstallCmd, []string{})
	assert.NoError(t, err)
}
