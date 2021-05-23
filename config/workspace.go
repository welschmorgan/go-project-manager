package config

import (
	"fmt"
	"math/rand"
	"os"

	"gopkg.in/yaml.v2"
)

type BranchNamesConfig map[string]string

type Workspace struct {
	Name               string            `yaml:"name"`
	Path               string            `yaml:"path"`
	Projects           []*Project        `yaml:"projects"`
	Author             *Person           `yaml:"author"`
	Manager            *Person           `yaml:"manager"`
	Developpers        []*Person         `yaml:"developpers"`
	BranchNames        BranchNamesConfig `yaml:"branch_names"`
	AcquireVersionFrom string            `yaml:"acquire_version_from"`
}

func NewWorkspace() *Workspace {
	return &Workspace{
		Name:        fmt.Sprintf("workspace %d", rand.Int()),
		Path:        DefaultWorkspacesRoot,
		Projects:    []*Project{},
		Author:      nil,
		Manager:     nil,
		Developpers: []*Person{},
		BranchNames: BranchNamesConfig{
			"development": DefaultDevelopmentBranch,
			"production":  DefaultProductionBranch,
			"release":     DefaultReleaseBranch,
		},
		AcquireVersionFrom: DefaultAcquireVersionFrom,
	}
}

func NewWorkspaceWithValues(name, path string, projects []*Project, sourceControl string, author *Person, manager *Person, developpers []*Person, branchNames BranchNamesConfig, acquireVersionFrom string) *Workspace {
	return &Workspace{
		Name:               name,
		Path:               path,
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
