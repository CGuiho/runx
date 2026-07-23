package maintenance

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	SkillID              = "guiho-s-runx"
	PromptID             = "guiho-i-runx"
	ManagedStart         = "<!-- BEGIN RUNX — DO NOT EDIT THIS SECTION -->"
	ManagedEnd           = "<!-- END RUNX -->"
	MojibakeManagedStart = "<!-- BEGIN RUNX \u00e2\u20ac\u201d DO NOT EDIT THIS SECTION -->"
	LegacyManagedStart   = "<!-- BEGIN RUNX AGENT INSTRUCTIONS -->"
	LegacyManagedEnd     = "<!-- END RUNX AGENT INSTRUCTIONS -->"
)

type MaintenanceResult struct {
	Skills       []string `json:"skills"`
	Instructions []string `json:"instructions"`
}

func DefaultInstructionBlock() string {
	return ManagedStart + `
## RunX Command Catalog

Load the ` + "`guiho-s-runx`" + ` skill whenever discovering commands, creating or
updating catalog entries, validating ` + "`runx.yaml`" + `, inspecting command details,
or executing RunX commands.
Start with ` + "`runx check --format json`" + ` and ` + "`runx list --format json`" + `, select
stable UIDs, use ` + "`runx describe <uid>`" + `, and run
` + "`runx run --dry-run <uid>`" + ` before unfamiliar or side-effecting work.
RunX options precede the selector; post-selector tokens belong to the child.
` + ManagedEnd + "\n"
}

func SkillDirectories(homeDir string) []string {
	if homeDir == "" {
		homeDir = HomeDirectory()
	}
	return []string{
		filepath.Join(homeDir, ".agents", "skills", SkillID),
		filepath.Join(homeDir, ".claude", "skills", SkillID),
	}
}

func NearestAgentsPath(cwd string) string {
	abs, err := filepath.Abs(cwd)
	if err != nil {
		abs = cwd
	}
	curr := abs
	for {
		cand := filepath.Join(curr, "AGENTS.md")
		if PathExists(cand) {
			return cand
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return filepath.Join(abs, "AGENTS.md")
}

func MaintainAgentIntegration(cwd string, homeDir string, embeddedSkill string) (*MaintenanceResult, error) {
	result := &MaintenanceResult{
		Skills:       []string{},
		Instructions: []string{},
	}

	if embeddedSkill != "" {
		dirs := SkillDirectories(homeDir)
		for _, dir := range dirs {
			skillPath := filepath.Join(dir, "SKILL.md")
			existing, _ := ReadTextIfExists(skillPath)
			if existing != embeddedSkill {
				if err := WriteTextFileAtomic(skillPath, embeddedSkill); err == nil {
					result.Skills = append(result.Skills, skillPath)
				}
			}
		}
	}

	targetAgentsPath := NearestAgentsPath(cwd)
	existingText, _ := ReadTextIfExists(targetAgentsPath)
	nextText := ReplaceManagedBlock(existingText, DefaultInstructionBlock())

	if nextText != existingText {
		if err := WriteTextFileAtomic(targetAgentsPath, nextText); err == nil {
			result.Instructions = append(result.Instructions, targetAgentsPath)
		}
	}

	return result, nil
}

func ReplaceManagedBlock(existing string, newBlock string) string {
	stripped := strings.TrimRight(RemoveKnownManagedBlocks(existing, false), " \t\r\n")
	if stripped == "" {
		return newBlock
	}
	return stripped + "\n\n" + newBlock
}

func RemoveKnownManagedBlocks(existing string, includeLeadingWhitespace bool) string {
	output := existing
	pairs := [][2]string{
		{ManagedStart, ManagedEnd},
		{MojibakeManagedStart, ManagedEnd},
		{LegacyManagedStart, LegacyManagedEnd},
	}

	for _, pair := range pairs {
		start := regexp.QuoteMeta(pair[0])
		end := regexp.QuoteMeta(pair[1])
		prefix := ""
		if includeLeadingWhitespace {
			prefix = `\s*`
		}
		pattern := regexp.MustCompile(prefix + start + `(?s:.*?)` + end + `\s*`)
		replacement := ""
		if includeLeadingWhitespace {
			replacement = "\n"
		}
		output = pattern.ReplaceAllString(output, replacement)
	}

	return output
}

func SpawnAgentMaintenanceWorker(execPath string, cwd string) error {
	if os.Getenv("RUNX_DISABLE_AGENT_MAINTENANCE_WORKER") == "1" {
		return nil
	}

	if execPath == "" {
		var err error
		execPath, err = os.Executable()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(execPath, "--maintain-agent-integration-worker", cwd)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	return cmd.Start()
}
