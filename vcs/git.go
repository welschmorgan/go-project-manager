package vcs

type Git struct{ VersionControlSoftware }

func (s *Git) Name() string                           { return "Git" }
func (g *Git) Clone(url, path string) error           { return errNotYetImpl }
func (g *Git) Checkout(branch string) error           { return errNotYetImpl }
func (g *Git) Pull() error                            { return errNotYetImpl }
func (g *Git) Push() error                            { return errNotYetImpl }
func (g *Git) Tag(name, commit, message string) error { return errNotYetImpl }
func (g *Git) Merge(source, dest string) error        { return errNotYetImpl }
