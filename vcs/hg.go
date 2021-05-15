package vcs

import "github.com/welschmorgan/go-project-manager/models"

type Hg struct {
	VersionControlSoftware
	path string
}

func (h *Hg) Name() string                           { return "Hg" }
func (h *Hg) Path() string                           { return h.path }
func (h *Hg) Detect(path string) (bool, error)       { return false, errNotYetImpl }
func (h *Hg) Open(p string) error                    { return errNotYetImpl }
func (h *Hg) Clone(url, path string) error           { return errNotYetImpl }
func (h *Hg) Checkout(branch string) error           { return errNotYetImpl }
func (h *Hg) Pull() error                            { return errNotYetImpl }
func (h *Hg) Push() error                            { return errNotYetImpl }
func (h *Hg) Tag(name, commit, message string) error { return errNotYetImpl }
func (h *Hg) Merge(source, dest string) error        { return errNotYetImpl }
func (s *Hg) Authors() ([]*models.Person, error)     { return nil, errNotYetImpl }
