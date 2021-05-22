package vcs

import "github.com/welschmorgan/go-project-manager/models"

type Hg struct {
	VersionControlSoftware
	path string
	url  string
}

func (h *Hg) Name() string                     { return "Hg" }
func (h *Hg) Path() string                     { return h.path }
func (h *Hg) Url() string                      { return h.url }
func (h *Hg) Detect(path string) (bool, error) { return false, errNotYetImpl }
func (h *Hg) Open(p string) error              { return errNotYetImpl }
func (h *Hg) Clone(url, path string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Status(options StatusOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (h *Hg) Stash(options StashOptions) ([]string, error) {
	return nil, errNotYetImpl
}
func (h *Hg) Checkout(branch string, options VersionControlOptions) error { return errNotYetImpl }
func (h *Hg) Pull(options VersionControlOptions) error                    { return errNotYetImpl }
func (h *Hg) Push(options VersionControlOptions) error                    { return errNotYetImpl }
func (h *Hg) Tag(name, commit, message string, options VersionControlOptions) error {
	return errNotYetImpl
}
func (h *Hg) Merge(source, dest string, options VersionControlOptions) error { return errNotYetImpl }
func (s *Hg) Authors(options VersionControlOptions) ([]*models.Person, error) {
	return nil, errNotYetImpl
}
