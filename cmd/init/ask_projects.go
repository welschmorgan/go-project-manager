package init

import (
	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
)

func askProjects(wksp *models.Workspace) error {
	var err error
	var menu *ProjectMenu

	if menu, err = NewProjectMenu(wksp); err != nil {
		return err
	}

	var action string
	var project string
	done := false
	for !done {
		menu.Update()
		menu.RenderProjects()
		if action, err = ui.Select("Action", []string{"ActionQuit", "ActionAdd", "ActionRemove", "ActionEdit", "ActionClear"}, nil); err != nil {
			return err
		}
		if action == "ActionRemove" || action == "ActionEdit" {
			if project, err = ui.Select("Project", menu.names, nil); err != nil {
				return err
			}
		}
		defaultProject := models.Project{}
		if action == "ActionEdit" {
			defaultProject = *menu.Get(project)
		}
		if action == "ActionEdit" || action == "ActionAdd" {
			if res, err := ui.AskProject("Project", &defaultProject, nil); err != nil {
				return err
			} else if action == "ActionEdit" {
				if err := menu.ActionEdit(menu.indices[defaultProject.Name], res.Name, res.Path, res.Url, res.SourceControl); err != nil {
					return err
				}
			} else if action == "ActionAdd" {
				if err := menu.Create(res.Name, res.Path, res.Url, res.SourceControl); err != nil {
					return err
				}
			}
		}
		if action == "ActionRemove" {
			menu.ActionRemove(project)
		}
		if action == "ActionClear" {
			menu.ActionClear()
		}
		if action == "ActionQuit" {
			done = true
		}
	}
	wksp.Projects = menu.projects
	return nil
}
