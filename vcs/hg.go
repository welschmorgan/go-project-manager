package vcs

import (
	"github.com/welschmorgan/go-release-manager/config"
)

type Hg struct {
	VersionControlSoftware
	path string
	url  string
}

func (h *Hg) Name() string             { return "Hg" }
func (h *Hg) Path() string             { return h.path }
func (h *Hg) Url() string              { return h.url }
func (h *Hg) Detect(path string) error { return errNotYetImpl }
func (h *Hg) Open(p string) error      { return errNotYetImpl }
func (h *Hg) Clone(url, path string, options VersionControlOptions) error {
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

func (h *Hg) Initialize(path string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Commit(options VersionControlOptions) error {
	return errNotYetImpl
}

func (h *Hg) Stage(options VersionControlOptions) error {
	return errNotYetImpl
}

// Retrieve commits without parents
func (h *Hg) GetRootCommits() ([]string, error) {
	return nil, errNotYetImpl
}
