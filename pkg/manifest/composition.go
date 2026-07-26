package manifest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	maximumCatalogDepth = 32
	maximumCatalogBytes = 1024 * 1024
)

type LoadOptions struct {
	CWD        string
	ConfigPath string
	HomeDir    string
	HTTPClient *http.Client
}

type loadedCatalog struct {
	manifest Manifest
	path     string
	source   CatalogSource
	basePath string
}

type loadState struct {
	active     map[string]bool
	commands   []ResolvedCommand
	groups     []Group
	children   []Child
	identities map[string]int
	idOwners   map[string][]int
	client     *http.Client
}

func ResolveConfigPath(cwd, explicit, home string) (string, error) {
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("resolve working directory: %w", err)
		}
	}
	absCWD, err := filepath.Abs(cwd)
	if err != nil {
		return "", fmt.Errorf("resolve working directory: %w", err)
	}
	var candidates []string
	if explicit != "" {
		if filepath.IsAbs(explicit) {
			candidates = []string{explicit}
		} else {
			candidates = []string{filepath.Join(absCWD, explicit)}
		}
	} else {
		candidates = append(candidates, filepath.Join(absCWD, "runx.yaml"))
		if home == "" {
			home, _ = os.UserHomeDir()
		}
		if home != "" {
			candidates = append(candidates, filepath.Join(home, ".guiho", "runx", "runx.yaml"))
		}
	}
	for _, candidate := range candidates {
		absolute, err := filepath.Abs(candidate)
		if err == nil {
			if info, statErr := os.Stat(absolute); statErr == nil && !info.IsDir() {
				return absolute, nil
			}
		}
	}
	return "", fmt.Errorf("no runx.yaml found; checked %s", strings.Join(candidates, ", "))
}

func Load(ctx context.Context, opts LoadOptions) (*Catalog, error) {
	path, err := ResolveConfigPath(opts.CWD, opts.ConfigPath, opts.HomeDir)
	if err != nil {
		return nil, err
	}
	root, err := loadLocal(path)
	if err != nil {
		return nil, err
	}
	client := opts.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	state := &loadState{
		active: map[string]bool{}, identities: map[string]int{}, idOwners: map[string][]int{}, client: client,
	}
	if root.manifest.Parent != "" {
		if err := validateDeclaredParent(ctx, root, state); err != nil {
			return nil, err
		}
	}
	if err := expand(ctx, root, nil, state, 0); err != nil {
		return nil, err
	}

	lookup := map[string]int{}
	for identity, index := range state.identities {
		lookup[identity] = index
	}
	for id, owners := range state.idOwners {
		if len(owners) == 1 {
			lookup[id] = owners[0]
		}
	}
	for index := range state.commands {
		state.commands[index].Index = index + 1
	}
	return &Catalog{
		Version: root.manifest.Version, Namespace: root.manifest.Namespace, Scripts: root.manifest.Scripts,
		Parent: root.manifest.Parent, Path: root.path, Commands: state.commands, Groups: state.groups,
		Children: state.children, lookup: lookup,
	}, nil
}

func expand(ctx context.Context, catalog loadedCatalog, prefix []string, state *loadState, depth int) error {
	if depth > maximumCatalogDepth {
		return fmt.Errorf("catalog composition exceeds maximum depth of %d", maximumCatalogDepth)
	}
	if state.active[catalog.path] {
		return fmt.Errorf("catalog composition cycle detected at %s", catalog.path)
	}
	state.active[catalog.path] = true
	defer delete(state.active, catalog.path)

	for _, entry := range catalog.manifest.Commands {
		if entry.Group != "" {
			groupPath := append(append([]string{}, prefix...), entry.Group)
			selector := strings.Join(groupPath, "/")
			state.groups = append(state.groups, Group{Selector: selector, Summary: entry.Summary, CatalogPath: catalog.path, CatalogSource: catalog.source})
			if entry.Commands != nil {
				childCatalog := catalog
				childCatalog.manifest.Commands = entry.Commands
				if err := expandEntries(ctx, childCatalog, groupPath, state, depth+1); err != nil {
					return err
				}
				continue
			}
			child, err := loadReference(ctx, catalog, entry.RunX, state.client)
			if err != nil {
				return err
			}
			if child.manifest.Parent == "" {
				return fmt.Errorf("child catalog %s must declare its parent", child.path)
			}
			declared, _, err := resolveReference(child, child.manifest.Parent)
			if err != nil {
				return err
			}
			if declared != catalog.path {
				if !(child.source == SourceForeign && catalog.source == SourceLocal) {
					return fmt.Errorf("child catalog %s declares parent %s; expected %s", child.path, declared, catalog.path)
				}
				if err := validateDeclaredParent(ctx, child, state); err != nil {
					return err
				}
			}
			state.children = append(state.children, Child{Namespace: entry.Group, DeclaredNamespace: child.manifest.Namespace, Path: child.path, Source: child.source, Parent: catalog.path, DeclaredParent: declared})
			if err := expand(ctx, child, groupPath, state, depth+1); err != nil {
				return err
			}
			continue
		}
		if err := registerCommand(catalog, entry, prefix, state); err != nil {
			return err
		}
	}
	return nil
}

