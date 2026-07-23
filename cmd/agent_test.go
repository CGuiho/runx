package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func executeCmd(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	RootCmd.SetOut(buf)
	RootCmd.SetErr(buf)
	RootCmd.SetArgs(args)

	err := RootCmd.Execute()
	return buf.String(), err
}

func TestAgentSkillList(t *testing.T) {
	output, err := executeCmd("agent", "skill", "list", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var skills []SkillInfo
	if err := json.Unmarshal([]byte(output), &skills); err != nil {
		t.Fatalf("failed to parse output: %v, output: %s", err, output)
	}

	if len(skills) != 1 || skills[0].ID != "guiho-s-runx" {
		t.Errorf("unexpected skill list: %+v", skills)
	}
}

func TestAgentSkillListFilter(t *testing.T) {
	output, err := executeCmd("agent", "skill", "list", "--filter", "runx", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var skills []SkillInfo
	if err := json.Unmarshal([]byte(output), &skills); err != nil {
		t.Fatalf("failed to parse output: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill matching 'runx', got %d", len(skills))
	}

	outputEmpty, err := executeCmd("agent", "skill", "list", "--filter", "nonexistent_filter_xyz", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var emptySkills []SkillInfo
	_ = json.Unmarshal([]byte(outputEmpty), &emptySkills)
	if len(emptySkills) != 0 {
		t.Errorf("expected 0 skills matching nonexistent filter, got %d", len(emptySkills))
	}
}

func TestAgentSkillShow(t *testing.T) {
	output, err := executeCmd("agent", "skill", "show", "guiho-s-runx", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var details SkillDetails
	if err := json.Unmarshal([]byte(output), &details); err != nil {
		t.Fatalf("failed to parse skill show output: %v", err)
	}

	if details.ID != "guiho-s-runx" || details.Path != "skills/guiho-s-runx/SKILL.md" {
		t.Errorf("unexpected skill details: %+v", details)
	}

	_, errInvalid := executeCmd("agent", "skill", "show", "invalid-skill-id")
	if errInvalid == nil {
		t.Errorf("expected error for invalid skill id, got nil")
	}
}

func TestAgentSkillInstallUninstallUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-test-skill-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Install local
	installOut, err := executeCmd("agent", "skill", "install", "--local", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("install error: %v", err)
	}

	var installRes map[string][]string
	if err := json.Unmarshal([]byte(installOut), &installRes); err != nil {
		t.Fatalf("failed to parse install result: %v", err)
	}

	installedPaths := installRes["installed"]
	if len(installedPaths) != 2 {
		t.Fatalf("expected 2 installed skill paths, got %d", len(installedPaths))
	}

	for _, p := range installedPaths {
		if !pathExists(p) {
			t.Errorf("installed skill file missing at %s", p)
		}
	}

	// Update local
	updateOut, err := executeCmd("agent", "skill", "update", "--local", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	var updateRes map[string][]string
	if err := json.Unmarshal([]byte(updateOut), &updateRes); err != nil {
		t.Fatalf("failed to parse update result: %v", err)
	}

	if len(updateRes["updated"]) != 2 {
		t.Errorf("expected 2 updated paths, got %d", len(updateRes["updated"]))
	}

	// Uninstall local
	uninstallOut, err := executeCmd("agent", "skill", "uninstall", "--local", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("uninstall error: %v", err)
	}

	var uninstallRes map[string][]string
	if err := json.Unmarshal([]byte(uninstallOut), &uninstallRes); err != nil {
		t.Fatalf("failed to parse uninstall result: %v", err)
	}

	for _, p := range uninstallRes["removed"] {
		if pathExists(p) {
			t.Errorf("directory still exists after uninstall at %s", p)
		}
	}
}

func TestAgentInstructionShow(t *testing.T) {
	output, err := executeCmd("agent", "instruction", "show")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, managedStart) || !strings.Contains(output, managedEnd) {
		t.Errorf("instruction show output missing managed block markers: %s", output)
	}
}

func TestAgentInstructionApplyUpdateRemove(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "runx-test-instruction-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	agentsPath := filepath.Join(tempDir, "AGENTS.md")

	// Apply instructions
	applyOut, err := executeCmd("agent", "instruction", "apply", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("apply error: %v", err)
	}

	var applyRes map[string][]string
	if err := json.Unmarshal([]byte(applyOut), &applyRes); err != nil {
		t.Fatalf("failed to parse apply result: %v", err)
	}

	if len(applyRes["updated"]) != 1 || applyRes["updated"][0] != agentsPath {
		t.Fatalf("unexpected apply updated list: %+v", applyRes)
	}

	content, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("failed to read AGENTS.md: %v", err)
	}

	if !strings.Contains(string(content), managedStart) {
		t.Errorf("AGENTS.md missing managed start marker: %s", string(content))
	}

	// Update instructions (no changes expected if identical)
	updateOut, err := executeCmd("agent", "instruction", "update", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("update error: %v", err)
	}

	var updateRes map[string][]string
	if err := json.Unmarshal([]byte(updateOut), &updateRes); err != nil {
		t.Fatalf("failed to parse update result: %v", err)
	}

	if len(updateRes["updated"]) != 0 {
		t.Errorf("expected 0 changed files on idempotent update, got %d", len(updateRes["updated"]))
	}

	// Remove instructions
	removeOut, err := executeCmd("agent", "instruction", "remove", "--cwd", tempDir, "--format", "json")
	if err != nil {
		t.Fatalf("remove error: %v", err)
	}

	var removeRes map[string][]string
	if err := json.Unmarshal([]byte(removeOut), &removeRes); err != nil {
		t.Fatalf("failed to parse remove result: %v", err)
	}

	if len(removeRes["removed"]) != 1 {
		t.Fatalf("expected 1 removed file, got %d", len(removeRes["removed"]))
	}

	afterRemove, _ := os.ReadFile(agentsPath)
	if strings.Contains(string(afterRemove), managedStart) {
		t.Errorf("managed block still present after remove: %s", string(afterRemove))
	}
}

func TestAgentPromptList(t *testing.T) {
	// List full
	output, err := executeCmd("agent", "prompt", "list", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var prompts []PromptInfo
	if err := json.Unmarshal([]byte(output), &prompts); err != nil {
		t.Fatalf("failed to parse prompt list: %v", err)
	}

	if len(prompts) != 1 || prompts[0].ID != "guiho-i-runx" {
		t.Errorf("unexpected prompt list: %+v", prompts)
	}

	// List names only
	outputNames, err := executeCmd("agent", "prompt", "list", "--names", "--format", "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var names []string
	if err := json.Unmarshal([]byte(outputNames), &names); err != nil {
		t.Fatalf("failed to parse prompt names: %v", err)
	}

	if len(names) != 1 || names[0] != "guiho-i-runx" {
		t.Errorf("unexpected prompt names: %+v", names)
	}
}

func TestAgentPromptShow(t *testing.T) {
	output, err := executeCmd("agent", "prompt", "show", "guiho-i-runx")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(output, "RunX Agent Instruction") {
		t.Errorf("prompt show output missing expected header: %s", output)
	}

	_, errInvalid := executeCmd("agent", "prompt", "show", "invalid-prompt-id")
	if errInvalid == nil {
		t.Errorf("expected error for invalid prompt id, got nil")
	}
}
