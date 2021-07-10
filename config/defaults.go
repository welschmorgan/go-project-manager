package config

var (
	// Used for flags.
	DefaultWorkspaceFilename  string       = ".grlm/workspace.yaml"
	DefaultWorkspacesRoot     string       = "${home}/development"
	DefaultDevelopmentBranch  string       = "develop"
	DefaultProductionBranch   string       = "master"
	DefaultReleaseBranch      string       = "release/${version}"
	DefaultVerbose            VerboseLevel = NoVerbose
	DefaultDryRun             bool         = false
	DefaultInteractive        bool         = false
	DefaultAcquireVersionFrom string       = "package"
	DefaultLogFolder          string       = "${workspace}/.grlm/logs"

	instance *Config
)
