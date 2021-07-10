package vcs

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
)

var (
	errNotYetImpl = errors.New("not yet implemented")
)

type VersionControlSoftware interface {
	// Retrieve the name of this vcs
	Name() string

	// Retrieve the path of the current repository
	Path() fs.Path

	// Retrieve the url of the current repository
	Url() string

	// Detect if a given path can be handled by this VCS
	Detect(path fs.Path) error

	// Open a local repository, loading infos
	Open(path fs.Path) error

	// Initialize a new repository
	Initialize(path fs.Path, options VersionControlOptions) error

	// Clone a remote repository
	Clone(url string, path fs.Path, options VersionControlOptions) error

	// Fetch remote index
	FetchIndex(options VersionControlOptions) error

	// Add files to index
	Stage(options VersionControlOptions) error

	// Retrieve commits without parents
	RootCommits() ([]string, error)

	// Create a new commit
	Commit(options VersionControlOptions) error

	// Get the working tree status (dirty / clean)
	Status(options VersionControlOptions) ([]string, error)

	// Get the name of the currently checked out branch
	CurrentBranch() (string, error)

	// Get the hash of the current commit
	CurrentCommit(options VersionControlOptions) (hash, subject string, err error)

	// Get the hash of the current commit
	ExtractLog(options VersionControlOptions) (lines []string, err error)

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

	// List already created stashes
	ListStashes() ([]string, error)

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

func Detect(path fs.Path) (VersionControlSoftware, error) {
	for _, s := range All {
		if err := s.Detect(path); err != nil && err != errNotYetImpl {
			log.Errorf("error: %s: %s: %s\n", path, s.Name(), err.Error())
		} else {
			return instanciate(s), nil
		}
	}
	return nil, fmt.Errorf("cannot find suitable vcs, tried %v", AllNames)
}

func Open(path fs.Path) (VersionControlSoftware, error) {
	if vc, err := Detect(path); err != nil {
		return nil, err
	} else {
		if err := vc.Open(path); err != nil {
			return nil, err
		}
		return vc, nil
	}
}

func Initialize(n string, p fs.Path, options VersionControlOptions) (VersionControlSoftware, error) {
	v := Get(n)
	return v, v.Initialize(p, options)
}