func expandEntries(ctx context.Context, catalog loadedCatalog, prefix []string, state *loadState, depth int) error {
	if depth > maximumCatalogDepth {
		return fmt.Errorf("catalog composition exceeds maximum depth of %d", maximumCatalogDepth)
	}
	for _, entry := range catalog.manifest.Commands {
		if entry.Group != "" {
			groupPath := append(append([]string{}, prefix...), entry.Group)
			selector := strings.Join(groupPath, "/")
			state.groups = append(state.groups, Group{Selector: selector, Summary: entry.Summary, CatalogPath: catalog.path, CatalogSource: catalog.source})
			if entry.Commands != nil {
				nested := catalog
				nested.manifest.Commands = entry.Commands
				if err := expandEntries(ctx, nested, groupPath, state, depth+1); err != nil {
					return err
				}
			} else {
				child, err := loadReference(ctx, catalog, entry.RunX, state.client)
				if err != nil {
					return err
				}
				if child.manifest.Parent == "" {
					return fmt.Errorf("child catalog %s must declare its parent", child.path)
				}
				declared, _, err := resolveReference(child, child.manifest.Parent)
				if err != nil {
					return err
				}
				if declared != catalog.path {
					return fmt.Errorf("child catalog %s declares parent %s; expected %s", child.path, declared, catalog.path)
				}
				state.children = append(state.children, Child{Namespace: entry.Group, DeclaredNamespace: child.manifest.Namespace, Path: child.path, Source: child.source, Parent: catalog.path, DeclaredParent: declared})
				if err := expand(ctx, child, groupPath, state, depth+1); err != nil {
					return err
				}
			}
			continue
		}
		if err := registerCommand(catalog, entry, prefix, state); err != nil {
			return err
		}
	}
	return nil
}

func registerCommand(catalog loadedCatalog, entry Command, prefix []string, state *loadState) error {
	selector := strings.Join(append(append([]string{}, prefix...), entry.ID), "/")
	base := catalog.basePath
	cwd := filepath.Clean(filepath.Join(base, entry.CWD))
	if entry.CWD == "" {
		cwd = base
	}
	relative, err := filepath.Rel(base, cwd)
	if err != nil || filepath.IsAbs(relative) || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("command %s has a cwd outside its catalog directory", entry.UID)
	}
	index := len(state.commands)
	owner := selector
	for _, identity := range []string{entry.UID, selector} {
		if prior, ok := state.identities[identity]; ok && prior != index {
			return fmt.Errorf("command identity %q conflicts with another command", identity)
		}
		if owners := state.idOwners[identity]; len(owners) > 0 {
			return fmt.Errorf("command identity %q conflicts with a command ID shorthand", identity)
		}
		state.identities[identity] = index
	}
	_ = owner
	state.idOwners[entry.ID] = append(state.idOwners[entry.ID], index)
	shell := entry.Shell
	if shell == "" {
		shell = "auto"
	}
	confirm := entry.Confirm
	if confirm == "" {
		confirm = "never"
	}
	state.commands = append(state.commands, ResolvedCommand{
		UID: entry.UID, ID: entry.ID, Selector: selector, Summary: entry.Summary, Description: entry.Description,
		Command: entry.Command, CWD: cwd, Shell: shell, Tags: append([]string{}, entry.Tags...), Confirm: confirm,
		CatalogPath: catalog.path, CatalogSource: catalog.source,
	})
	return nil
}

func loadLocal(path string) (loadedCatalog, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return loadedCatalog{}, err
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return loadedCatalog{}, fmt.Errorf("read manifest %q: %w", absolute, err)
	}
	manifest, err := ParseManifestBytes(data)
	if err != nil {
		return loadedCatalog{}, fmt.Errorf("invalid RunX configuration %s: %w", absolute, err)
	}
	return loadedCatalog{manifest: *manifest, path: absolute, source: SourceLocal, basePath: filepath.Dir(absolute)}, nil
}

func loadReference(ctx context.Context, owner loadedCatalog, reference string, client *http.Client) (loadedCatalog, error) {
	resolved, source, err := resolveReference(owner, reference)
	if err != nil {
		return loadedCatalog{}, err
	}
	if source == SourceLocal {
		return loadLocal(resolved)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, resolved, nil)
	if err != nil {
		return loadedCatalog{}, err
	}
	req.Header.Set("Accept", "text/yaml, text/plain")
	response, err := client.Do(req)
	if err != nil {
		return loadedCatalog{}, fmt.Errorf("load foreign RunX catalog %s: %w", resolved, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return loadedCatalog{}, fmt.Errorf("load foreign RunX catalog %s: HTTP %d", resolved, response.StatusCode)
	}
	if response.ContentLength > maximumCatalogBytes {
		return loadedCatalog{}, fmt.Errorf("foreign RunX catalog exceeds %d bytes: %s", maximumCatalogBytes, resolved)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, maximumCatalogBytes+1))
	if err != nil {
		return loadedCatalog{}, fmt.Errorf("read foreign RunX catalog %s: %w", resolved, err)
	}
	if len(data) > maximumCatalogBytes {
		return loadedCatalog{}, fmt.Errorf("foreign RunX catalog exceeds %d bytes: %s", maximumCatalogBytes, resolved)
	}
	manifest, err := ParseManifestBytes(data)
	if err != nil {
		return loadedCatalog{}, fmt.Errorf("invalid RunX configuration %s: %w", resolved, err)
	}
	return loadedCatalog{manifest: *manifest, path: resolved, source: SourceForeign, basePath: owner.basePath}, nil
}

