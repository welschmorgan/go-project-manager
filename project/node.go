package project

import (
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/config"
)

type NodeProjectAccessor struct {
	ProjectAccessor
	infos   *config.Project
	pkgFile string
	pkg     NodePackage
}

func NewNodeProjectAccessor(infos *config.Project) *NodeProjectAccessor {
	return &NodeProjectAccessor{
		infos:   infos,
		pkg:     NodePackage{},
		pkgFile: filepath.Join(infos.Path, "package.json"),
	}
}

func (a *NodeProjectAccessor) Infos() *config.Project {
	return a.infos
}

func (a *NodeProjectAccessor) Detect() (bool, error) {
	if err := a.pkg.ReadFile(a.pkgFile); err != nil {
		return false, err
	}
	return true, nil
}

func (a *NodeProjectAccessor) CurrentVersion() (string, error) {
	return a.pkg.Version()
}
