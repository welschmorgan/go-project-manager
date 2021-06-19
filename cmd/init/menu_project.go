package init

import (
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

type ProjectMenu struct {
	*ui.CRUDMenu
}

func validateProject(k, v string) error {
	switch k {
	case "Name":
		return ui.StrMustBeNonEmpty(v)
	case "Path":
		if err := ui.StrMustBeNonEmpty(v); err != nil {
			return err
		}
	case "Url":
		return nil
	case "SourceControl":
		if err := ui.StrMustBeNonEmpty(v); err != nil {
			return err
		}
		ok := false
		names := []string{}
		for _, s := range vcs.All {
			names = append(names, s.Name())
			if s.Name() == v {
				ok = true
			}
		}
		if !ok {
			return fmt.Errorf("unknown vcs '%s', allowed: [%v]", v, names)
		}
	}
	return nil
}

func NewProjectMenu(workspace *config.Workspace) (*ProjectMenu, error) {
	if menu, err := ui.NewCRUDMenu(
		workspace,
		"Projects", "Name", &config.Project{
			Name: fmt.Sprintf("project %2.2d", rand.Int()),
			Path: workspace.GetPath() + "/",
		},
		[]ui.ObjValidator{
			validateProject,
		},
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
					if err = m.Edit(id, config.NewProject(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name())); err != nil {
						return err
					}
				} else {
					m.Create(config.NewProject(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name()))
				}
			}
		}
	}
	return nil
}
