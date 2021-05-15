package cmd

import (
	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/config"
	"github.com/welschmorgan/go-project-manager/ui"
)

var (
	initCmd = &cobra.Command{
		Use:   "init [sub]",
		Short: "Initialize the current folder as a workspace",
		Long: `Initialize the current folder and turns it into a workspace.
This will write 'grlm.worspace.yaml' and will interactively ask a few questions.
`,
		Run: func(cmd *cobra.Command, args []string) {
			wksp := config.Workspace{}
			if res, err := ui.Select("Source Control Type", []string{"git"}, nil); err != nil {
				panic(err.Error())
			} else {
				wksp.SourceControl = res
			}
		},
	}
)
