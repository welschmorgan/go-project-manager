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

type InitOptions struct {
	Bare bool
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
	Save             bool
	List             bool
	Apply            bool
	Pop              bool
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
	Delete    bool
	Annotated bool
	Message   string
	Commit    string
}

type ResetOptions struct {
	VersionControlOptions
	Hard   bool
	Commit string
}
type DeleteBranchOptions struct {
	VersionControlOptions
	Local      bool
	Remote     bool
	RemoteName string
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

	// Initialize a new repository
	Initialize(path string, options VersionControlOptions) error

	// Clone a remote repository
	Clone(url, path string, options VersionControlOptions) error

	// Get the working tree status (dirty / clean)
	Status(options VersionControlOptions) ([]string, error)

	// Get the name of the currently checked out branch
	CurrentBranch() (string, error)

	// Checkout a specific branch
	Checkout(branch string, options VersionControlOptions) error

	// Reset a branch to a specific commit
	Reset(options VersionControlOptions) error

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

	// Delete repository branch
	DeleteBranch(name string, options VersionControlOptions) error

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
func runCommand(name string, args ...string) (exitCode int, stdout []string, stderr []string, error error) {
	var bufStdout bytes.Buffer
	var bufStderr bytes.Buffer
	if config.Get().Verbose || config.Get().DryRun {
		argStr := ""
		for _, a := range args {
			if len(argStr) > 0 {
				argStr += ", "
			}
			argStr += fmt.Sprintf("%q", a)
		}
		fmt.Printf("%s* exec: %q %s\n", strings.Repeat("\t", config.Get().Indent), name, argStr)
	}
	ret := []string{}
	var errs []string
	if !config.Get().DryRun {
		cmd := exec.Command(name, args...)
		cmd.Stderr = &bufStderr
		cmd.Stdout = &bufStdout
		if err := cmd.Run(); err != nil {
			return cmd.ProcessState.ExitCode(), nil, strings.Split(bufStderr.String(), "\n"), err
		}
		exitCode = cmd.ProcessState.ExitCode()
		lines := map[string]bool{}
		for line, err := bufStdout.ReadString('\n'); err == nil; line, err = bufStdout.ReadString('\n') {
			line = strings.TrimSpace(line)
			if ok := lines[line]; !ok {
				lines[line] = true
				ret = append(ret, line)
			}
		}
		if len(strings.TrimSpace(bufStderr.String())) > 0 {
			errs = strings.Split(strings.TrimSpace(bufStderr.String()), "\n")
		}
	}
	return exitCode, ret, errs, nil
}

func dumpCommandErrors(exitCode int, errs []string) {
	level := ""
	color := ""
	indent := strings.Repeat("\t", config.Get().Indent)
	if exitCode != 0 {
		level = "error"
		color = "\033[1;31m"
	} else {
		level = "warning"
		color = "\033[1;33m"
	}
	shouldPrint := level == "error" || config.Get().Verbose
	if !shouldPrint {
		return
	}
	if len(errs) > 0 {
		if len(errs) == 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s%s: %v\n", indent, color, level, "\033[0m", errs[0])
		} else {
			errStr := ""
			numErrs := 0
			for _, err := range errs {
				if len(strings.TrimSpace(err)) > 0 {
					if len(errStr) > 0 {
						errStr += "\n"
					}
					errStr += fmt.Sprintf("%s\t- %s", indent, strings.TrimSpace(err))
					numErrs += 1
				}
			}
			fmt.Fprintf(os.Stderr, "%s%s%d %s(s)%s:\n%s\n", indent, color, len(errs), level, "\033[0m", errStr)
		}
	}
}

var All = []VersionControlSoftware{
	&Git{},
	&Svn{},
	&Hg{},
}

var AllNames = []string{}

func init() {
	for _, v := range All {
		AllNames = append(AllNames, v.Name())
	}
}

func instanciate(a VersionControlSoftware) VersionControlSoftware {
	inst := reflect.New(reflect.TypeOf(a).Elem())
	return inst.Interface().(VersionControlSoftware)
}

func Get(n string) VersionControlSoftware {
	for _, s := range All {
		if s.Name() == n {
			return instanciate(s)
		}
	}
	return nil
}

func Detect(path string) (VersionControlSoftware, error) {
	for _, s := range All {
		ok, err := s.Detect(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %s", path, err.Error())
		}
		if ok {
			return instanciate(s), nil
		}
	}
	return nil, fmt.Errorf("unknown vcs for folder '%s'", path)
}

func Open(path string) (VersionControlSoftware, error) {
	if vc, err := Detect(path); err != nil {
		return nil, err
	} else {
		if err := vc.Open(path); err != nil {
			return nil, err
		}
		return vc, nil
	}
}

func Initialize(n, p string, options VersionControlOptions) (VersionControlSoftware, error) {
	v := Get(n)
	return v, v.Initialize(p, options)
}
