package root

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag"
	initCommand "github.com/welschmorgan/go-release-manager/cmd/init"
	releaseCommand "github.com/welschmorgan/go-release-manager/cmd/release"
	undoCommand "github.com/welschmorgan/go-release-manager/cmd/undo"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
)

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

	// ⑤ Define the CLI flag parameters for your wrapped enum flag.
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

	// define workspaces root
	Command.PersistentFlags().StringVar(&config.Get().WorkspacesRoot, "workspaces_root", config.Get().WorkspacesRoot, "The root folder where to find workspaces")
	viper.BindPFlag("workspaces_root", Command.PersistentFlags().Lookup("workspaces_root"))

	// define log output dir
	Command.PersistentFlags().StringVar(&config.Get().LogFolder, "log_folder", config.Get().LogFolder, "change where to write logs")
	viper.BindPFlag("log_folder", Command.PersistentFlags().Lookup("log_folder"))

	// Command.ActionAddCommand(addCmd)
	Command.AddCommand(initCommand.Command)
	Command.AddCommand(releaseCommand.Command)
	Command.AddCommand(undoCommand.Command)
}

func initConfig() {
	if config.Get().CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config.Get().CfgFile)
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

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("configuration error: %s", err))
		}
	}

	if len(config.Get().WorkingDirectory) != 0 {
		config.Get().Workspace.SetPath(config.Get().WorkingDirectory)
		config.Get().Workspace.Name = filepath.Base(config.Get().Workspace.Path())
	}
	if _, err := os.Stat(config.Get().WorkingDirectory); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(config.Get().WorkingDirectory, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "\033[1;33merror\033[0m: %s\n", err.Error())
		}
	}

	if err := log.Setup(); err != nil {
		panic(err.Error())
	}
	log.Debugf("[\033[1;34m+\033[0m] Using config file: %s\n", viper.ConfigFileUsed())

	if err := os.Chdir(config.Get().WorkingDirectory); err != nil {
		panic(err.Error())
	}
	config.Get().WorkspacePath = filepath.Join(config.Get().WorkingDirectory, config.Get().WorkspaceFilename)
	if _, err := os.Stat(config.Get().WorkspacePath); err == nil || os.IsExist(err) {
		log.Infof("[\033[1;34m+\033[0m] Using local config file: %s\n", config.Get().WorkspacePath)
		if err = config.Get().Workspace.ReadFile(config.Get().WorkspacePath); err != nil {
			panic(err.Error())
		}
		dotDir := regexp.MustCompile(`^\s*\./`)
		wkspDir := config.Get().Workspace.Path()
		if !strings.HasSuffix(wkspDir, string(os.PathSeparator)) {
			wkspDir += string(os.PathSeparator)
		}
		for _, proj := range config.Get().Workspace.Projects {
			proj.Path = dotDir.ReplaceAllString(proj.Path, wkspDir)
			proj.Path = strings.ReplaceAll(proj.Path, "$WORKSPACE", wkspDir)
		}
	}

	config.Get().Workspace.Versionning.PreReleasePrefix = config.Get().Versionning.PreReleasePrefix

	if content, err := json.MarshalIndent(*config.Get(), "", "  "); err != nil {
		panic(err.Error())
	} else {
		log.Debugf("[\033[1;34m+\033[0m] Configuration: %s\n", content)
	}
}
