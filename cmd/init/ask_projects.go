package init

import "github.com/welschmorgan/go-project-manager/config"

func askProjects(wksp *config.Workspace) error {
	var err error
	var menu *ProjectMenu
	if menu, err = NewProjectMenu(wksp); err != nil {
		panic(err)
	} else {
		if err := menu.Render(); err != nil {
			panic(err)
		}
	}
	// wksp.Projects = menu.projects
	return nil
}
