package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag"
	initCommand "github.com/welschmorgan/go-release-manager/cmd/init"
	releaseCommand "github.com/welschmorgan/go-release-manager/cmd/release"
	undoCommand "github.com/welschmorgan/go-release-manager/cmd/undo"
	versionCommand "github.com/welschmorgan/go-release-manager/cmd/version"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/fs"
	"github.com/welschmorgan/go-release-manager/log"
)

var workspacesRoot string
var logFolder string

var Command = &cobra.Command{
	Use:          "grlm [commands]",
	Short:        "Release multiple projects in a single go",
	Long:         `GRLM allows releasing multiple projects declared in a workspace`,
	SilenceUsage: true,
}

// Execute executes the root command.
func Execute() error {
	return Command.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	if cwd, err := os.Getwd(); err != nil {
		panic(err.Error())
	} else {
		config.Get().WorkingDirectory = cwd
	}
	// config file
	Command.PersistentFlags().StringVarP(&config.Get().CfgFile, "config", "c", config.Get().CfgFile, "config file (default is $HOME/.grlm.yaml)")
	viper.BindPFlag("config", Command.PersistentFlags().Lookup("config"))

	// verbose

	var VerboseLevels = map[config.VerboseLevel][]string{}
	for _, v := range config.VerboseLevels {
		VerboseLevels[v] = v.TextualRepresentations()
	}
	Command.PersistentFlags().VarP(
		enumflag.New(&config.Get().Verbose, "verbose", VerboseLevels, enumflag.EnumCaseInsensitive),
		"verbose",
		"v",
		"show additionnal log messages; can be 'none', 'low', 'normal', 'high', 'max'")
	viper.BindPFlag("verbose", Command.PersistentFlags().Lookup("verbose"))

	// dry run
	Command.PersistentFlags().BoolVarP(&config.Get().DryRun, "dry_run", "n", config.Get().DryRun, "simulate commande execution, do not execute them")
	viper.BindPFlag("dry_run", Command.PersistentFlags().Lookup("dry_run"))

	// change working dir
	Command.PersistentFlags().StringVarP(&config.Get().WorkingDirectory, "working_directory", "C", config.Get().WorkingDirectory, "change working directory")
	viper.BindPFlag("working_directory", Command.PersistentFlags().Lookup("working_directory"))

	workspacesRoot = config.Get().WorkspacesRoot.Expand()
	logFolder = config.Get().LogFolder.Expand()

	// define workspaces root
	Command.PersistentFlags().StringVar(&workspacesRoot, "workspaces_root", workspacesRoot, "The root folder where to find workspaces")
	viper.BindPFlag("workspaces_root", Command.PersistentFlags().Lookup("workspaces_root"))

	// define log output dir
	Command.PersistentFlags().StringVar(&logFolder, "log_folder", logFolder, "change where to write logs")
	viper.BindPFlag("log_folder", Command.PersistentFlags().Lookup("log_folder"))

	// Command.ActionAddCommand(addCmd)
	Command.AddCommand(initCommand.Command)
	Command.AddCommand(releaseCommand.Command)
	Command.AddCommand(versionCommand.Command)
	Command.AddCommand(undoCommand.Command)
}

func updateWorkspacePaths() (err error) {
	cfg := config.Get()
	cfg.WorkspacesRoot = fs.Path(workspacesRoot)
	cfg.LogFolder = fs.Path(logFolder)

	if cfg.WorkingDirectory, err = filepath.Abs(cfg.WorkingDirectory); err != nil {
		return err
	}

	if len(cfg.WorkingDirectory) != 0 {
		cfg.Workspace.Path = fs.Path(cfg.WorkingDirectory)
		cfg.Workspace.Name = filepath.Base(cfg.Workspace.Path.Raw())
	}
	if _, err := os.Stat(cfg.WorkingDirectory); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(cfg.WorkingDirectory, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "\033[1;33merror\033[0m: %s\n", err.Error())
		}
	}
	if path, err := filepath.Abs(filepath.Join(cfg.WorkingDirectory, cfg.WorkspaceFilename)); err != nil {
		return err
	} else {
		cfg.WorkspacePath = fs.Path(path)
	}
	fs.PutPathEnv("workspace", cfg.Workspace.Path.Expand)

	if cfg.WorkspacePath.Exists() {
		log.Infof("[\033[1;34m+\033[0m] Using local config file: %s\n", cfg.WorkspacePath)
		if err = cfg.Workspace.ReadFile(cfg.WorkspacePath.Expand()); err != nil {
			return err
		}
		for _, proj := range cfg.Workspace.Projects {
			proj.Path = proj.Path.TrimSpace()
			proj.Path = proj.Path.ReplaceAll("./", "${workspace}/")
		}
	}
	return nil
}

func initConfig() {
	cfg := config.Get()
	if cfg.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg.CfgFile)
	} else {
		viper.SetConfigName("grlm")        // name of config file (without extension)
		viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/grlm/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.grlm") // call multiple times to add many search paths
		viper.AddConfigPath(".")           // optionally look for config in the working directory
	}

	viper.AutomaticEnv()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	var err error
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("configuration error: %s", err))
		}
	}

	log.Debugf("[\033[1;34m+\033[0m] Using config file: %s\n", viper.ConfigFileUsed())

	if err = updateWorkspacePaths(); err != nil {
		panic(err.Error())
	}

	if err = log.Setup(); err != nil {
		panic(err.Error())
	}

	if err = os.Chdir(cfg.WorkingDirectory); err != nil {
		panic(err.Error())
	}

	cfg.Workspace.Versionning.PreReleasePrefix = cfg.Versionning.PreReleasePrefix

	log.Debugf("[\033[1;34m+\033[0m] Configuration: %s\n", cfg.Json())
}
