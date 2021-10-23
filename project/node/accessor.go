package node

import (
	"encoding/json"
	"fmt"
	io_fs "io/fs"
	"os"
	"strings"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

type ProjectAccessor struct {
	accessor.ProjectAccessor
	path    fs.Path
	pkgFile string
	pkg     Package
}

func (a *ProjectAccessor) AccessorName() string {
	return "Node"
}

func (a *ProjectAccessor) Path() fs.Path {
	return a.path
}

func (a *ProjectAccessor) Open(p fs.Path) error {
	a.path = p
	a.pkg = Package{}
	a.pkgFile = p.Join(a.DescriptionFile()).Expand()
	return a.pkg.ReadFile(a.pkgFile)
}

func (a *ProjectAccessor) Detect(p fs.Path) (bool, error) {
	fname := p.Join(a.DescriptionFile())
	if _, err := fname.Stat(); err != nil {
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

func (a *ProjectAccessor) WriteVersion(v *version.Version) (err error) {
	var entries []io_fs.DirEntry
	if entries, err = os.ReadDir(a.Path().Expand()); err != nil {
		return
	}
	var content []byte
	var contentStr string
	var path fs.Path
	var currentVersion version.Version
	if currentVersion, err = a.ReadVersion(); err != nil {
		return
	}
	var fi os.FileInfo
	for _, e := range entries {
		for _, vm := range a.VersionManipulators() {
			path = a.Path().Join(e.Name())
			if fi, err = path.Stat(); err != nil {
				return
			}
			if !fi.IsDir() {
				if content, err = os.ReadFile(path.Expand()); err != nil {
					return
				}
				contentStr = string(content)
				if vm.Detect(path.Expand(), contentStr, &currentVersion) {
					if contentStr, err = vm.Update(path.Expand(), contentStr, &currentVersion, v); err != nil {
						return
					}
					if err = os.WriteFile(path.Expand(), []byte(contentStr), 0755); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

// Retrieve the possible list of files that need version updates
func (a *ProjectAccessor) VersionManipulators() map[string]accessor.VersionFileManipulator {
	return map[string]accessor.VersionFileManipulator{
		a.DescriptionFile(): {
			Detect: a.detectPackageVersion,
			Update: a.updatePackageVersion,
		},
	}
}

func (a *ProjectAccessor) detectPackageVersion(file, content string, currentVersion *version.Version) bool {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		log.Errorf("Failed to parse %s content: %s", file, err.Error())
		return false
	}
	if _, ok := data["version"]; !ok {
		log.Errorf("No version key found in %s", file)
		return false
	}
	return true
}

func (a *ProjectAccessor) updatePackageVersion(file, content string, currentVersion, nextVersion *version.Version) (string, error) {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return "", fmt.Errorf("failed to parse %s content: %s", file, err.Error())
	}
	if _, ok := data["version"]; !ok {
		return "", fmt.Errorf("no version key found in %s", file)
	}
	data["version"] = nextVersion.String()
	if output, err := json.MarshalIndent(data, "", "  "); err != nil {
		return "", err
	} else {
		return string(output), nil
	}
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
