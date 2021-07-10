package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    string
	pkgFile string
	pkg     Package
}

func (a *ProjectAccessor) AccessorName() string {
	return "Node"
}

func (a *ProjectAccessor) Path() string {
	return a.path
}

func (a *ProjectAccessor) Open(p string) error {
	a.path = p
	a.pkg = Package{}
	a.pkgFile = filepath.Join(p, a.DescriptionFile())
	return a.pkg.ReadFile(a.pkgFile)
}

func (a *ProjectAccessor) Detect(p string) (bool, error) {
	fname := filepath.Join(p, a.DescriptionFile())
	if _, err := os.Stat(fname); err != nil {
		return false, err
	}
	return true, nil
}

func (a *ProjectAccessor) ReadVersion() (v version.Version, err error) {
	var vs string
	if vs, err = a.pkg.Version(); err != nil {
		return nil, err
	}
	vs = strings.Replace(vs, "-SNAPSHOT", "", 1)
	if v = version.Parse(vs); v == nil {
		return nil, fmt.Errorf("failed to parse version from '%s'", vs)
	}
	return v, nil
}

func (a *ProjectAccessor) Version() (string, error) {
	return a.pkg.Version()
}

func (a *ProjectAccessor) Name() (string, error) {
	return a.pkg.Name()
}

func (a *ProjectAccessor) DescriptionFile() string {
	return "package.json"
}

func (a *ProjectAccessor) Scaffold(ctx *accessor.FinalizationContext) error {
	// exec.RunCommand("")
	return nil
}
