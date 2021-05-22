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
	StartingPoint    string
}

type PullOptions struct {
	VersionControlOptions
	Force    bool
	All      bool
	ListTags bool
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

type ListTagsOptions struct {
	VersionControlOptions
	SortByTaggerDate    bool
	SortByCommitterDate bool
}

type TagOptions struct {
	VersionControlOptions
	Annotated bool
	Message   string
	Commit    string
}

type VersionControlSoftware interface {
	// Retrieve the name of this vcs
	Name() string

	// Retrieve the path of the current repository
	Path() string

	// Retrieve the url of the current repository
	Url() string

	// Detect if a given path can be handled by this VCS
	Detect(path string) (bool, error)

	// Open a local repository, loading infos
	Open(path string) error

	// Clone a remote repository
	Clone(url, path string, options VersionControlOptions) error

	// Get the working tree status (dirty / clean)
	Status(options VersionControlOptions) ([]string, error)

	// Get the name of the currently checked out branch
	CurrentBranch() (string, error)

	// Checkout a specific branch
	Checkout(branch string, options VersionControlOptions) error

	// Pull sources from remote repository
	Pull(options VersionControlOptions) error

	// Push sources to remote repository
	Push(options VersionControlOptions) error

	// Create a new tag
	Tag(name string, options VersionControlOptions) error

	// Merge source into dest branch
	Merge(source, dest string, options VersionControlOptions) error

	// Create a new stash from the working tree
	Stash(options VersionControlOptions) ([]string, error)

	// List repository branches
	ListBranches(options VersionControlOptions) ([]string, error)

	// List repository authors (scan all commits)
	ListAuthors(options VersionControlOptions) ([]*config.Person, error)

	// List repository remote urls
	ListRemotes(options VersionControlOptions) (map[string]string, error)

	// List tags
	ListTags(options VersionControlOptions) ([]string, error)
}

// Run a command using os.exec. It returns the split stdout, potentially an error, and split stderr
func runCommand(name string, args ...string) ([]string, []string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if config.Get().Verbose || config.Get().DryRun {
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
