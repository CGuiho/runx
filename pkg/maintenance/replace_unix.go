//go:build !windows

package maintenance

import "os"

func replaceFileAtomic(source, destination string) error { return os.Rename(source, destination) }
