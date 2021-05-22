package init

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/cmd/config"
	"github.com/welschmorgan/go-project-manager/ui"
	"gopkg.in/yaml.v2"
)

var Command = &cobra.Command{
	Use:   "init [sub]",
	Short: "Initialize the current folder as a workspace",
	Long: `Initialize the current folder and turns it into a workspace.
This will write '.grlm-workspace.yaml' and will interactively ask a few questions.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(config.WorkspacePath); err == nil {
			if ret, err := ui.AskYN("Workspace already initialized, do you want to reconfigure it"); err != nil {
				return err
			} else if !ret {
				return errors.New("abort")
			}
		}
		if err = askName(&config.Workspace); err != nil {
			return err
		}
		if err = askPath(&config.Workspace); err != nil {
			return err
		}
		if err = askProjects(&config.Workspace); err != nil {
			return err
		}
		if err = askAuthor(&config.Workspace); err != nil {
			return err
		}
		if err = askManager(&config.Workspace); err != nil {
			return err
		}
		if err = askDeveloppers(&config.Workspace); err != nil {
			return err
		}
		if yaml, err := yaml.Marshal(&config.Workspace); err != nil {
			panic(err.Error())
		} else {
			if err := os.WriteFile(config.WorkspacePath, yaml, 0755); err != nil {
				return err
			}
			fmt.Printf("Written '%s':\n%s\n", config.WorkspacePath, yaml)
		}
		return nil
	},
}
