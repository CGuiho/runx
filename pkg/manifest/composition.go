package manifest

import (
	"fmt"
)

// IndexManifest flattens the Manifest V2 command hierarchy into a map of selector/UID to ResolvedCommand.
func IndexManifest(m *Manifest, catalogPath string) (map[string]ResolvedCommand, error) {
	resolved := make(map[string]ResolvedCommand)
	if err := traverseCommands(m.Commands, m.Namespace, catalogPath, resolved, 1); err != nil {
		return nil, err
	}
	return resolved, nil
}

func traverseCommands(cmds []Command, prefix string, catalogPath string, out map[string]ResolvedCommand, depth int) error {
	if depth > 32 {
		return fmt.Errorf("exceeded maximum nested group depth of 32 levels")
	}

	for _, c := range cmds {
		if c.Group != "" {
			groupPrefix := prefix + ":" + c.Group
			if len(c.Commands) > 0 {
				if err := traverseCommands(c.Commands, groupPrefix, catalogPath, out, depth+1); err != nil {
					return err
				}
			}
			continue
		}

		if c.ID == "" {
			return fmt.Errorf("command leaf requires an 'id' field")
		}

		selector := prefix + ":" + c.ID
		uid := c.UID
		if uid == "" {
			uid = selector
		}

		if _, exists := out[selector]; exists {
			return fmt.Errorf("duplicate command selector '%s' detected", selector)
		}

		cmdObj := ResolvedCommand{
			UID:         uid,
			FullID:      selector,
			Description: c.Description,
			Run:         c.Run,
			CWD:         c.CWD,
			Confirm:     c.Confirm,
			CatalogPath: catalogPath,
		}

		out[selector] = cmdObj
		if uid != selector {
			out[uid] = cmdObj
		}
	}

	return nil
}
