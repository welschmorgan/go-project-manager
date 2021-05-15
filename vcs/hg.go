package vcs

type Hg struct{}

func (s *Hg) Name() string                           { return "Hg" }
func (h *Hg) Clone(url, path string) error           { return errNotYetImpl }
func (h *Hg) Checkout(branch string) error           { return errNotYetImpl }
func (h *Hg) Pull() error                            { return errNotYetImpl }
func (h *Hg) Push() error                            { return errNotYetImpl }
func (h *Hg) Tag(name, commit, message string) error { return errNotYetImpl }
func (h *Hg) Merge(source, dest string) error        { return errNotYetImpl }
