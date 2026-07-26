//go:build !windows

package update

import "os"

func replaceFileAtomic(source, destination string) error { return os.Rename(source, destination) }
