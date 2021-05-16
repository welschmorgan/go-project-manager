package init

import (
	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

type DevelopperMenu struct {
	*ui.CRUDMenu
}

func NewDevelopperMenu(workspace *models.Workspace) (*DevelopperMenu, error) {
	if menu, err := ui.NewCRUDMenu(
		workspace,
		"Developpers", "Name", models.Person{},
		[]ui.CRUDAction{ui.ActionQuit, ui.ActionAdd, ui.ActionEdit, ui.ActionRemove, ui.ActionClear},
		map[uint8]string{
			ui.ActionAdd.Id:    "Add new developper",
			ui.ActionEdit.Id:   "Edit existing developper",
			ui.ActionRemove.Id: "Remove existing developper",
			ui.ActionClear.Id:  "Clear developpers",
		}); err != nil {
		return nil, err
	} else {
		return &DevelopperMenu{
			CRUDMenu: menu,
		}, nil
	}
}
func (m *DevelopperMenu) Discover() error {
	for _, project := range m.Workspace.Projects {
		s := vcs.Get(project.SourceControl)
		if err := s.Open(project.Path); err != nil {
			return err
		}
		if projectDeveloppers, err := s.Authors(nil); err != nil {
			return err
		} else {
			for _, tmpDev := range projectDeveloppers {
				if id, ok := m.Indices[tmpDev.Name]; ok {
					if err := m.Edit(id, m.Items[id]); err != nil {
						return err
					}
				} else {
					m.Create(tmpDev)
				}
			}
		}
	}
	return nil
}
