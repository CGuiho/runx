package manifest

// Manifest represents the root structure of a runx.yaml manifest V2.
type Manifest struct {
	Version     string    `yaml:"version"`
	Namespace   string    `yaml:"namespace"`
	Description string    `yaml:"description,omitempty"`
	Parent      string    `yaml:"parent,omitempty"`
	Commands    []Command `yaml:"commands"`
}

// Command represents a single command leaf or group in a manifest.
type Command struct {
	UID         string    `yaml:"uid,omitempty"`
	ID          string    `yaml:"id,omitempty"`
	Group       string    `yaml:"group,omitempty"`
	Description string    `yaml:"description,omitempty"`
	Run         string    `yaml:"run,omitempty"`
	RunX        string    `yaml:"runx,omitempty"`
	CWD         string    `yaml:"cwd,omitempty"`
	Confirm     string    `yaml:"confirm,omitempty"`
	Commands    []Command `yaml:"commands,omitempty"`
}

// ResolvedCommand represents a flattened, executable command ready for execution.
type ResolvedCommand struct {
	UID         string `json:"uid"`
	FullID      string `json:"full_id"`
	Description string `json:"description"`
	Run         string `json:"run"`
	CWD         string `json:"cwd"`
	Confirm     string `json:"confirm"`
	CatalogPath string `json:"catalog_path"`
}
