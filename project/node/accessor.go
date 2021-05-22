package node

import (
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/project/accessor"
)

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    string
	pkgFile string
	pkg     Package
}

func (a *ProjectAccessor) Name() string {
	return "Node"
}

func (a *ProjectAccessor) Path() string {
	return a.path
}

func (a *ProjectAccessor) Open(p string) error {
	a.path = p
	a.pkg = Package{}
	a.pkgFile = filepath.Join(p, "package.json")
	return a.pkg.ReadFile(a.pkgFile)
}

func (a *ProjectAccessor) Detect(p string) (bool, error) {
	fname := filepath.Join(p, "package.json")
	if _, err := os.Stat(fname); err != nil {
		return false, err
	}
	return true, nil
}

func (a *ProjectAccessor) CurrentVersion() (string, error) {
	return a.pkg.Version()
}

func (a *ProjectAccessor) CurrentName() (string, error) {
	return a.pkg.Name()
}
