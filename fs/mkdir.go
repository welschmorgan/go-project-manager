package fs

import "os"

// Mkdir creates all elements towards a directory.
// If the directory already exists, it exists gracefully
// It returns an error if the directory creation failed
func Mkdir(p string) error {
	if _, err := os.Stat(p); err == nil || os.IsExist(err) {
		return err
	}
	return os.MkdirAll(p, 0755)
}
