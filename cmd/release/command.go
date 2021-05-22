package release

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/cmd/config"
)

var Command = &cobra.Command{
	Use:   "release [OPTIONS...]",
	Short: "Release all projects included in this workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if _, err := os.Stat(config.WorkspacePath); err != nil && os.IsNotExist(err) {
			panic("Workspace has not been initialized, run `grlm init`")
		}
		return nil
	},
}
