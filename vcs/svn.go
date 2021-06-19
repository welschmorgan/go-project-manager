package vcs

import (
	"errors"

	"github.com/welschmorgan/go-release-manager/config"
)

type Svn struct {
	VersionControlSoftware
	path string
	url  string
}

func (s *Svn) Name() string                     { return "Svn" }
func (s *Svn) Path() string                     { return s.path }
func (s *Svn) Url() string                      { return s.url }
func (s *Svn) Detect(path string) (bool, error) { return false, errNotYetImpl }
func (s *Svn) Open(p string) error              { return errNotYetImpl }
func (s *Svn) Clone(url, path string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (s *Svn) Status(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (s *Svn) Stash(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (s *Svn) Checkout(branch string, options VersionControlOptions) error { return errNotYetImpl }
func (s *Svn) Pull(options VersionControlOptions) error                    { return errNotYetImpl }
func (s *Svn) Push(options VersionControlOptions) error                    { return errNotYetImpl }
func (s *Svn) Tag(name string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (s *Svn) Merge(source, dest string, options VersionControlOptions) error { return errNotYetImpl }
func (s *Svn) ListAuthors(options VersionControlOptions) ([]*config.Person, error) {
	return nil, errNotYetImpl
}
func (s *Svn) CurrentBranch() (string, error) {
	return "", errNotYetImpl
}
func (s *Svn) DeleteBranch(name string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (s *Svn) Reset(options VersionControlOptions) error {
	return errNotYetImpl
}
func (s *Svn) ListTags(options VersionControlOptions) ([]string, error) {
	return nil, errNotYetImpl
}

func (s *Svn) Initialize(path string, options VersionControlOptions) error {
	return errors.New("Not yet implemented")
}

func (s *Svn) Commit(options VersionControlOptions) error {
	return errors.New("Not yet implemented")
}
