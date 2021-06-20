package maven

import (
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/project/accessor"
)

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    string
	pomFile string
	pom     *POMFile
}

func (a *ProjectAccessor) AccessorName() string {
	return "Maven"
}

func (a *ProjectAccessor) Path() string {
	return a.path
}

func (a *ProjectAccessor) DescriptionFile() string {
	return "pom.xml"
}

func (a *ProjectAccessor) Scaffold(ctx *accessor.FinalizationContext) error {
	a.path = ctx.Project.Path
	a.pom = NewPOMFile()
	a.pomFile = filepath.Join(ctx.Project.Path, a.DescriptionFile())
	return a.Scaffolder().Scaffold(ctx)
}

func (a *ProjectAccessor) Scaffolder() accessor.Scaffolder {
	return &MavenScaffolder{}
}

func (a *ProjectAccessor) Open(p string) error {
	a.path = p
	a.pom = NewPOMFile()
	a.pomFile = filepath.Join(p, a.DescriptionFile())
	return a.pom.ReadFile(a.pomFile)
}

func (a *ProjectAccessor) Detect(p string) (bool, error) {
	fname := filepath.Join(p, a.DescriptionFile())
	if _, err := os.Stat(fname); err != nil {
		return false, err
	}
	return true, nil
}

func (a *ProjectAccessor) Version() (string, error) {
	return a.pom.Root.Version, nil
}

func (a *ProjectAccessor) Name() (string, error) {
	return a.pom.Root.ArtifactId, nil
}
