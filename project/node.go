package project

import (
	"os"
	"path/filepath"
)

type NodeProjectAccessor struct {
	ProjectAccessor
	path    string
	pkgFile string
	pkg     NodePackage
}

func (a *NodeProjectAccessor) Name() string {
	return "Node"
}

func (a *NodeProjectAccessor) Path() string {
	return a.path
}

func (a *NodeProjectAccessor) Open(p string) error {
	a.path = p
	a.pkg = NodePackage{}
	a.pkgFile = filepath.Join(p, "package.json")
	return a.pkg.ReadFile(a.pkgFile)
}

func (a *NodeProjectAccessor) Detect(p string) (bool, error) {
	fname := filepath.Join(p, "package.json")
	if _, err := os.Stat(fname); err != nil {
		return false, err
	}
	return true, nil
}

func (a *NodeProjectAccessor) CurrentVersion() (string, error) {
	return a.pkg.Version()
}
