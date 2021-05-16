package vcs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
)

var (
	errNotYetImpl = errors.New("not yet implemented")
)

type VersionControlSoftware interface {
	Name() string
	Path() string
	Url() string
	Detect(path string) (bool, error)
	Open(path string) error
	Clone(url, path string) error
	Checkout(branch string) error
	Pull() error
	Push() error
	Tag(name, commit, message string) error
	Merge(source, dest string) error
	Authors() ([]*models.Person, error)
	Remotes() (map[string]string, error)
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

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func runCommand(name string, args ...string) ([]string, error, []string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err, strings.Split(stderr.String(), "\n")
	}
	ret := []string{}
	lines := map[string]bool{}
	for line, err := stdout.ReadString('\n'); err == nil; line, err = stdout.ReadString('\n') {
		line = strings.TrimSpace(line)
		if ok := lines[line]; !ok {
			lines[line] = true
			ret = append(ret, line)
		}
	}
	return ret, nil, strings.Split(stderr.String(), "\n")
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
	println("detecting vcs for", path)
	for _, s := range All {
		println("\ttrying", s.Name())
		ok, err := s.Detect(path)
		if err != nil {
			println("\t\t-> error:", err.Error())
		}
		if ok {
			println("\t\t-> ok!")
		} else {
			println("\t\t-> not ok!")
		}
		if ok {
			return s, s.Open(path)
		}
	}
	return nil, fmt.Errorf("unknown vcs for folder '%s'", path)
}
