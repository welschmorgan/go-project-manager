package config

import (
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/version"
)

type Config struct {
	Workspace
	Indent            int
	WorkspacesRoot    string
	Verbose           VerboseLevel
	CfgFile           string
	WorkingDirectory  string
	WorkspaceFilename string
	WorkspacePath     string
	DryRun            bool
	Interactive       bool
	LogFolder         string
	ReleaseType       version.VersionPart
}

func NewConfig() (*Config, error) {
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		return &Config{
			Workspace:         *NewWorkspace(),
			WorkspacesRoot:    DefaultWorkspacesRoot,
			Verbose:           DefaultVerbose,
			Indent:            0,
			CfgFile:           "",
			WorkingDirectory:  cwd,
			WorkspaceFilename: DefaultWorkspaceFilename,
			WorkspacePath:     filepath.Join(cwd, DefaultWorkspaceFilename),
			DryRun:            DefaultDryRun,
			Interactive:       DefaultInteractive,
			LogFolder:         DefaultLogFolder,
			ReleaseType:       version.Minor,
		}, nil
	}
}

func Get() *Config {
	return instance
}

func init() {
	var err error
	if instance, err = NewConfig(); err != nil {
		panic(err.Error())
	}
}
