// Package config defines RunX's separate strict project and global
// configuration contracts, including the mandatory agent-evolution policy.
package config

import (
	"fmt"
)

// Policy values allowed for every agent.evolution field.
const (
	PolicyDisabled       = "disabled"
	PolicyAlwaysAsk      = "always-ask"
	PolicyAlwaysProceed  = "always-proceed"
	DefaultPolicy        = PolicyAlwaysAsk
	policySchemaEnum     = PolicyDisabled + ", " + PolicyAlwaysAsk + ", " + PolicyAlwaysProceed
	schemaBaseURLPattern = "https://github.com/CGuiho/runx/releases/download/v%s/%s"
)

// IssuesPolicy governs GitHub issue creation by AI agents.
type IssuesPolicy struct {
	Bugs         string `yaml:"bugs" json:"bugs"`
	Improvements string `yaml:"improvements" json:"improvements"`
	Reviews      string `yaml:"reviews" json:"reviews"`
}

// Evolution is the mandatory agent-evolution policy contract.
type Evolution struct {
	Upgrade string       `yaml:"upgrade" json:"upgrade"`
	Issues  IssuesPolicy `yaml:"issues" json:"issues"`
}

// Agent is the top-level agent policy namespace.
type Agent struct {
	Evolution Evolution `yaml:"evolution" json:"evolution"`
}

func normalizePolicy(value string) (string, error) {
	switch value {
	case "":
		return DefaultPolicy, nil
	case PolicyDisabled, PolicyAlwaysAsk, PolicyAlwaysProceed:
		return value, nil
	default:
		return "", fmt.Errorf("invalid agent.evolution value %q: must be one of %s", value, policySchemaEnum)
	}
}

func (e Evolution) validate() error {
	for _, field := range []struct {
		name  string
		value string
	}{
		{"agent.evolution.upgrade", e.Upgrade},
		{"agent.evolution.issues.bugs", e.Issues.Bugs},
		{"agent.evolution.issues.improvements", e.Issues.Improvements},
		{"agent.evolution.issues.reviews", e.Issues.Reviews},
	} {
		if _, err := normalizePolicy(field.value); err != nil {
			return fmt.Errorf("%s: %w", field.name, err)
		}
	}
	return nil
}

// ProjectConfig is the strict project-root `runx.yaml` policy overlay.
//
// The catalog fields belong to pkg/manifest; this type models only the
// convention-owned configuration surface so both contracts stay independent.
type ProjectConfig struct {
	Schema string `yaml:"schema,omitempty" json:"schema,omitempty"`
	Agent  Agent  `yaml:"agent,omitempty" json:"agent,omitempty"`
}

// GlobalConfig is the strict global `runx.global.yaml` baseline.
type GlobalConfig struct {
	Schema string `yaml:"schema,omitempty" json:"schema,omitempty"`
	Agent  Agent  `yaml:"agent,omitempty" json:"agent,omitempty"`
}

// Effective resolves the policy with project values overriding matching global
// values and every missing value inheriting the global value or the default.
type Effective struct {
	Upgrade string       `json:"upgrade"`
	Issues  IssuesPolicy `json:"issues"`
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

// resolveField applies project-over-global precedence to one policy field.
func resolveField(projectValue, globalValue, name string) (string, error) {
	value, source := projectValue, "project"
	if value == "" {
		value, source = globalValue, "global"
	}
	if value == "" {
		return DefaultPolicy, nil
	}
	policy, err := normalizePolicy(value)
	if err != nil {
		return "", fmt.Errorf("agent.evolution.%s (%s): %w", name, source, err)
	}
	return policy, nil
}

// Resolve overlays project over global and applies defaults. It fails when any
// explicit value violates the policy enum.
func Resolve(project ProjectConfig, global GlobalConfig) (Effective, error) {
	resolved := Effective{}
	fields := []struct {
		name         string
		projectValue string
		globalValue  string
		set          func(string)
	}{
		{"upgrade", project.Agent.Evolution.Upgrade, global.Agent.Evolution.Upgrade, func(v string) { resolved.Upgrade = v }},
		{"issues.bugs", project.Agent.Evolution.Issues.Bugs, global.Agent.Evolution.Issues.Bugs, func(v string) { resolved.Issues.Bugs = v }},
		{"issues.improvements", project.Agent.Evolution.Issues.Improvements, global.Agent.Evolution.Issues.Improvements, func(v string) { resolved.Issues.Improvements = v }},
		{"issues.reviews", project.Agent.Evolution.Issues.Reviews, global.Agent.Evolution.Issues.Reviews, func(v string) { resolved.Issues.Reviews = v }},
	}
	for _, field := range fields {
		policy, err := resolveField(field.projectValue, field.globalValue, field.name)
		if err != nil {
			return Effective{}, err
		}
		field.set(policy)
	}
	return resolved, nil
}

// Validate checks a decoded document against its semantic contract.
func Validate(document any) error {
	switch typed := document.(type) {
	case ProjectConfig:
		return typed.Agent.Evolution.validate()
	case *ProjectConfig:
		if typed == nil {
			return fmt.Errorf("nil project configuration")
		}
		return typed.Agent.Evolution.validate()
	case GlobalConfig:
		return typed.Agent.Evolution.validate()
	case *GlobalConfig:
		if typed == nil {
			return fmt.Errorf("nil global configuration")
		}
		return typed.Agent.Evolution.validate()
	default:
		return fmt.Errorf("unsupported configuration document type %T", document)
	}
}

// SchemaURL returns the version-pinned remote schema reference for a release.
func SchemaURL(version, schemaFile string) string {
	return fmt.Sprintf(schemaBaseURLPattern, version, schemaFile)
}

const (
	// ProjectSchemaFile is the project configuration schema asset name.
	ProjectSchemaFile = "runx.schema.json"
	// GlobalSchemaFile is the global configuration schema asset name.
	GlobalSchemaFile = "runx.global.schema.json"
)
