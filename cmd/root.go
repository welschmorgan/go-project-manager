package cmd

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile        string
	workspacesRoot string
	verbose        bool

	rootCmd = &cobra.Command{
		Use:   "grlm [commands]",
		Short: "Release multiple projects in a single go",
		Long:  `GRLM allows releasing multiple projects declared in a workspace`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grlm.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "show additionnal log messages")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")
	rootCmd.PersistentFlags().StringVar(&workspacesRoot, "workspaces-root", "$HOME/projects", "The root folder where to find workspaces")
	viper.BindPFlag("workspaces_root", rootCmd.PersistentFlags().Lookup("workspaces-root"))

	// rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
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
	if verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
