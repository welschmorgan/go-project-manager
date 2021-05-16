package init

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/welschmorgan/go-project-manager/models"
	"github.com/welschmorgan/go-project-manager/ui"
	"github.com/welschmorgan/go-project-manager/vcs"
)

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