func resolveReference(owner loadedCatalog, reference string) (string, CatalogSource, error) {
	if strings.HasPrefix(strings.ToLower(reference), "https://") {
		normalized, err := normalizeGitHubURL(reference)
		return normalized, SourceForeign, err
	}
	if filepath.IsAbs(reference) {
		return "", "", fmt.Errorf("RunX catalog references must be relative paths or full GitHub URLs: %s", reference)
	}
	if owner.source == SourceForeign {
		base, _ := url.Parse(owner.path)
		relative, err := url.Parse(reference)
		if err != nil {
			return "", "", err
		}
		resolved, err := normalizeGitHubURL(base.ResolveReference(relative).String())
		if err != nil {
			return "", "", err
		}
		if githubSourceRoot(resolved) != githubSourceRoot(owner.path) {
			return "", "", fmt.Errorf("relative foreign RunX catalog reference escapes its GitHub owner/repository/ref root: %s", reference)
		}
		return resolved, SourceForeign, nil
	}
	return filepath.Clean(filepath.Join(filepath.Dir(owner.path), reference)), SourceLocal, nil
}

func normalizeGitHubURL(value string) (string, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme != "https" {
		return "", fmt.Errorf("foreign RunX catalog URL must use HTTPS: %s", value)
	}
	if parsed.Host == "github.com" {
		parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
		if len(parts) < 5 || parts[2] != "blob" {
			return "", fmt.Errorf("GitHub RunX catalog URL must use /owner/repository/blob/ref/path: %s", value)
		}
		return "https://raw.githubusercontent.com/" + strings.Join(append(parts[:2], parts[3:]...), "/"), nil
	}
	if parsed.Host != "raw.githubusercontent.com" {
		return "", fmt.Errorf("foreign RunX catalogs must use a full GitHub URL: %s", value)
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("invalid raw GitHub RunX catalog URL: %s", value)
	}
	parsed.RawQuery, parsed.Fragment = "", ""
	return parsed.String(), nil
}

func githubSourceRoot(value string) string {
	parsed, _ := url.Parse(value)
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) > 3 {
		parts = parts[:3]
	}
	return parsed.Scheme + "://" + parsed.Host + "/" + strings.Join(parts, "/")
}

func validateDeclaredParent(ctx context.Context, child loadedCatalog, state *loadState) error {
	parentPath, source, err := resolveReference(child, child.manifest.Parent)
	if err != nil {
		return err
	}
	var parent loadedCatalog
	if source == SourceLocal {
		parent, err = loadLocal(parentPath)
	} else {
		parent, err = loadReference(ctx, child, child.manifest.Parent, state.client)
	}
	if err != nil {
		return err
	}
	var references []string
	collectReferences(parent.manifest.Commands, &references)
	for _, reference := range references {
		resolved, _, resolveErr := resolveReference(parent, reference)
		if resolveErr == nil && resolved == child.path {
			validation := &loadState{active: map[string]bool{}, identities: map[string]int{}, idOwners: map[string][]int{}, client: state.client}
			return expand(ctx, parent, nil, validation, 0)
		}
	}
	return fmt.Errorf("parent catalog %s does not declare child %s", parent.path, child.path)
}

func collectReferences(entries []Command, output *[]string) {
	for _, entry := range entries {
		if entry.RunX != "" {
			*output = append(*output, entry.RunX)
		}
		collectReferences(entry.Commands, output)
	}
}

// IndexManifest retains the small package API used by older callers while
// applying the v2 validation and selector rules to a single local catalog.
func IndexManifest(manifest *Manifest, catalogPath string) (map[string]ResolvedCommand, error) {
	base, _ := filepath.Abs(catalogPath)
	state := &loadState{active: map[string]bool{}, identities: map[string]int{}, idOwners: map[string][]int{}}
	loaded := loadedCatalog{manifest: *manifest, path: base, source: SourceLocal, basePath: filepath.Dir(base)}
	if err := expand(context.Background(), loaded, nil, state, 0); err != nil {
		return nil, err
	}
	result := map[string]ResolvedCommand{}
	for index := range state.commands {
		state.commands[index].Index = index + 1
		result[state.commands[index].Selector] = state.commands[index]
		result[state.commands[index].UID] = state.commands[index]
	}
	for id, owners := range state.idOwners {
		if len(owners) == 1 {
			result[id] = state.commands[owners[0]]
		}
	}
	return result, nil
}

func SortedCommands(commands []ResolvedCommand) []ResolvedCommand {
	result := append([]ResolvedCommand{}, commands...)
	sort.SliceStable(result, func(i, j int) bool { return result[i].Index < result[j].Index })
	return result
}
