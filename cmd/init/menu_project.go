package init

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

type ProjectMenu struct {
	*ui.CRUDMenu
}

func NewProjectMenu(workspace *models.Workspace) (*ProjectMenu, error) {
	if menu, err := ui.NewCRUDMenu(
		workspace,
		"Projects", "Name", &models.Project{},
		[]ui.CRUDAction{ui.ActionQuit, ui.ActionAdd, ui.ActionEdit, ui.ActionRemove, ui.ActionClear},
		map[uint8]string{
			ui.ActionAdd.Id:    "Add new project",
			ui.ActionEdit.Id:   "Edit existing project",
			ui.ActionRemove.Id: "Remove existing project",
			ui.ActionClear.Id:  "Clear projects",
		}); err != nil {
		return nil, err
	} else {
		return &ProjectMenu{
			CRUDMenu: menu,
		}, nil
	}
}
func (m *ProjectMenu) Discover() error {
	var cwd string
	var err error
	var entries []fs.DirEntry
	if cwd, err = os.Getwd(); err != nil {
		return err
	}
	if entries, err = os.ReadDir(cwd); err != nil {
		return err
	}
	for _, dir := range entries {
		if dir.IsDir() && !strings.HasPrefix(strings.TrimSpace(dir.Name()), ".") {
			if sourceControl, err := vcs.Open(filepath.Join(cwd, dir.Name())); err != nil {
				log.Printf("failed to open folder '%s'", err.Error())
			} else {
				if id, ok := m.Indices[dir.Name()]; ok {
					if err = m.Edit(id, models.NewProject(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name())); err != nil {
						return err
					}
				} else {
					m.Create(models.NewProject(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name()))
				}
			}
		}
	}
	return nil
}
