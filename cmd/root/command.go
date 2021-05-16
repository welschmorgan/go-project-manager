package root

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/welschmorgan/go-project-manager/cmd/config"
	initCommand "github.com/welschmorgan/go-project-manager/cmd/init"
)

var Command = &cobra.Command{
	Use:   "grlm [commands]",
	Short: "Release multiple projects in a single go",
	Long:  `GRLM allows releasing multiple projects declared in a workspace`,
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
		config.WorkingDirectory = cwd
	}

	Command.PersistentFlags().StringVarP(&config.CfgFile, "config", "c", "", "config file (default is $HOME/.grlm.yaml)")
	Command.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "show additionnal log messages")
	Command.PersistentFlags().StringVarP(&config.WorkingDirectory, "change-directory", "C", config.WorkingDirectory, "change working directory")
	// Command.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// Command.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// Command.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", Command.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", Command.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")
	Command.PersistentFlags().StringVar(&config.WorkspacesRoot, "workspaces-root", "$HOME/projects", "The root folder where to find workspaces")
	viper.BindPFlag("workspaces_root", Command.PersistentFlags().Lookup("workspaces-root"))

	// Command.ActionAddCommand(addCmd)
	Command.AddCommand(initCommand.Command)
}

func initConfig() {
	if config.CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(config.CfgFile)
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
		}
		panic(fmt.Errorf("configuration error: %s", err))
	}
	if config.Verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	if _, err := os.Stat(config.WorkingDirectory); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(config.WorkingDirectory, 0755); err != nil {
			fmt.Printf("error: %s\n", err.Error())
		}
	}
	if err := os.Chdir(config.WorkingDirectory); err != nil {
		panic(err.Error())
	}
}
