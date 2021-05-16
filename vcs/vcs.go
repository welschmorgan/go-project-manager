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
	for _, s := range All {
		ok, err := s.Detect(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %s", path, err.Error())
		}
		if ok {
			return s, s.Open(path)
		}
	}
	return nil, fmt.Errorf("unknown vcs for folder '%s'", path)
}
