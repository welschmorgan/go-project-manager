package maven

import (
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/ui"
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

func (a *ProjectAccessor) Initialize(p string, proj *config.Project) error {
	a.path = p
	a.pom = NewPOMFile()
	a.pomFile = filepath.Join(p, a.DescriptionFile())
	println("------[ Maven POM ]-------")
	var ans string
	var err error
	if ans, err = ui.Ask("\tModelVersion", DefaultPOMModel); err != nil {
		return err
	} else {
		a.pom.Root.SetModelVersion(ParseModelVersion(ans))
	}
	if ans, err = ui.Ask("\tArtifactId", proj.Name); err != nil {
		return err
	} else {
		a.pom.Root.ArtifactId = ans
	}
	if ans, err = ui.Ask("\tGroupId", "com."); err != nil {
		return err
	} else {
		a.pom.Root.GroupId = ans
	}
	if ans, err = ui.Ask("\tVersion", DefaultPOMVersion); err != nil {
		return err
	} else {
		a.pom.Root.Version = ans
	}
	return a.pom.WriteFile(a.pomFile)
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
