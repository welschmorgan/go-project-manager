package vcs

import "errors"

var (
	errNotYetImpl = errors.New("not yet implemented")
)

type VersionControlSoftware interface {
	Name() string
	Clone(url, path string) error
	Checkout(branch string) error
	Pull() error
	Push() error
	Tag(name, commit, message string) error
	Merge(source, dest string) error
}

var VersionControlSoftwares = []VersionControlSoftware{
	&Git{},
	&Svn{},
	&Hg{},
}
