package fs

import (
	"errors"
	"os"
)

var dirStack []string = []string{}

func Pushd(newDir string) (string, error) {
	if oldCwd, err := os.Getwd(); err != nil {
		return "", err
	} else {
		dirStack = append(dirStack, oldCwd)
	}
	if err := os.Chdir(newDir); err != nil {
		return "", err
	}
	return newDir, nil
}

func Popd() (string, error) {
	if len(dirStack) == 0 {
		return "", errors.New("no directory in stack")
	}
	newDir := dirStack[len(dirStack)-1]
	if err := os.Chdir(newDir); err != nil {
		return "", err
	}
	dirStack = dirStack[0 : len(dirStack)-1]
	return newDir, nil
}
