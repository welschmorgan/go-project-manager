package config

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/version"
	"gopkg.in/yaml.v2"
)

type BranchNamesConfig map[string]string

type Versionning struct {
	PreReleasePrefix string
}
type Workspace struct {
	Name               string            `yaml:"name"`
	path               string            `yaml:"path"`
	Projects           []*Project        `yaml:"projects"`
	Author             *Person           `yaml:"author"`
	Manager            *Person           `yaml:"manager"`
	Developpers        []*Person         `yaml:"developpers"`
	BranchNames        BranchNamesConfig `yaml:"branch_names"`
	Versionning        Versionning       `yaml:"versionning"`
	AcquireVersionFrom string            `yaml:"acquire_version_from"`
}

func NewWorkspace() *Workspace {
	path := DefaultWorkspacesRoot
	name := fmt.Sprintf("workspace %5d", rand.Int())
	return &Workspace{
		Name:        name,
		path:        path,
		Projects:    []*Project{},
		Author:      nil,
		Manager:     nil,
		Developpers: []*Person{},
		BranchNames: BranchNamesConfig{
			"development": DefaultDevelopmentBranch,
			"production":  DefaultProductionBranch,
			"release":     DefaultReleaseBranch,
		},
		Versionning: Versionning{
			PreReleasePrefix: version.PreReleasePrefix,
		},
		AcquireVersionFrom: DefaultAcquireVersionFrom,
	}
}

func NewWorkspaceWithValues(name, path string, projects []*Project, sourceControl string, author *Person, manager *Person, developpers []*Person, branchNames BranchNamesConfig, acquireVersionFrom string) *Workspace {
	return &Workspace{
		Name:               name,
		path:               path,
		Projects:           projects,
		Author:             author,
		Manager:            manager,
		Developpers:        developpers,
		BranchNames:        branchNames,
		AcquireVersionFrom: acquireVersionFrom,
	}
}

func (w *Workspace) ReadFile(path string) error {
	if content, err := os.ReadFile(path); err != nil {
		return err
	} else if err = w.Read(content); err != nil {
		return err
	}
	return nil
}

func (w *Workspace) Read(b []byte) error {
	return yaml.Unmarshal(b, w)
}

func (w *Workspace) WriteFile(path string) error {
	if yaml, err := w.Write(); err != nil {
		return err
	} else {
		if err := os.WriteFile(path, yaml, 0755); err != nil {
			return err
		}
		fmt.Printf("Written '%s':\n%s\n", path, yaml)
	}
	return nil
}

func (w *Workspace) Write() ([]byte, error) {
	return yaml.Marshal(w)
}

func (w *Workspace) Path() string {
	return w.path
}

func (w *Workspace) SetPath(p string) {
	w.path = p
	fs.SanitizePath(w.path, nil)
}

func (w *Workspace) LogFolder() string {
	return fs.SanitizePath(Get().LogFolder, map[string]string{
		"WORKSPACE": w.Path(),
	})
}
