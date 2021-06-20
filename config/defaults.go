package config

var (
	// Used for flags.
	DefaultWorkspaceFilename  string       = ".grlm-workspace.yaml"
	DefaultWorkspacesRoot     string       = "$HOME/development"
	DefaultDevelopmentBranch  string       = "develop"
	DefaultProductionBranch   string       = "master"
	DefaultReleaseBranch      string       = "release/$VERSION"
	DefaultVerbose            VerboseLevel = NoVerbose
	DefaultDryRun             bool         = false
	DefaultInteractive        bool         = false
	DefaultAcquireVersionFrom string       = "package"
	DefaultLogFolder          string       = "$WORKSPACE/logs"

	instance *Config
)
