package maven

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
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

func (a *ProjectAccessor) Name() (string, error) {
	return a.pom.Root.ArtifactId, nil
}

func (a *ProjectAccessor) ReadVersion() (v version.Version, err error) {
	vs := a.pom.Root.Version
	vs = strings.Replace(vs, "-SNAPSHOT", "", 1)
	if v = version.Parse(vs); v == nil {
		return nil, fmt.Errorf("failed to parse version from '%s'", vs)
	}
	return v, nil
}

func (a *ProjectAccessor) WriteVersion(v *version.Version) (err error) {
	var entries []fs.DirEntry
	if entries, err = os.ReadDir(a.Path()); err != nil {
		return
	}
	var content []byte
	var contentStr string
	var path string
	var currentVersion version.Version
	if currentVersion, err = a.ReadVersion(); err != nil {
		return
	}
	var fi os.FileInfo
	for _, e := range entries {
		for _, vm := range a.VersionManipulators() {
			path = filepath.Join(a.Path(), e.Name())
			if fi, err = os.Stat(path); err != nil {
				return
			}
			if !fi.IsDir() {
				if content, err = os.ReadFile(path); err != nil {
					return
				}
				contentStr = string(content)
				if vm.Detect(path, contentStr, &currentVersion) {
					if contentStr, err = vm.Update(path, contentStr, &currentVersion, v); err != nil {
						return
					}
					if err = os.WriteFile(path, []byte(contentStr), 0755); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func (a *ProjectAccessor) detectPOMVersion(file, content string, currentVersion *version.Version) bool {
	return strings.Contains(content, fmt.Sprintf("<version>%s</version>", currentVersion))
}

func (a *ProjectAccessor) updatePOMVersion(file, content string, currentVersion, nextVersion *version.Version) (string, error) {
	return strings.ReplaceAll(content, fmt.Sprintf("<version>%s</version>", currentVersion), fmt.Sprintf("<version>%s</version>", nextVersion)), nil
}

func (a *ProjectAccessor) detectSonarVersion(file, content string, currentVersion *version.Version) bool {
	return strings.Contains(content, fmt.Sprintf("version: %s", currentVersion))
}

func (a *ProjectAccessor) updateSonarVersion(file, content string, currentVersion, nextVersion *version.Version) (string, error) {
	return strings.ReplaceAll(content, fmt.Sprintf("version: %s", currentVersion), fmt.Sprintf("version: %s", nextVersion)), nil
}

// func (a *ProjectAccessor) detectJenkinsfile(file, content string, currentVersion *version.Version) bool {
// 	return strings.Contains(content, fmt.Sprintf("<version>%s</version>", currentVersion))
// }

// func (a *ProjectAccessor) updateJendetectJenkinsfile(file, content string, currentVersion *version.Version) bool {
// 	strings.ReplaceAll(content, fmt.Sprintf("<version>%s</version>", currentVersion), fmt.Sprintf("<version>%s</version>", nextVersion))
// }

// Retrieve the possible list of files that need version updates
func (a *ProjectAccessor) VersionManipulators() map[string]accessor.VersionFileManipulator {
	return map[string]accessor.VersionFileManipulator{
		a.DescriptionFile(): {
			Detect: a.detectPOMVersion,
			Update: a.updatePOMVersion,
		},

		"sonar-project.properties": {
			Detect: a.detectSonarVersion,
			Update: a.updateSonarVersion,
		},
		// "Jenkinsfile": {
		// 	Detect: a.detectJenkinsfile(),
		// 	Update: a.updateJenkinsfile(),
		// }
	}
}
