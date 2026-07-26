package manifest

// Manifest is the strict on-disk RunX manifest v2 contract.
type Manifest struct {
	Version   string    `yaml:"version" json:"version"`
	Namespace string    `yaml:"namespace" json:"namespace"`
	Scripts   Scripts   `yaml:"scripts" json:"scripts"`
	Parent    string    `yaml:"parent,omitempty" json:"parent,omitempty"`
	Commands  []Command `yaml:"commands" json:"commands"`
}

type Scripts struct {
	Directory string `yaml:"directory" json:"directory"`
}

// Command represents either one executable leaf or one group. Semantic
// validation enforces the exact shape for the selected kind.
type Command struct {
	UID         string    `yaml:"uid,omitempty" json:"uid,omitempty"`
	ID          string    `yaml:"id,omitempty" json:"id,omitempty"`
	Group       string    `yaml:"group,omitempty" json:"group,omitempty"`
	Summary     string    `yaml:"summary,omitempty" json:"summary,omitempty"`
	Description string    `yaml:"description,omitempty" json:"description,omitempty"`
	Command     string    `yaml:"command,omitempty" json:"command,omitempty"`
	CWD         string    `yaml:"cwd,omitempty" json:"cwd,omitempty"`
	Shell       string    `yaml:"shell,omitempty" json:"shell,omitempty"`
	Tags        []string  `yaml:"tags,omitempty" json:"tags,omitempty"`
	Confirm     string    `yaml:"confirm,omitempty" json:"confirm,omitempty"`
	Commands    []Command `yaml:"commands,omitempty" json:"commands,omitempty"`
	RunX        string    `yaml:"runx,omitempty" json:"runx,omitempty"`
}

type CatalogSource string

const (
	SourceLocal   CatalogSource = "local"
	SourceForeign CatalogSource = "foreign"
)

type ResolvedCommand struct {
	Index         int           `json:"index"`
	UID           string        `json:"uid"`
	ID            string        `json:"id"`
	Selector      string        `json:"selector"`
	Summary       string        `json:"summary"`
	Description   string        `json:"description"`
	Command       string        `json:"command"`
	CWD           string        `json:"cwd"`
	Shell         string        `json:"shell"`
	Tags          []string      `json:"tags"`
	Confirm       string        `json:"confirm"`
	CatalogPath   string        `json:"catalogPath"`
	CatalogSource CatalogSource `json:"catalogSource"`
}

type Group struct {
	Selector      string        `json:"selector"`
	Summary       string        `json:"summary"`
	CatalogPath   string        `json:"catalogPath"`
	CatalogSource CatalogSource `json:"catalogSource"`
}

type Child struct {
	Namespace         string        `json:"namespace"`
	DeclaredNamespace string        `json:"declaredNamespace"`
	Path              string        `json:"path"`
	Source            CatalogSource `json:"source"`
	Parent            string        `json:"parent"`
	DeclaredParent    string        `json:"declaredParent"`
}

type Catalog struct {
	Version   string            `json:"version"`
	Namespace string            `json:"namespace"`
	Scripts   Scripts           `json:"scripts"`
	Parent    string            `json:"parent,omitempty"`
	Path      string            `json:"path"`
	Commands  []ResolvedCommand `json:"commands"`
	Groups    []Group           `json:"groups"`
	Children  []Child           `json:"children"`
	lookup    map[string]int
}

func (catalog *Catalog) Resolve(selector string) (ResolvedCommand, bool) {
	index, ok := catalog.lookup[selector]
	if !ok {
		return ResolvedCommand{}, false
	}
	return catalog.Commands[index], true
}
