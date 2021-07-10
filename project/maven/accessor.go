package maven

import (
	"fmt"
	io_fs "io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    fs.Path
	pomFile string
	pom     *POMProject
}

func (a *ProjectAccessor) AccessorName() string {
	return "Maven"
}

func (a *ProjectAccessor) Path() fs.Path {
	return a.path
}

func (a *ProjectAccessor) DescriptionFile() string {
	return "pom.xml"
}

func (a *ProjectAccessor) Scaffold(ctx *accessor.FinalizationContext) error {
	if err := a.Open(ctx.Project.Path); err != nil {
		return err
	}
	return a.Scaffolder().Scaffold(ctx)
}

func (a *ProjectAccessor) Scaffolder() accessor.Scaffolder {
	return &MavenScaffolder{}
}

func (a *ProjectAccessor) Open(p fs.Path) error {
	a.path = p
	a.pom = NewPOMProject()
	a.pomFile = p.Join(a.DescriptionFile()).Expand()
	return a.pom.ReadFile(a.pomFile)
}

func (a *ProjectAccessor) Detect(p fs.Path) (bool, error) {
	fname := p.Join(a.DescriptionFile())
	if _, err := fname.Stat(); err != nil {
		return false, err
	}
	return true, nil
}

func (a *ProjectAccessor) Name() (string, error) {
	return a.pom.ArtifactId, nil
}

func (a *ProjectAccessor) ReadVersion() (v version.Version, err error) {
	vs := a.pom.Version
	vs = strings.Replace(vs, "-SNAPSHOT", "", 1)
	if v = version.Parse(vs); v == nil {
		return nil, fmt.Errorf("failed to parse version from '%s'", vs)
	}
	return v, nil
}

func (a *ProjectAccessor) WriteVersion(v *version.Version) (err error) {
	var entries []io_fs.DirEntry
	var content []byte
	var contentStr string
	var path fs.Path
	var currentVersion version.Version
	if currentVersion, err = a.ReadVersion(); err != nil {
		return
	}
	var fi os.FileInfo
	entries, err = a.Path().ReadDir()
	for _, e := range entries {
		for _, vm := range a.VersionManipulators() {
			path = a.Path().Join(e.Name())
			if fi, err = path.Stat(); err != nil {
				return
			}
			if !fi.IsDir() {
				if content, err = path.ReadFile(); err != nil {
					return
				}
				contentStr = string(content)
				if vm.Detect(path.Expand(), contentStr, &currentVersion) {
					if contentStr, err = vm.Update(path.Expand(), contentStr, &currentVersion, v); err != nil {
						return
					}
					if err = path.WriteFile([]byte(contentStr)); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func (a *ProjectAccessor) detectPOMVersion(file, content string, currentVersion *version.Version) bool {
	needle := fmt.Sprintf(`<version>\s*%s\s*(-SNAPSHOT|)</version>`, currentVersion)
	re := regexp.MustCompile(needle)
	return re.MatchString(content)
}

func (a *ProjectAccessor) updatePOMVersion(file, content string, currentVersion, nextVersion *version.Version) (string, error) {
	needle := fmt.Sprintf(`<version>\s*%s(-SNAPSHOT|)\s*</version>`, currentVersion)
	re := regexp.MustCompile(needle)
	return re.ReplaceAllString(content, fmt.Sprintf("<version>%s$1</version>", nextVersion)), nil
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
