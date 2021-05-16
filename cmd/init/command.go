package init

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/cmd/config"
	"github.com/welschmorgan/go-project-manager/models"
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
		wksp := models.Workspace{}
		path := filepath.Join(config.WorkingDirectory, config.WorkspaceFilename)
		if _, err := os.Stat(path); err == nil {
			if ret, err := ui.AskYN("Workspace already initialized, do you want to reconfigure it"); err != nil {
				return err
			} else if !ret {
				return errors.New("abort")
			}
			if content, err := os.ReadFile(path); err != nil {
				return err
			} else if err = yaml.Unmarshal(content, &wksp); err != nil {
				return err
			}
		}
		if err = askName(&wksp); err != nil {
			return err
		}
		if err = askPath(&wksp); err != nil {
			return err
		}
		if err = askProjects(&wksp); err != nil {
			return err
		}
		if err = askAuthor(&wksp); err != nil {
			return err
		}
		if err = askManager(&wksp); err != nil {
			return err
		}
		if err = askDeveloppers(&wksp); err != nil {
			return err
		}
		if yaml, err := yaml.Marshal(&wksp); err != nil {
			panic(err.Error())
		} else {
			if err := os.WriteFile(path, yaml, 0755); err != nil {
				return err
			}
			fmt.Printf("Written '%s':\n%s\n", path, yaml)
		}
		return nil
	},
}
