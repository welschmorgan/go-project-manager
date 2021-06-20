package fs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type DirStackEntry struct {
	Previous string
	Next     string
}

var dirStack []DirStackEntry = []DirStackEntry{}

func Pushd(newDir string) (string, error) {
	if len(strings.TrimSpace(newDir)) == 0 {
		panic("Cannot pushd to empty directory")
	}
	nextIsCurrentDir := false
	if oldCwd, err := os.Getwd(); err != nil {
		return "", err
	} else {
		if nextIsCurrentDir = oldCwd == newDir; !nextIsCurrentDir {
			dirStack = append(dirStack, DirStackEntry{oldCwd, newDir})
		}
	}
	if !nextIsCurrentDir {
		if err := os.Chdir(newDir); err != nil {
			return "", err
		}
	}
	return newDir, nil
}

func DumpDirStack(w io.Writer) error {
	if _, err := w.Write([]byte("DirStack:\n")); err != nil {
		return err
	}
	for i, e := range dirStack {
		if _, err := fmt.Fprintf(w, "\t[%d] %s -> %s\n", i, e.Previous, e.Next); err != nil {
			return err
		}
	}
	return nil
}

func Popd() (string, error) {
	if len(dirStack) == 0 {
		return "", errors.New("no directory in stack")
	}
	dirEntry := dirStack[len(dirStack)-1]
	if err := os.Chdir(dirEntry.Previous); err != nil {
		return "", err
	}
	dirStack = dirStack[0 : len(dirStack)-1]
	return dirEntry.Previous, nil
}
