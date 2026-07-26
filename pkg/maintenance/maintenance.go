package maintenance

import (
	"fmt"
	"os"
	"path/filepath"
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

func MaintainAgentIntegration(cwd string, homeDir string, embeddedSkill string) (*MaintenanceResult, error) {
	result := &MaintenanceResult{
		Skills:       []string{},
		Instructions: []string{},
	}

	if embeddedSkill != "" {
		dirs := SkillDirectories(homeDir)
		for _, dir := range dirs {
			skillPath := filepath.Join(dir, "SKILL.md")
			existing, err := ReadTextIfExists(skillPath)
			if err != nil {
				return result, err
			}
			if existing != embeddedSkill {
				if err := WriteTextFileAtomic(skillPath, embeddedSkill); err != nil {
					return result, err
				}
				result.Skills = append(result.Skills, skillPath)
			}
		}
	}

	repositoryRoot, found := RepositoryRoot(cwd)
	if !found {
		return result, nil
	}
	for _, instructionPath := range InstructionTargets(repositoryRoot) {
		existingText, err := ReadTextIfExists(instructionPath)
		if err != nil {
			return result, err
		}
		nextText, err := ReplaceManagedBlockStrict(existingText, DefaultInstructionBlock())
		if err != nil {
			return result, fmt.Errorf("update managed RunX block in %s: %w", instructionPath, err)
		}
		if nextText != existingText {
			if err := WriteTextFileAtomic(instructionPath, nextText); err != nil {
				return result, err
			}
			result.Instructions = append(result.Instructions, instructionPath)
		}
	}

	return result, nil
}

func RepositoryRoot(cwd string) (string, bool) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return "", false
		}
	}
	current, err := filepath.Abs(cwd)
	if err != nil {
		return "", false
	}
	for {
		if PathExists(filepath.Join(current, ".git")) {
			return current, true
		}
		parent := filepath.Dir(current)
		if parent == current {
			return "", false
		}
		current = parent
	}
}

func InstructionTargets(repositoryRoot string) []string {
	agents := filepath.Join(repositoryRoot, "AGENTS.md")
	claude := filepath.Join(repositoryRoot, "CLAUDE.md")
	agentsExists, claudeExists := PathExists(agents), PathExists(claude)
	if agentsExists && claudeExists {
		return []string{agents, claude}
	}
	if claudeExists {
		return []string{claude}
	}
	return []string{agents}
}

type markerPair struct{ start, end string }

var managedMarkerPairs = []markerPair{{ManagedStart, ManagedEnd}, {MojibakeManagedStart, ManagedEnd}, {LegacyManagedStart, LegacyManagedEnd}}

func ReplaceManagedBlockStrict(existing, newBlock string) (string, error) {
	start, end, found, err := managedBlockBounds(existing)
	if err != nil {
		return existing, err
	}
	lineEnding := detectLineEnding(existing)
	block := strings.ReplaceAll(strings.ReplaceAll(newBlock, "\r\n", "\n"), "\n", lineEnding)
	if found {
		return existing[:start] + block + existing[end:], nil
	}
	if existing == "" {
		return block, nil
	}
	separator := lineEnding + lineEnding
	if strings.HasSuffix(existing, lineEnding+lineEnding) {
		separator = ""
	} else if strings.HasSuffix(existing, lineEnding) {
		separator = lineEnding
	}
	return existing + separator + block, nil
}

func RemoveManagedBlockStrict(existing string) (string, error) {
	start, end, found, err := managedBlockBounds(existing)
	if err != nil || !found {
		return existing, err
	}
	return existing[:start] + existing[end:], nil
}

func managedBlockBounds(existing string) (int, int, bool, error) {
	foundStart, foundEnd := -1, -1
	for _, pair := range managedMarkerPairs {
		startCount, endCount := strings.Count(existing, pair.start), strings.Count(existing, pair.end)
		if startCount == 0 {
			continue
		}
		if foundStart >= 0 || startCount != 1 || endCount != 1 {
			return 0, 0, false, fmt.Errorf("malformed or duplicate managed RunX markers")
		}
		foundStart = strings.Index(existing, pair.start)
		endIndex := strings.Index(existing[foundStart+len(pair.start):], pair.end)
		if endIndex < 0 {
			return 0, 0, false, fmt.Errorf("managed RunX start marker has no matching end marker")
		}
		foundEnd = foundStart + len(pair.start) + endIndex + len(pair.end)
		if strings.HasPrefix(existing[foundEnd:], "\r\n") {
			foundEnd += 2
		} else if strings.HasPrefix(existing[foundEnd:], "\n") {
			foundEnd++
		}
	}
	allStarts := strings.Count(existing, ManagedStart) + strings.Count(existing, MojibakeManagedStart) + strings.Count(existing, LegacyManagedStart)
	allEnds := strings.Count(existing, ManagedEnd) + strings.Count(existing, LegacyManagedEnd)
	if foundStart < 0 {
		if allEnds != 0 {
			return 0, 0, false, fmt.Errorf("managed RunX end marker has no matching start marker")
		}
		return 0, 0, false, nil
	}
	if allStarts != 1 || allEnds != 1 {
		return 0, 0, false, fmt.Errorf("malformed or duplicate managed RunX markers")
	}
	return foundStart, foundEnd, true, nil
}

func detectLineEnding(value string) string {
	if strings.Contains(value, "\r\n") {
		return "\r\n"
	}
	return "\n"
}
