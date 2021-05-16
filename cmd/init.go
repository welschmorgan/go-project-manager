package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"log"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
	"gopkg.in/yaml.v2"
)

func strMustBeNonEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("value must be non-empty")
	}
	return nil
}

func strMustNotContainOnlySpaces(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("value must not contain only spaces")
	}
	return nil
}

func pathMustExist(p string) error {
	if _, err := os.Stat(p); err != nil && os.IsNotExist(err) {
		return errors.New("path does not exist")
	}
	return nil
}

func pathMustBeDir(p string) error {
	if fi, err := os.Stat(p); err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	return nil
}

func askName(wksp *models.Workspace) error {
	var dir string
	var err error
	if len(strings.TrimSpace(wksp.Name)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		wksp.Name = path.Base(dir)
	}
	if wksp.Name, err = ui.Ask("Name", wksp.Name, strMustNotContainOnlySpaces); err != nil {
		return err
	}
	return nil
}

func askPath(wksp *models.Workspace) error {
	var dir string
	var err error
	if len(strings.TrimSpace(wksp.Path)) == 0 {
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		wksp.Path = dir
	}
	if wksp.Path, err = ui.Ask("Path", wksp.Path, strMustBeNonEmpty, strMustNotContainOnlySpaces, pathMustBeDir); err != nil {
		return err
	}
	return nil
}

func askProjects(wksp *models.Workspace) error {
	var projects []*models.Project = wksp.Projects
	if projects == nil {
		projects = []*models.Project{}
	}
	var projectNames []string = make([]string, 0)
	var projectIds map[string]int = map[string]int{}
	var entries []fs.DirEntry
	var cwd string
	var err error
	getProjectId := func(name string) int {
		for id, p := range projects {
			if p != nil && p.Name == name {
				return id
			}
		}
		return -1
	}
	printProjects := func() {
		projectNames = []string{}
		projectIds = map[string]int{}
		s := fmt.Sprintf("Found %d projects: ", len(projects))
		for id, proj := range projects {
			if id > 0 {
				s += ", "
			}
			s += proj.Name
			projectNames = append(projectNames, proj.Name)
			projectIds[proj.Name] = len(projectNames) - 1
		}
		println(s)
	}
	if cwd, err = os.Getwd(); err != nil {
		return err
	} else if entries, err = os.ReadDir(cwd); err != nil {
		return err
	} else {
		for _, dir := range entries {
			if dir.IsDir() && !strings.HasPrefix(strings.TrimSpace(dir.Name()), ".") {
				if sourceControl, err := vcs.Open(filepath.Join(cwd, dir.Name())); err != nil {
					log.Printf("failed to open folder '%s'", err.Error())
				} else {
					id := getProjectId(dir.Name())
					if id != -1 {
						project := projects[id]
						if len(strings.TrimSpace(project.Name)) == 0 {
							project.Name = dir.Name()
						}
						if len(strings.TrimSpace(project.Path)) == 0 {
							project.Path = sourceControl.Path()
						}
						if len(strings.TrimSpace(project.SourceControl)) == 0 {
							project.SourceControl = sourceControl.Name()
						}
						if len(strings.TrimSpace(project.SourceControl)) == 0 {
							project.Url = sourceControl.Url()
						}
					} else {
						projects = append(projects, models.NewProject(dir.Name(), sourceControl.Path(), sourceControl.Url(), sourceControl.Name()))
					}
				}
			}
		}
	}

	var action string
	var project string
	done := false
	for !done {
		printProjects()
		if action, err = ui.Select("Action", []string{"Quit", "Add", "Remove", "Edit", "Clear"}, nil); err != nil {
			return err
		}
		if action == "Remove" || action == "Edit" {
			if project, err = ui.Select("Project", projectNames, nil); err != nil {
				return err
			}
		}
		defaultProject := models.Project{}
		if action == "Edit" {
			id := projectIds[project]
			defaultProject.Name = projects[id].Name
			defaultProject.Path = projects[id].Path
			defaultProject.Url = projects[id].Url
			defaultProject.SourceControl = projects[id].SourceControl
		}
		if action == "Edit" || action == "Add" {
			if res, err := ui.AskProject("Project", &defaultProject, nil); err != nil {
				return err
			} else if len(strings.TrimSpace(res.Name)) == 0 {
				done = true
			} else {
				projects = append(projects, res)
				if action == "Edit" {
					oldId := projectIds[defaultProject.Name]
					projects = append(projects[:oldId], projects[oldId+1:]...)
				}
			}
		}
		if action == "Remove" {
			id := projectIds[project]
			projects = append(projects[:id], projects[id+1:]...)
		}
		if action == "Clear" {
			projects = []*models.Project{}
		}
		if action == "Quit" {
			done = true
		}
	}
	wksp.Projects = projects
	return nil
}

func askAuthor(wksp *models.Workspace) error {
	var currentUser *user.User
	var err error
	var defaultAuthor *models.Person = wksp.Author
	if defaultAuthor != nil && len(strings.TrimSpace(defaultAuthor.Name)) == 0 {
		if currentUser, err = user.Current(); err != nil {
			return err
		}
		defaultAuthor.Name = currentUser.Name
		if len(strings.TrimSpace(defaultAuthor.Name)) == 0 {
			defaultAuthor.Name = currentUser.Username
		}
	}
	if wksp.Author, err = ui.AskPerson("Author", defaultAuthor); err != nil {
		return err
	}
	return nil
}

