package vcs

import "github.com/welschmorgan/go-project-manager/models"

type Svn struct {
	VersionControlSoftware
	path string
}

func (s *Svn) Name() string                           { return "Svn" }
func (s *Svn) Path() string                           { return s.path }
func (s *Svn) Detect(path string) (bool, error)       { return false, errNotYetImpl }
func (s *Svn) Open(p string) error                    { return errNotYetImpl }
func (s *Svn) Clone(url, path string) error           { return errNotYetImpl }
func (s *Svn) Checkout(branch string) error           { return errNotYetImpl }
func (s *Svn) Pull() error                            { return errNotYetImpl }
func (s *Svn) Push() error                            { return errNotYetImpl }
func (s *Svn) Tag(name, commit, message string) error { return errNotYetImpl }
func (s *Svn) Merge(source, dest string) error        { return errNotYetImpl }
func (s *Svn) Authors() ([]*models.Person, error)     { return nil, errNotYetImpl }
