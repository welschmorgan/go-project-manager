package init

import (
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type DevelopperMenu struct {
	*ui.CRUDMenu
}

func validateDevelopper(k, v string) error {
	switch k {
	case "Name":
		return ui.StrMustBeNonEmpty(v)
	case "Email":
		return ui.StrMustBeNonEmpty(v)
	case "Phone":
		return nil
	}
	return nil
}

func NewDevelopperMenu(workspace *config.Workspace) (*DevelopperMenu, error) {
	if menu, err := ui.NewCRUDMenu(
		workspace,
		"Developpers", "Name", config.Person{},
		[]ui.ObjValidator{
			validateDevelopper,
		},
		[]ui.CRUDAction{ui.ActionQuit, ui.ActionAdd, ui.ActionEdit, ui.ActionRemove, ui.ActionClear},
		map[uint8]string{
			ui.ActionAdd.Id:    "Add new developper",
			ui.ActionEdit.Id:   "Edit existing developper",
			ui.ActionRemove.Id: "Remove existing developper",
			ui.ActionClear.Id:  "Clear developpers",
		}, map[string]ui.ItemFieldType{
			"Name":  ui.NewItemFieldType(ui.ItemFieldText, ""),
			"Email": ui.NewItemFieldType(ui.ItemFieldText, ""),
			"Phone": ui.NewItemFieldType(ui.ItemFieldText, ""),
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
		if projectDeveloppers, err := s.ListAuthors(nil); err != nil {
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
