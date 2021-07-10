package version

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/project/accessor"
	"github.com/welschmorgan/go-release-manager/version"
)

var projects []*config.Project

var Command = &cobra.Command{
	Use:   "version [project_name...]",
	Short: "Show project versions",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		unknownProjects := []string{}
		if len(args) > 0 {
			projects = make([]*config.Project, 0)
			for i := 0; i < len(args); i++ {
				var projectFound *config.Project = nil
				for _, p := range config.Get().Workspace.Projects {
					if strings.EqualFold(p.Name, args[0]) {
						projectFound = p
					}
				}
				if projectFound == nil {
					unknownProjects = append(unknownProjects, args[0])
				} else {
					projects = append(projects, projectFound)
				}
			}
		} else {
			projects = config.Get().Workspace.Projects
		}
		if len(unknownProjects) > 0 {
			return fmt.Errorf("unknown projects given in workspace '%s': %s", config.Get().Workspace.Name, strings.Join(unknownProjects, ", "))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cfg := config.Get()
		if !cfg.Workspace.Initialized {
			panic("Workspace has not been initialized yet, run `grlm init`")
		}
		for _, p := range projects {
			if a, err := accessor.Open(p.Path); err != nil {
				return err
			} else {
				var v version.Version
				if v, err = a.ReadVersion(); err != nil {
					return err
				} else {
					fmt.Printf("%s: %s\n", p.Name, v)
				}
			}
		}
		return nil
	},
}
