package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Workspace
	WorkspacesRoot    string
	Verbose           bool
	CfgFile           string
	WorkingDirectory  string
	WorkspaceFilename string
	WorkspacePath     string
	DryRun            bool
}

func NewConfig() (*Config, error) {
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		return &Config{
			Workspace:         *NewWorkspace(),
			WorkspacesRoot:    DefaultWorkspacesRoot,
			Verbose:           DefaultVerbose,
			CfgFile:           "",
			WorkingDirectory:  cwd,
			WorkspaceFilename: DefaultWorkspaceFilename,
			WorkspacePath:     filepath.Join(cwd, DefaultWorkspaceFilename),
			DryRun:            DefaultDryRun,
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
