package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
)

var Command = &cobra.Command{
	Use:   "init [sub]",
	Short: "Initialize the current folder as a workspace",
	Long: `Initialize the current folder and turns it into a workspace.
This will write '.grlm/workspace.yaml' and will interactively ask a few questions.
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(config.Get().WorkspacePath); err == nil {
			if ret, err := ui.AskYN("Workspace already initialized, do you want to reconfigure it"); err != nil {
				return err
			} else if !ret {
				return errors.New("abort")
			}
		}
		os.MkdirAll(filepath.Dir(config.Get().WorkspaceFilename), 0755)
		println("----------------[ General Infos ]--------------")
		if err = askName(&config.Get().Workspace); err != nil {
			return err
		}
		if err = askPath(&config.Get().Workspace); err != nil {
			return err
		}
		println("----------------[ Projects ]--------------")
		if err = askProjects(&config.Get().Workspace); err != nil {
			return err
		}
		println("----------------[ Author ]--------------")
		if err = askAuthor(&config.Get().Workspace); err != nil {
			return err
		}
		println("----------------[ Managers ]--------------")
		if err = askManager(&config.Get().Workspace); err != nil {
			return err
		}
		println("done asking managers")
		println("----------------[ Developpers ]--------------")
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
		return config.Get().Workspace.WriteFile(config.Get().WorkspacePath)
	},
}
