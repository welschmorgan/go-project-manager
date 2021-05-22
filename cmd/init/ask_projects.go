package init

import (
	"github.com/welschmorgan/go-project-manager/models"
)

func askProjects(wksp *models.Workspace) error {
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
