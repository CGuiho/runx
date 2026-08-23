package installstate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ReadPointer loads and strictly validates the current pointer. A missing
// pointer returns (nil, nil) so callers can distinguish "not installed".
func ReadPointer() (*Pointer, error) {
	path, err := CurrentPointerPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pointer %s: %w", path, err)
	}
	var pointer Pointer
	decoder := json.NewDecoder(bytesReader(withoutUTF8BOM(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pointer); err != nil {
		return nil, fmt.Errorf("decode pointer %s: %w", path, err)
	}
	if err := pointer.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pointer %s: %w", path, err)
	}
	return &pointer, nil
}

// withoutUTF8BOM accepts pointers produced by legacy Windows PowerShell
// installers while keeping the JSON decoder and pointer schema strict.
func withoutUTF8BOM(data []byte) []byte {
	if len(data) >= 3 && data[0] == 0xef && data[1] == 0xbb && data[2] == 0xbf {
		return data[3:]
	}
	return data
}

// WritePointer atomically replaces the current pointer.
func WritePointer(pointer Pointer) error {
	if err := pointer.Validate(); err != nil {
		return err
	}
	path, err := CurrentPointerPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(pointer, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return WriteFileAtomic(path, data, 0o644)
}

// Artifact is one entry of the installed-artifacts ledger: the authority for
// what this installation owns and where projections live.
type Artifact struct {
	ID          string   `json:"id"`
	Version     string   `json:"version"`
	Path        string   `json:"path"`
	SHA256      string   `json:"sha256,omitempty"`
	Kind        string   `json:"kind"`
	Projections []string `json:"projections,omitempty"`
}

// Ledger is the installed ownership manifest stored at
// $HOME/.guiho/runx/installed-artifacts.json.
type Ledger struct {
	Protocol  int        `json:"protocol"`
	Version   string     `json:"version"`
	Artifacts []Artifact `json:"artifacts"`
}

// Validate enforces ledger invariants: protocol, version, unique IDs, and
// owned paths confined to the CLI home or the single launcher path.
func (l Ledger) Validate() error {
	if l.Protocol != ProtocolVersion {
		return fmt.Errorf("unsupported ledger protocol %d: want %d", l.Protocol, ProtocolVersion)
	}
	if !isValidVersionName(l.Version) {
		return fmt.Errorf("ledger version %q is not a valid SemVer", l.Version)
	}
	cliDir, err := CLIDir()
	if err != nil {
		return err
	}
	launcherPath, err := LauncherPath()
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	for _, artifact := range l.Artifacts {
		if artifact.ID == "" {
			return fmt.Errorf("ledger artifact with empty ID")
		}
		if seen[artifact.ID] {
			return fmt.Errorf("duplicate ledger artifact ID %q", artifact.ID)
		}
		seen[artifact.ID] = true
		if !pathOwnedBy(artifact.Path, cliDir) && filepath.Clean(artifact.Path) != filepath.Clean(launcherPath) {
			return fmt.Errorf("ledger artifact %q path %q is outside CLI-owned roots", artifact.ID, artifact.Path)
		}
	}
	return nil
}

func pathOwnedBy(path, root string) bool {
	cleanPath := filepath.Clean(path)
	cleanRoot := filepath.Clean(root)
	rel, err := filepath.Rel(cleanRoot, cleanPath)
	if err != nil {
		return false
	}
	return rel != ".." && !filepath.IsAbs(rel) && !hasTraversalPrefix(rel)
}

func hasTraversalPrefix(rel string) bool {
	return rel == ".." || len(rel) >= 3 && rel[:3] == ".."+string(filepath.Separator)
}

// ReadLedger loads the installed ledger. A missing file returns (nil, nil).
func ReadLedger() (*Ledger, error) {
	path, err := InstalledLedgerPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read ledger %s: %w", path, err)
	}
	var ledger Ledger
	decoder := json.NewDecoder(bytesReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&ledger); err != nil {
		return nil, fmt.Errorf("decode ledger %s: %w", path, err)
	}
	if err := ledger.Validate(); err != nil {
		return nil, fmt.Errorf("invalid ledger %s: %w", path, err)
	}
	return &ledger, nil
}

// WriteLedger atomically replaces the installed ledger after validation.
func WriteLedger(ledger Ledger) error {
	if err := ledger.Validate(); err != nil {
		return err
	}
	path, err := InstalledLedgerPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(ledger, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return WriteFileAtomic(path, data, 0o644)
}

// ReadPointerIn loads and validates the pointer under an explicit home.
func ReadPointerIn(home string) (*Pointer, error) {
	data, err := os.ReadFile(CurrentPointerPathIn(home))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read pointer: %w", err)
	}
	var pointer Pointer
	decoder := json.NewDecoder(bytesReader(withoutUTF8BOM(data)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&pointer); err != nil {
		return nil, fmt.Errorf("decode pointer: %w", err)
	}
	if err := pointer.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pointer: %w", err)
	}
	return &pointer, nil
}

// WritePointerIn atomically replaces the pointer under an explicit home.
func WritePointerIn(home string, pointer Pointer) error {
	if err := pointer.Validate(); err != nil {
		return err
	}
	data, err := json.MarshalIndent(pointer, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return WriteFileAtomic(CurrentPointerPathIn(home), data, 0o644)
}
