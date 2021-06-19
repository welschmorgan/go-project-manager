package maven

import (
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
)

const DefaultPOMModel = "4.0.0"

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    string
	pomFile string
	pom     POMFile
}

func (a *ProjectAccessor) Name() string {
	return "Maven"
}

func (a *ProjectAccessor) Path() string {
	return a.path
}

func (a *ProjectAccessor) Initialize(p string, proj *config.Project) error {
	a.path = p
	a.pom = NewPOMFile(DefaultPOMModel)
	a.pomFile = filepath.Join(p, "pom.xml")
	a.pom.Root.ArtifactId = proj.Name
	a.pom.Root.Version = "0.1.0-SNAPSHOT"
	return a.pom.WriteFile(a.pomFile)
}

func (a *ProjectAccessor) Open(p string) error {
	a.path = p
	a.pom = NewPOMFile(DefaultPOMModel)
	a.pomFile = filepath.Join(p, "pom.xml")
	return a.pom.ReadFile(a.pomFile)
}

func (a *ProjectAccessor) Detect(p string) (bool, error) {
	fname := filepath.Join(p, "package.json")
	if _, err := os.Stat(fname); err != nil {
		return false, err
	}
	return true, nil
}

func (a *ProjectAccessor) CurrentVersion() (string, error) {
	return a.pom.Root.Version, nil
}

func (a *ProjectAccessor) CurrentName() (string, error) {
	return a.pom.Root.ArtifactId, nil
}