func askManager(wksp *models.Workspace) error {
	var username string
	var err error
	if wksp.Manager != nil {
		username = wksp.Manager.Name
	}
	if len(strings.TrimSpace(username)) == 0 && wksp.Author != nil {
		username = wksp.Author.Name
	}
	if wksp.Manager, err = ui.AskPerson("Manager", &models.Person{Name: username}); err != nil {
		return err
	}
	return nil
}

func askDeveloppers(wksp *models.Workspace) error {
	var developpers []*models.Person = wksp.Developpers
	if developpers == nil {
		developpers = []*models.Person{}
	}
	developperNames := []string{}
	developperIds := map[string]int{}
	getDevelopperId := func(name string) int {
		for id, d := range developpers {
			if d != nil && d.Name == name {
				return id
			}
		}
		return -1
	}
	for _, project := range wksp.Projects {
		s := vcs.Get(project.SourceControl)
		if err := s.Open(project.Path); err != nil {
			return err
		}
		if projectDeveloppers, err := s.Authors(nil); err != nil {
			return err
		} else {
			for _, tmpDev := range projectDeveloppers {
				id := getDevelopperId(tmpDev.Name)
				if id != -1 {
					dev := developpers[id]
					if len(strings.TrimSpace(dev.Name)) == 0 {
						dev.Name = tmpDev.Name
					}
					if len(strings.TrimSpace(dev.Phone)) == 0 {
						dev.Phone = tmpDev.Phone
					}
					if len(strings.TrimSpace(dev.Email)) == 0 {
						dev.Email = tmpDev.Email
					}
				} else {
					developpers = append(developpers, tmpDev)
				}
			}
		}
	}
	printDeveloppers := func() {
		developperNames = []string{}
		developperIds = map[string]int{}
		s := fmt.Sprintf("Found %d developpers: ", len(developpers))
		for id, a := range developpers {
			if id > 0 {
				s += ", "
			}
			s += a.Name
			developperNames = append(developperNames, a.Name)
			developperIds[a.Name] = id
		}
		println(s)
	}
	done := false
	var action string
	var developperName string
	var err error
	for !done {
		printDeveloppers()
		if action, err = ui.Select("Action", []string{"Quit", "Add", "Remove", "Edit", "Clear"}, nil); err != nil {
			return err
		}
		if action == "Edit" || action == "Remove" {
			if developperName, err = ui.Select("Developper", developperNames, nil); err != nil {
				return err
			}
		}
		defaultDevelopper := models.Person{}
		if action == "Edit" {
			developper := developpers[developperIds[developperName]]
			defaultDevelopper.Name = developper.Name
			defaultDevelopper.Email = developper.Email
			defaultDevelopper.Phone = developper.Phone
		}
		if action == "Edit" || action == "Add" {
			if auth, err := ui.AskPerson("Developper", &defaultDevelopper, nil); err != nil {
				return err
			} else {
				oldId := developperIds[auth.Name]
				developpers = append(developpers, auth)
				if action == "Edit" {
					developpers = append(developpers[:oldId], developpers[oldId+1:]...)
				}
			}
		}
		if action == "Remove" {
			id := developperIds[developperName]
			developpers = append(developpers[:id], developpers[id+1:]...)
		}
		if action == "Clear" {
			developpers = []*models.Person{}
		}
		if action == "Quit" {
			done = true
		}
	}
	wksp.Developpers = developpers
	return nil
}

var (
	initCmd = &cobra.Command{
		Use:   "init [sub]",
		Short: "Initialize the current folder as a workspace",
		Long: `Initialize the current folder and turns it into a workspace.
This will write '.grlm-workspace.yaml' and will interactively ask a few questions.
`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			wksp := models.Workspace{}
			path := filepath.Join(workingDirectory, workspaceFilename)
			if _, err := os.Stat(path); err == nil {
				if ret, err := ui.AskYN("Workspace already initialized, do you want to reconfigure it"); err != nil {
					return err
				} else if !ret {
					return errors.New("abort")
				}
				if content, err := os.ReadFile(path); err != nil {
					return err
				} else if err = yaml.Unmarshal(content, &wksp); err != nil {
					return err
				}
			}
			if err = askName(&wksp); err != nil {
				return err
			}
			if err = askPath(&wksp); err != nil {
				return err
			}
			if err = askProjects(&wksp); err != nil {
				return err
			}
			if err = askAuthor(&wksp); err != nil {
				return err
			}
			if err = askManager(&wksp); err != nil {
				return err
			}
			if err = askDeveloppers(&wksp); err != nil {
				return err
			}
			if yaml, err := yaml.Marshal(&wksp); err != nil {
				panic(err.Error())
			} else {
				if err := os.WriteFile(path, yaml, 0755); err != nil {
					return err
				}
				fmt.Printf("Written '%s':\n%s\n", path, yaml)
			}
			return nil
		},
	}
)
