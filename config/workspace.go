package config

import (
	"fmt"
	"math/rand"
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
