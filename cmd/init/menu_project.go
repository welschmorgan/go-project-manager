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
	"github.com/welschmorgan/go-release-manager/project/accessor"
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
		"Type":          ui.NewItemFieldType(ui.ItemFieldList, accessor.GetAllNames()),
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
		if err = m.Discover(); err != nil {
			return nil, err
		}
		return m, nil
	}
}

func (m *ProjectMenu) CreateFinalizationContext(workspace *config.Workspace, proj *config.Project) (*accessor.FinalizationContext, error) {
	if len(proj.Type) == 0 {
		if ans, err := ui.Select("What type of project is it", accessor.GetAllNames()); err != nil {
			return nil, err
		} else {
			proj.Type = ans
		}
	}
	a := accessor.Get(proj.Type)
	a.Open(proj.Path)
	v := vcs.Get(proj.SourceControl)
	ctx := accessor.NewFinalizationContext(v, a, proj, workspace)
	if err := v.Detect(proj.Path); err == nil {
		ctx.RepositoryInitialized = true
		vcs.SHOW_ERRORS = false
		rootCommits := []string{}
		if rootCommits, err = ctx.VC.GetRootCommits(); err == nil && len(rootCommits) > 0 {
			ctx.InitialCommitExists = true
		}
		fmt.Printf("context: %+v\n", *ctx)
		fmt.Printf("rootCommits: %v\n", rootCommits)
		vcs.SHOW_ERRORS = true
	}
	return ctx, nil
}

func (m *ProjectMenu) FinalizeProject(workspace *config.Workspace, item interface{}) error {
	var err error
	var ctx *accessor.FinalizationContext
	proj := item.(config.Project)
	if ctx, err = m.CreateFinalizationContext(workspace, &proj); err != nil {
		return err
	}
	if err = m.createProjectFolder(ctx); err != nil {
		return err
	}
	if err = m.initializeVCS(ctx); err != nil {
		return err
	}
	if ctx.RepositoryInitialized {
		if err = m.checkBranches(ctx); err != nil {
			return err
		}
		if ctx.InitialCommitExists {
			var checkoutOpts vcs.VersionControlOptions = nil

			if !ctx.DevelopExists {
				checkoutOpts = vcs.CheckoutOptions{
					CreateBranch:  true,
					StartingPoint: workspace.BranchNames["production"],
				}
			}
			if err = ctx.VC.Checkout(workspace.BranchNames["development"], checkoutOpts); err != nil {
				return err
			}
			ctx.DevelopExists = true
		}
	}
	a := accessor.Get(ctx.Project.Type)
	if err = a.Scaffold(ctx); err != nil {
		return err
	}
	if ctx.InitialCommitExists && !ctx.DevelopExists {
		if err = ctx.VC.Checkout(workspace.BranchNames["development"], vcs.CheckoutOptions{
			CreateBranch:  true,
			StartingPoint: workspace.BranchNames["production"],
		}); err != nil {
			return err
		}
	}
	return nil
}

func (m *ProjectMenu) createProjectFolder(ctx *accessor.FinalizationContext) error {
	if fi, err := os.Stat(ctx.Project.Path); err != nil {
		if os.IsNotExist(err) {
			if ctx.UserWantsFolderCreation, _ = ui.AskYN("Project folder does not exist, do you want to create it"); ctx.UserWantsFolderCreation {
				if err := os.MkdirAll(ctx.Project.Path, 0755); err != nil {
					return err
				}
			} else {
				return nil
			}
		} else {
			return err
		}
	} else if !fi.IsDir() {
		return fmt.Errorf("%s: not a directory", ctx.Project.Path)
	}
	return nil
}

func (m *ProjectMenu) initializeVCS(ctx *accessor.FinalizationContext) (err error) {
	if !ctx.RepositoryInitialized {
		if ctx.UserWantsVCSInit, _ = ui.AskYN(ctx.Project.SourceControl + " not initialized, do it now"); ctx.UserWantsVCSInit {
			if err = ctx.VC.Initialize(ctx.Project.Path, nil); err != nil {
				return err
			}
			ctx.RepositoryInitialized = true
		} else {
			log.Println("Skipped VCS initialization...")
			return nil
		}
	}
	if err = ctx.VC.Open(ctx.Project.Path); err != nil {
		return err
	}
	if ctx.InitialCommitExists {
		fmt.Printf("%s already initialized\n", ctx.Project.SourceControl)
	}
	return nil
}

func (m *ProjectMenu) checkBranches(ctx *accessor.FinalizationContext) (err error) {
	if ctx.InitialCommitExists {
		if branches, err := ctx.VC.ListBranches(vcs.BranchOptions{}); err != nil {
			return err
		} else {
			for _, branch := range branches {
				if branch == ctx.Workspace.BranchNames["development"] {
					ctx.DevelopExists = true
				} else if branch == ctx.Workspace.BranchNames["production"] {
					ctx.MasterExists = true
				}
			}
		}
		if ctx.InitialCommitExists && !ctx.MasterExists {
			return fmt.Errorf("production branch '%s' wasn't found in '%s'", ctx.Workspace.BranchNames["production"], ctx.Project.Path)
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
		for _, n := range accessor.GetAllNames() {
			a := accessor.Get(n)
			if ok, err := a.Detect(path); err == nil && ok {
				projType = a.AccessorName()
				break
			}
		}
	}
	if len(projType) == 0 {
		if ans, err := ui.Select("Unknown project type, please pick a type", accessor.GetAllNames()); err != nil {
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
					log.Printf("failed to open folder '%s', %s", projFolder, err.Error())
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
				if err = m.Create(newItem); err != nil {
					return err
				}
			}
		}
	}
	return m.CRUDMenu.Discover()
}
