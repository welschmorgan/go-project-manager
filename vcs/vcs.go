package vcs

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
)

var (
	errNotYetImpl = errors.New("not yet implemented")
)

type VersionControlOptions interface{}

func getOptions(options, defaults VersionControlOptions) (VersionControlOptions, error) {
	if options == nil {
		return defaults, nil
	}
	optType := reflect.TypeOf(options)
	defType := reflect.TypeOf(defaults)
	if defType.Name() != optType.Name() {
		return nil, fmt.Errorf("options are of wrong type, expected %s but got %s", defType.Name(), optType.Name())
	}
	return options, nil
}

type CloneOptions struct {
	Branch   string
	Insecure bool
}
type CheckoutOptions struct {
	VersionControlOptions
	CreateBranch     bool
	UpdateIfExisting bool
}

type PullOptions struct {
	VersionControlOptions
	Force bool
	All   bool
	Tags  bool
}

type PushOptions struct {
	VersionControlOptions
	Force bool
	All   bool
}

type MergeOptions struct {
	VersionControlOptions
	NoFastForward   bool
	FastForwardOnly bool
}

type StatusOptions struct {
	VersionControlOptions
	Short bool
}

type StashOptions struct {
	VersionControlOptions
	IncludeUntracked bool
	Message          string
}

type BranchOptions struct {
	VersionControlOptions
	All           bool
	Verbose       bool
	SetUpstreamTo string
}

type VersionControlSoftware interface {
	Name() string
	Path() string
	Url() string
	Detect(path string) (bool, error)
	Open(path string) error
	Clone(url, path string, options VersionControlOptions) error
	Status(options VersionControlOptions) ([]string, error)
	Branch(options VersionControlOptions) ([]string, error)
	Checkout(branch string, options VersionControlOptions) error
	Pull(options VersionControlOptions) error
	Push(options VersionControlOptions) error
	Tag(name, commit, message string, options VersionControlOptions) error
	Merge(source, dest string, options VersionControlOptions) error
	Stash(options VersionControlOptions) ([]string, error)
	Authors(options VersionControlOptions) ([]*config.Person, error)
	Remotes(options VersionControlOptions) (map[string]string, error)
}

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func runCommand(name string, args ...string) ([]string, []string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if config.Get().Verbose {
		argStr := ""
		for _, a := range args {
			if len(argStr) > 0 {
				argStr += ", "
			}
			argStr += fmt.Sprintf("%q", a)
		}
		fmt.Printf("* exec: %q %s\n", name, argStr)
	}
	ret := []string{}
	var errs []string
	if !config.Get().DryRun {
		cmd := exec.Command(name, args...)
		cmd.Stderr = &stderr
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			return nil, strings.Split(stderr.String(), "\n"), err
		}
		lines := map[string]bool{}
		for line, err := stdout.ReadString('\n'); err == nil; line, err = stdout.ReadString('\n') {
			line = strings.TrimSpace(line)
			if ok := lines[line]; !ok {
				lines[line] = true
				ret = append(ret, line)
			}
		}
		if len(strings.TrimSpace(stderr.String())) > 0 {
			errs = strings.Split(strings.TrimSpace(stderr.String()), "\n")
		}
	}
	return ret, errs, nil
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
