package init

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
	"gopkg.in/yaml.v2"
)

var Command = &cobra.Command{
	Use:   "init [sub]",
	Short: "Initialize the current folder as a workspace",
	Long: `Initialize the current folder and turns it into a workspace.
This will write '.grlm-workspace.yaml' and will interactively ask a few questions.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(config.Get().WorkspacePath); err == nil {
			if ret, err := ui.AskYN("Workspace already initialized, do you want to reconfigure it"); err != nil {
				return err
			} else if !ret {
				return errors.New("abort")
			}
		}
		if err = askName(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askPath(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askProjects(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askAuthor(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askManager(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askDeveloppers(&config.Get().Workspace); err != nil {
			return err
		}

		for _, proj := range config.Get().Workspace.Projects {
			if _, err := os.Stat(proj.Path); err != nil {
				if os.IsNotExist(err) {
					if ok, err := ui.AskYN(fmt.Sprintf("'%s' does not exist, create project folder", proj.Path)); err != nil {
						return err
					} else if ok {
						if err = os.MkdirAll(proj.Path, 0755); err != nil {
							return err
						}
					}
				}
			}
		}

		if yaml, err := yaml.Marshal(&config.Get().Workspace); err != nil {
			panic(err.Error())
		} else {
			if err := os.WriteFile(config.Get().WorkspacePath, yaml, 0755); err != nil {
				return err
			}
			fmt.Printf("Written '%s':\n%s\n", config.Get().WorkspacePath, yaml)
		}
		return nil
	},
}
