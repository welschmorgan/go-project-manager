package vcs

import (
	"errors"
	"fmt"
	"os"

	"github.com/welschmorgan/go-project-manager/models"
)

var (
	errNotYetImpl = errors.New("not yet implemented")
)

type VersionControlSoftware interface {
	Name() string
	Path() string
	Detect(path string) (bool, error)
	Open(path string) error
	Clone(url, path string) error
	Checkout(branch string) error
	Pull() error
	Push() error
	Tag(name, commit, message string) error
	Merge(source, dest string) error
	Authors() ([]*models.Person, error)
}

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

var All = []VersionControlSoftware{
	&Git{},
	&Svn{},
	&Hg{},
}

func Get(n string) VersionControlSoftware {
	for _, s := range All {
		if s.Name() == n {
			return s
		}
	}
	return nil
}

func Open(path string) (VersionControlSoftware, error) {
	for _, s := range All {
		if ok, _ := s.Detect(path); ok {
			return s, s.Open(path)
		}
	}
	return nil, fmt.Errorf("unknown vcs for folder '%s'", path)
}
