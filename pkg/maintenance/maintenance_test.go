package maintenance

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNearestAgentsPath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-agents-path-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	subDir := filepath.Join(tempDir, "a", "b", "c")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	agentsPath := filepath.Join(tempDir, "a", "AGENTS.md")
	require.NoError(t, os.WriteFile(agentsPath, []byte("# Test Agents"), 0644))

	found := NearestAgentsPath(subDir)
	assert.Equal(t, agentsPath, found)
}

func TestMaintainAgentIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-maint-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	homeDir := filepath.Join(tempDir, "home")
	projectDir := filepath.Join(tempDir, "proj")
	require.NoError(t, os.MkdirAll(projectDir, 0755))

	mockSkillContent := "---\nname: guiho-s-runx\n---\n# Skill Content"

	res, err := MaintainAgentIntegration(projectDir, homeDir, mockSkillContent)
	require.NoError(t, err)
	require.NotNil(t, res)

	// Check skills installed
	assert.Len(t, res.Skills, 2) // .agents and .claude
	for _, skillFile := range res.Skills {
		content, err := ReadTextIfExists(skillFile)
		require.NoError(t, err)
		assert.Equal(t, mockSkillContent, content)
	}

	// Check AGENTS.md created/updated
	assert.Len(t, res.Instructions, 1)
	agentsFile := res.Instructions[0]
	agentsContent, err := ReadTextIfExists(agentsFile)
	require.NoError(t, err)
	assert.Contains(t, agentsContent, ManagedStart)
	assert.Contains(t, agentsContent, ManagedEnd)

	// Idempotency check - second run should produce no changes
	res2, err := MaintainAgentIntegration(projectDir, homeDir, mockSkillContent)
	require.NoError(t, err)
	assert.Empty(t, res2.Skills)
	assert.Empty(t, res2.Instructions)
}

func TestReplaceManagedBlock(t *testing.T) {
	existing := "# Existing Documentation\n\nSome user notes here.\n\n" + ManagedStart + "\nOld content\n" + ManagedEnd
	newBlock := DefaultInstructionBlock()

	replaced := ReplaceManagedBlock(existing, newBlock)
	assert.Contains(t, replaced, "# Existing Documentation")
	assert.Contains(t, replaced, "Some user notes here.")
	assert.Contains(t, replaced, ManagedStart)
	assert.NotContains(t, replaced, "Old content")
}
