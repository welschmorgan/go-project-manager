package vcs

import (
	"fmt"
	"reflect"
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

type CommitOptions struct {
	Signed     bool
	Message    string
	AllowEmpty bool
	StageFiles bool
}

type StageOptions struct {
	All              bool
	AllAlreadyStaged bool
	Files            []string
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

type CurrentCommitOptions struct {
	ShortHash bool
}

type Sorter struct {
	Column     string
	Descending bool
}

type ExtractLogOptions struct {
	Limit  int
	Format string
	Branch string
}
