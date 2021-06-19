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
	"github.com/welschmorgan/go-release-manager/project"
	"github.com/welschmorgan/go-release-manager/ui"
	"github.com/welschmorgan/go-release-manager/vcs"
)

const projectMenuItemsKey = "Projects"
const projectMenuItemsSubKey = "Name"

func projectMenuDefaultItem(w *config.Workspace) *config.Project {
	return &config.Project{
		Name: fmt.Sprintf("project %2.2d", rand.Int()),
		Path: w.GetPath() + "/project",
	}
}

var projectMenuActions = []ui.CRUDAction{ui.ActionQuit, ui.ActionAdd, ui.ActionEdit, ui.ActionRemove, ui.ActionClear}
var projectMenuActionLabels = map[uint8]string{
	ui.ActionAdd.Id:    "Add new project",
	ui.ActionEdit.Id:   "Edit existing project",
	ui.ActionRemove.Id: "Remove existing project",
	ui.ActionClear.Id:  "Clear projects",
}

func projectMenuItemFieldTypes(w *config.Workspace) map[string]ui.ItemFieldType {
	defaultItem := projectMenuDefaultItem(w)
	return map[string]ui.ItemFieldType{
		"Type":          ui.NewItemFieldType(ui.ItemFieldList, project.AllNames),
		"Name":          ui.NewItemFieldType(ui.ItemFieldText, defaultItem.Name),
		"Path":          ui.NewItemFieldType(ui.ItemFieldText, defaultItem.Path),
		"Url":           ui.NewItemFieldType(ui.ItemFieldText, defaultItem.Url),
		"SourceControl": ui.NewItemFieldType(ui.ItemFieldList, vcs.AllNames),
	}
}

type ProjectMenu struct {
	*ui.CRUDMenu
}

func validateProject(k, v string) error {
	switch k {
	case "Type":
		return ui.StrMustBeNonEmpty(v)
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
		projectMenuItemsKey, projectMenuItemsSubKey,
		projectMenuDefaultItem(workspace),
		[]ui.ObjValidator{validateProject},
		projectMenuActions,
		projectMenuActionLabels,
		projectMenuItemFieldTypes(workspace),
		nil,
		false); err != nil {
		return nil, err
	} else {
		m := &ProjectMenu{
			CRUDMenu: menu,
		}
		m.Finalizer = m.FinalizeProject
		m.Discover()
		return m, nil
	}
}

func (m *ProjectMenu) FinalizeProject(item interface{}) error {
	projItem := item.(config.Project)
	if fi, err := os.Stat(projItem.Path); err != nil {
		if os.IsNotExist(err) {
			if createFolder, _ := ui.AskYN("Project folder does not exist, do you want to create it"); createFolder {
				if err := os.MkdirAll(projItem.Path, 0755); err != nil {
					return err
				}
			} else {
				return nil
			}
		} else {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("%s: not a directory", projItem.Path)
	}
	v := vcs.Get(projItem.SourceControl)
	if v == nil {
		return fmt.Errorf("%s: unknown vcs", projItem.SourceControl)
	}
	vcsFirstInit := false
	if ok, err := v.Detect(projItem.Path); err != nil || !ok {
		vcsFirstInit = true
		if initGit, _ := ui.AskYN(projItem.SourceControl + " not initialized, do it now"); initGit {
			if err = v.Initialize(projItem.Path, nil); err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("%s already initialized\n", projItem.SourceControl)
	}
	if len(projItem.Type) == 0 {
		if ans, err := ui.Select("What type of project is it", project.AllNames); err != nil {
			return err
		} else {
			projItem.Type = ans
		}
	}
	accessor := project.Get(projItem.Type)
	if ok, err := accessor.Detect(projItem.Path); err != nil || !ok {
		fmt.Printf("Initializing %s project...\n", accessor.AccessorName())
		if err = accessor.Initialize(projItem.Path, &projItem); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProjectMenu) DetectProjectAccessor(path string) (string, error) {
	projType := ""
	for _, p := range m.Workspace.Projects {
		if p.Path == path {
			projType = strings.TrimSpace(p.Type)
			break
		}
	}
	if len(projType) == 0 {
		for _, n := range project.AllNames {
			a := project.Get(n)
			if ok, err := a.Detect(path); err == nil && ok {
				projType = a.AccessorName()
				break
			}
		}
	}
	if len(projType) == 0 {
		if ans, err := ui.Select("Unknown project type, please pick a type", project.AllNames); err != nil {
			return "", err
		} else {
			projType = ans
		}
	}
	return projType, nil
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
	var knownProject *config.Project = nil
	var sourceControl vcs.VersionControlSoftware = nil
	var projFolder = ""
	var projUrl = ""
	var projVCS = ""
	for _, dir := range entries {
		if dir.IsDir() && !strings.HasPrefix(strings.TrimSpace(dir.Name()), ".") {
			knownProject = nil
			sourceControl = nil
			projFolder = filepath.Join(cwd, dir.Name())
			projUrl = ""
			projVCS = ""
			// try and find a corresponding workspace project declaration
			for _, p := range m.Workspace.Projects {
				if p.Path == projFolder {
					knownProject = p
					break
				}
			}
			// if workspace project found, retrieve configured vcs
			if knownProject != nil {
				sourceControl = vcs.Get(knownProject.SourceControl)
			}
			// detect and open the project using VCS
			if sourceControl == nil {
				if sourceControl, err = vcs.Open(projFolder); err != nil {
					log.Printf("failed to open folder '%s' using %s, %s", projFolder, sourceControl.Name(), err.Error())
				}
			} else {
				if err = sourceControl.Open(projFolder); err != nil {
					log.Printf("failed to open folder '%s' using %s, %s", projFolder, sourceControl.Name(), err.Error())
				}
			}
			// extract url and source control name from detected VCS
			if sourceControl != nil {
				projUrl = sourceControl.Url()
				projVCS = sourceControl.Name()
			}
			// find project accessor
			projType := ""
			if projType, err = m.DetectProjectAccessor(projFolder); err != nil {
				return err
			}
			// create new entry
			newItem := *config.NewProject(projType, dir.Name(), projFolder, projUrl, projVCS)
			if knownProject != nil {
				newItem.Name = knownProject.Name
				newItem.Path = knownProject.Path
			}
			if id, ok := m.Indices[newItem.Name]; ok {
				if err = m.Edit(id, newItem); err != nil {
					return err
				}
			} else {
				if err := m.Create(newItem); err != nil {
					return err
				}
			}
		}
	}
	return m.CRUDMenu.Discover()
}
