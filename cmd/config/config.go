package config

import (
	"os"

	"github.com/welschmorgan/go-project-manager/models"
)

var (
	// Used for flags.
	DefaultWorkspaceFilename string = ".grlm-workspace.yaml"
	DefaultDevelopmentBranch string = "develop"
	DefaultProductionBranch  string = "master"
	DefaultVerbose           bool   = false

	config *AppConfig
)

func init() {
	var err error
	if config, err = NewAppConfig(); err != nil {
		panic(err.Error())
	}
}

type BranchNamesConfig map[string]string

type ReleaseConfig struct {
	BrancheNames BranchNamesConfig
}

type AppConfig struct {
	WorkspacesRoot    string
	Release           ReleaseConfig
	Verbose           bool
	CfgFile           string
	WorkingDirectory  string
	WorkspaceFilename string
	WorkspacePath     string
	Workspace         models.Workspace
}

func NewAppConfig() (*AppConfig, error) {
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		return &AppConfig{
			WorkspacesRoot: "$HOME/development",
			Release: ReleaseConfig{
				BrancheNames: BranchNamesConfig{
					"development": DefaultDevelopmentBranch,
					"production":  DefaultProductionBranch,
				},
			},
			Verbose:           DefaultVerbose,
			CfgFile:           "",
			WorkingDirectory:  cwd,
			WorkspaceFilename: DefaultWorkspaceFilename,
		}, nil
	}
}

func Get() *AppConfig {
	return config
}
