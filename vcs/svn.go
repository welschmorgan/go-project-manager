package vcs

type Svn struct{}

func (s *Svn) Name() string                           { return "Svn" }
func (s *Svn) Clone(url, path string) error           { return errNotYetImpl }
func (s *Svn) Checkout(branch string) error           { return errNotYetImpl }
func (s *Svn) Pull() error                            { return errNotYetImpl }
func (s *Svn) Push() error                            { return errNotYetImpl }
func (s *Svn) Tag(name, commit, message string) error { return errNotYetImpl }
func (s *Svn) Merge(source, dest string) error        { return errNotYetImpl }
