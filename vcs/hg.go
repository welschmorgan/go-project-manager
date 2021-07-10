package vcs

import (
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
)

type Hg struct {
	VersionControlSoftware
	path fs.Path
	url  string
}

func (h *Hg) Name() string              { return "Hg" }
func (h *Hg) Path() fs.Path             { return h.path }
func (h *Hg) Url() string               { return h.url }
func (h *Hg) Detect(path fs.Path) error { return errNotYetImpl }
func (h *Hg) Open(p fs.Path) error      { return errNotYetImpl }
func (h *Hg) Clone(url string, path fs.Path, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Status(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (h *Hg) Stash(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (h *Hg) Checkout(branch string, options VersionControlOptions) error { return errNotYetImpl }
func (h *Hg) Pull(options VersionControlOptions) error                    { return errNotYetImpl }
func (h *Hg) Push(options VersionControlOptions) error                    { return errNotYetImpl }
func (h *Hg) Tag(name string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) CurrentBranch() (string, error) {
	return "", errNotYetImpl
}
func (h *Hg) Merge(source, dest string, options VersionControlOptions) error { return errNotYetImpl }
func (h *Hg) ListAuthors(options VersionControlOptions) ([]*config.Person, error) {
	return nil, errNotYetImpl
}

func (h *Hg) DeleteBranch(name string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Reset(options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) ListTags(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}

func (h *Hg) Initialize(path fs.Path, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Commit(options VersionControlOptions) error {
	return errNotYetImpl
}

func (h *Hg) Stage(options VersionControlOptions) error {
	return errNotYetImpl
}

// Retrieve commits without parents
func (h *Hg) RootCommits() ([]string, error) {
	return nil, errNotYetImpl
}

// Get the hash of the current commit
func (h *Hg) CurrentCommit(options VersionControlOptions) (hash, subject string, err error) {
	err = errNotYetImpl
	return
}

// Get the hash of the current commit
func (h *Hg) ExtractLog(options VersionControlOptions) (lines []string, err error) {
	err = errNotYetImpl
	return
}

// List already created stashes
func (h *Hg) ListStashes() (lines []string, err error) {
	err = errNotYetImpl
	return
}

// Fetch remote index
func (h *Hg) FetchIndex(options VersionControlOptions) error {
	return errNotYetImpl
}
