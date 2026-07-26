package maintenance

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMaintainAgentIntegration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-maint-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	homeDir := filepath.Join(tempDir, "home")
	projectDir := filepath.Join(tempDir, "proj")
	require.NoError(t, os.MkdirAll(projectDir, 0755))
	require.NoError(t, os.Mkdir(filepath.Join(projectDir, ".git"), 0755))

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

func TestMaintainAgentIntegrationUpdatesBothFilesAndPreservesCRLF(t *testing.T) {
	root := t.TempDir()
	home := filepath.Join(root, "home")
	repository := filepath.Join(root, "repository")
	nested := filepath.Join(repository, "packages", "worker")
	require.NoError(t, os.MkdirAll(filepath.Join(repository, ".git"), 0o755))
	require.NoError(t, os.MkdirAll(nested, 0o755))
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		require.NoError(t, os.WriteFile(filepath.Join(repository, name), []byte("# User content\r\n\r\nKeep me.\r\n"), 0o644))
	}

	first, err := MaintainAgentIntegration(nested, home, "skill\n")
	require.NoError(t, err)
	assert.Len(t, first.Skills, 2)
	assert.Len(t, first.Instructions, 2)
	for _, name := range []string{"AGENTS.md", "CLAUDE.md"} {
		content, readErr := os.ReadFile(filepath.Join(repository, name))
		require.NoError(t, readErr)
		text := string(content)
		assert.Contains(t, text, "# User content\r\n\r\nKeep me.\r\n")
		assert.Contains(t, text, strings.ReplaceAll(DefaultInstructionBlock(), "\n", "\r\n"))
		assert.NotContains(t, strings.ReplaceAll(text, "\r\n", ""), "\n")
	}
	second, err := MaintainAgentIntegration(nested, home, "skill\n")
	require.NoError(t, err)
	assert.Empty(t, second.Skills)
	assert.Empty(t, second.Instructions)
}

func TestMaintainAgentIntegrationUsesExistingClaudeOnly(t *testing.T) {
	repository := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(repository, ".git"), 0o755))
	claude := filepath.Join(repository, "CLAUDE.md")
	require.NoError(t, os.WriteFile(claude, []byte("# Claude\n"), 0o644))
	result, err := MaintainAgentIntegration(repository, filepath.Join(repository, "home"), "skill")
	require.NoError(t, err)
	assert.Equal(t, []string{claude}, result.Instructions)
	assert.False(t, PathExists(filepath.Join(repository, "AGENTS.md")))
}

func TestMaintainAgentIntegrationRejectsMalformedMarkers(t *testing.T) {
	repository := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(repository, ".git"), 0o755))
	agents := filepath.Join(repository, "AGENTS.md")
	original := "# Keep\n\n" + ManagedStart + "\nbroken\n"
	require.NoError(t, os.WriteFile(agents, []byte(original), 0o644))
	_, err := MaintainAgentIntegration(repository, filepath.Join(repository, "home"), "skill")
	require.Error(t, err)
	content, readErr := os.ReadFile(agents)
	require.NoError(t, readErr)
	assert.Equal(t, original, string(content))
}

func TestMaintainAgentIntegrationSkipsInstructionsOutsideRepository(t *testing.T) {
	directory := t.TempDir()
	result, err := MaintainAgentIntegration(directory, filepath.Join(directory, "home"), "skill")
	require.NoError(t, err)
	assert.Len(t, result.Skills, 2)
	assert.Empty(t, result.Instructions)
	assert.False(t, PathExists(filepath.Join(directory, "AGENTS.md")))
}

func TestReplaceManagedBlockStrict(t *testing.T) {
	existing := "# Existing Documentation\n\nSome user notes here.\n\n" + ManagedStart + "\nOld content\n" + ManagedEnd
	newBlock := DefaultInstructionBlock()

	replaced, err := ReplaceManagedBlockStrict(existing, newBlock)
	require.NoError(t, err)
	assert.Contains(t, replaced, "# Existing Documentation")
	assert.Contains(t, replaced, "Some user notes here.")
	assert.Contains(t, replaced, ManagedStart)
	assert.NotContains(t, replaced, "Old content")
}
