package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/version"
)

type Config struct {
	Workspace
	Indent            int
	WorkspacesRoot    fs.Path
	Verbose           VerboseLevel
	CfgFile           string
	WorkingDirectory  string
	WorkspaceFilename string
	WorkspacePath     fs.Path
	DryRun            bool
	Interactive       bool
	LogFolder         fs.Path
	ReleaseType       version.VersionPart
}

func NewConfig() (*Config, error) {
	if cwd, err := os.Getwd(); err != nil {
		return nil, err
	} else {
		return &Config{
			Workspace:         *NewWorkspace(),
			WorkspacesRoot:    fs.Path(DefaultWorkspacesRoot),
			Verbose:           DefaultVerbose,
			Indent:            0,
			CfgFile:           "",
			WorkingDirectory:  cwd,
			WorkspaceFilename: DefaultWorkspaceFilename,
			WorkspacePath:     fs.Path(filepath.Join(cwd, DefaultWorkspaceFilename)),
			DryRun:            DefaultDryRun,
			Interactive:       DefaultInteractive,
			LogFolder:         fs.Path(DefaultLogFolder),
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

func (c *Config) Json() string {
	if jsonCfg, err := json.MarshalIndent(*c, "", "  "); err != nil {
		panic(err.Error())
	} else {
		return string(jsonCfg)
	}
}
